package filecoord

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestRewriteNoChange(t *testing.T) {
	dir := canonicalBackendTempDir(t)
	lockRoot := filepath.Join(dir, "locks")
	target := filepath.Join(dir, "target.txt")
	initialContent := []byte("constant content")
	if err := os.WriteFile(target, initialContent, 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := Rewrite(context.Background(), target, lockRoot, func(current []byte, exists bool) ([]byte, error) {
		if !exists {
			t.Fatal("expected file to exist")
		}
		return initialContent, nil
	})
	if err != nil {
		t.Fatalf("Rewrite() = %v, want nil", err)
	}
	if res.Outcome != OutcomeNoChange {
		t.Fatalf("res.Outcome = %v, want OutcomeNoChange", res.Outcome)
	}
	if res.BytesChanged {
		t.Fatal("res.BytesChanged = true, want false")
	}
}

func TestRewriteCommittedNewAndExisting(t *testing.T) {
	dir := canonicalBackendTempDir(t)
	lockRoot := filepath.Join(dir, "locks")
	target := filepath.Join(dir, "target.txt")

	// 1. New file creation
	res1, err := Rewrite(context.Background(), target, lockRoot, func(current []byte, exists bool) ([]byte, error) {
		if exists {
			t.Fatal("expected file not to exist initially")
		}
		return []byte("first edition\n"), nil
	})
	if err != nil {
		t.Fatalf("Rewrite(new) = %v", err)
	}
	if res1.Outcome != OutcomeCommitted || !res1.BytesChanged {
		t.Fatalf("res1 = %+v, want OutcomeCommitted", res1)
	}

	data, err := os.ReadFile(target)
	if err != nil || string(data) != "first edition\n" {
		t.Fatalf("read target = %q, %v", data, err)
	}

	// 2. Existing file update
	res2, err := Rewrite(context.Background(), target, lockRoot, func(current []byte, exists bool) ([]byte, error) {
		if !exists || string(current) != "first edition\n" {
			t.Fatalf("unexpected current state: exists=%v, current=%q", exists, current)
		}
		return []byte("second edition\n"), nil
	})
	if err != nil {
		t.Fatalf("Rewrite(update) = %v", err)
	}
	if res2.Outcome != OutcomeCommitted || !res2.BytesChanged {
		t.Fatalf("res2 = %+v, want OutcomeCommitted", res2)
	}

	data2, err := os.ReadFile(target)
	if err != nil || string(data2) != "second edition\n" {
		t.Fatalf("read target updated = %q, %v", data2, err)
	}
}

func TestRewriteConflictConcurrentModification(t *testing.T) {
	dir := canonicalBackendTempDir(t)
	lockRoot := filepath.Join(dir, "locks")
	target := filepath.Join(dir, "target.txt")
	if err := os.WriteFile(target, []byte("initial\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := Rewrite(context.Background(), target, lockRoot, func(current []byte, exists bool) ([]byte, error) {
		// Simulate concurrent external modification during mutator execution
		if writeErr := os.WriteFile(target, []byte("racing writer bytes\n"), 0o644); writeErr != nil {
			t.Fatal(writeErr)
		}
		return []byte("my update\n"), nil
	})
	if err != nil {
		t.Fatalf("Rewrite() returned unexpected error: %v", err)
	}
	if res.Outcome != OutcomeConflict {
		t.Fatalf("res.Outcome = %v, want OutcomeConflict", res.Outcome)
	}
	if res.ConflictReason == "" {
		t.Fatal("res.ConflictReason is empty, want descriptive reason")
	}

	// Verify racing writer bytes were preserved and not clobbered
	data, _ := os.ReadFile(target)
	if string(data) != "racing writer bytes\n" {
		t.Fatalf("target was clobbered: %q", data)
	}
}

func TestRewriteConflictSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink fixtures unavailable on windows")
	}
	dir := canonicalBackendTempDir(t)
	lockRoot := filepath.Join(dir, "locks")
	realFile := filepath.Join(dir, "real.txt")
	if err := os.WriteFile(realFile, []byte("real content"), 0o644); err != nil {
		t.Fatal(err)
	}
	symlink := filepath.Join(dir, "link.txt")
	if err := os.Symlink(realFile, symlink); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}

	res, err := Rewrite(context.Background(), symlink, lockRoot, func(current []byte, exists bool) ([]byte, error) {
		return []byte("new content"), nil
	})
	if err != nil {
		t.Fatalf("Rewrite(symlink) = %v, want nil error with OutcomeConflict", err)
	}
	if res.Outcome != OutcomeConflict {
		t.Fatalf("res.Outcome = %v, want OutcomeConflict", res.Outcome)
	}
}

func TestRewriteConflictLockBusy(t *testing.T) {
	dir := canonicalBackendTempDir(t)
	lockRoot := filepath.Join(dir, "locks")
	target := filepath.Join(dir, "target.txt")

	// Acquire the lock first
	lease, err := Acquire(context.Background(), target, lockRoot)
	if err != nil {
		t.Fatalf("Acquire() = %v", err)
	}
	defer func() { _ = lease.Release() }()

	res, err := Rewrite(context.Background(), target, lockRoot, func(current []byte, exists bool) ([]byte, error) {
		return []byte("new"), nil
	})
	if err != nil {
		t.Fatalf("Rewrite() = %v", err)
	}
	if res.Outcome != OutcomeConflict {
		t.Fatalf("res.Outcome = %v, want OutcomeConflict", res.Outcome)
	}
}

func TestRewriteContextCancelled(t *testing.T) {
	dir := canonicalBackendTempDir(t)
	lockRoot := filepath.Join(dir, "locks")
	target := filepath.Join(dir, "target.txt")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	res, err := Rewrite(ctx, target, lockRoot, func(current []byte, exists bool) ([]byte, error) {
		return []byte("new"), nil
	})
	if res.Outcome != OutcomeFailed {
		t.Fatalf("res.Outcome = %v, want OutcomeFailed", res.Outcome)
	}
	if !errors.Is(err, ErrOperational) && !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled/ErrOperational", err)
	}
}
