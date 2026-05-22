package cli

import (
	"context"
	"strings"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/model"
)

func TestAntigravityCollisionCheckIncludesCLI(t *testing.T) {
	checks := antigravityCollisionCheck([]model.AgentID{model.AgentGeminiCLI, model.AgentAntigravityCLI})
	if len(checks) != 1 {
		t.Fatalf("antigravityCollisionCheck() len = %d, want 1", len(checks))
	}

	check := checks[0]
	if !check.Soft {
		t.Fatal("antigravityCollisionCheck() should return a soft warning")
	}
	if check.ID != "verify:antigravity:rules-collision" {
		t.Fatalf("check ID = %q, want verify:antigravity:rules-collision", check.ID)
	}

	err := check.Run(context.Background())
	if err == nil {
		t.Fatal("check.Run() error = nil, want warning message")
	}
	message := err.Error()
	for _, want := range []string{
		"Antigravity CLI intentionally uses the Gemini-compatible global prompt surface",
		"last synced SDD orchestrator owns the shared gentle-ai:sdd-orchestrator section",
		"Prefer Antigravity CLI for new installs",
	} {
		if !strings.Contains(message, want) {
			t.Fatalf("warning message missing %q; got:\n%s", want, message)
		}
	}
}

func TestAntigravityCollisionCheckNoWarningWithoutGemini(t *testing.T) {
	checks := antigravityCollisionCheck([]model.AgentID{model.AgentAntigravityCLI})
	if len(checks) != 0 {
		t.Fatalf("antigravityCollisionCheck() len = %d, want 0", len(checks))
	}
}
