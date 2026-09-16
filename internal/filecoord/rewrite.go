package filecoord

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"

	"github.com/gentleman-programming/gentle-ai/v2/internal/components/filemerge"
)

// Outcome represents the typed result of a cooperative rewrite operation.
type Outcome string

const (
	OutcomeNoChange    Outcome = "no_change"
	OutcomeCommitted   Outcome = "committed"
	OutcomeConflict    Outcome = "conflict"
	OutcomeUnsupported Outcome = "unsupported"
	OutcomeFailed      Outcome = "failed"
)

// RewriteResult reports the outcome and details of a cooperative rewrite.
// NOTE: It is an explicit cooperative optimistic contract, NOT an arbitrary-writer
// CAS; a non-cooperating writer racing the atomic publication window cannot be
// portably prevented.
type RewriteResult struct {
	Outcome        Outcome
	BytesChanged   bool
	ConflictReason string
	Durability     string
}

// Mutator receives the current file content and existence status, returning
// the desired replacement bytes or an error.
type Mutator func(current []byte, exists bool) (newBytes []byte, err error)

// Rewrite coordinates a safe, cooperative optimistic mutation of target:
// 1. Takes the shared cooperative lock.
// 2. Captures an initial no-follow snapshot.
// 3. Computes the replacement via mutate.
// 4. Returns OutcomeNoChange if bytes are unchanged.
// 5. Revalidates point-in-time identity under lock; returns OutcomeConflict if changed.
// 6. Atomically publishes by pathname.
// 7. Reads back to verify durability before reporting OutcomeCommitted.
func Rewrite(ctx context.Context, target, lockRoot string, mutate Mutator) (*RewriteResult, error) {
	lease, err := Acquire(ctx, target, lockRoot)
	if err != nil {
		if errors.Is(err, ErrBusy) {
			return &RewriteResult{Outcome: OutcomeConflict, ConflictReason: "cooperative lock busy"}, nil
		}
		if errors.Is(err, ErrUnsupported) {
			return &RewriteResult{Outcome: OutcomeUnsupported, ConflictReason: "cooperative lock unsupported"}, nil
		}
		return &RewriteResult{Outcome: OutcomeFailed}, err
	}
	defer func() { _ = lease.Release() }()

	snap, err := ReadSnapshot(target)
	if err != nil {
		if errors.Is(err, ErrSymlinkTarget) || errors.Is(err, ErrNonRegularTarget) || errors.Is(err, ErrOversizedTarget) {
			return &RewriteResult{Outcome: OutcomeConflict, ConflictReason: err.Error()}, nil
		}
		return &RewriteResult{Outcome: OutcomeFailed}, err
	}

	newBytes, err := mutate(snap.Bytes, snap.Exists)
	if err != nil {
		return &RewriteResult{Outcome: OutcomeFailed}, fmt.Errorf("mutate %q: %w", target, err)
	}

	if snap.Exists && bytes.Equal(snap.Bytes, newBytes) {
		return &RewriteResult{
			Outcome:      OutcomeNoChange,
			BytesChanged: false,
		}, nil
	}

	// Revalidate under lock immediately before publication
	reval, err := ReadSnapshot(target)
	if err != nil {
		return &RewriteResult{
			Outcome:        OutcomeConflict,
			ConflictReason: fmt.Sprintf("revalidation failed: %v", err),
		}, nil
	}
	if reval.Exists != snap.Exists || !reval.Identity.Matches(snap.Identity) || !bytes.Equal(reval.Bytes, snap.Bytes) {
		return &RewriteResult{
			Outcome:        OutcomeConflict,
			ConflictReason: "target modified concurrently prior to publication",
		}, nil
	}

	perm := fs.FileMode(0o644)
	if snap.Exists && snap.Mode.Perm() != 0 {
		perm = snap.Mode.Perm()
	}

	// Atomically stage, publish, and sync directory
	writeRes, err := filemerge.WriteFileAtomic(target, newBytes, perm)
	if err != nil {
		return &RewriteResult{
			Outcome:      OutcomeFailed,
			BytesChanged: writeRes.Changed,
		}, err
	}

	// Read back to guarantee durability and honest reporting
	postSnap, err := ReadSnapshot(target)
	if err != nil || !postSnap.Exists || !bytes.Equal(postSnap.Bytes, newBytes) {
		return &RewriteResult{
			Outcome:      OutcomeFailed,
			BytesChanged: true,
		}, fmt.Errorf("read-back verification failed for %q", target)
	}

	return &RewriteResult{
		Outcome:      OutcomeCommitted,
		BytesChanged: true,
		Durability:   "synced",
	}, nil
}
