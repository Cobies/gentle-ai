package screens

import (
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/model"
)

func TestAgentOptionsHidesLegacyAntigravity(t *testing.T) {
	options := AgentOptions()

	seenCLI := false
	for _, option := range options {
		if option == model.AgentAntigravity {
			t.Fatal("AgentOptions() should hide legacy Antigravity; Antigravity CLI supersedes it in the TUI")
		}
		if option == model.AgentAntigravityCLI {
			seenCLI = true
		}
	}

	if !seenCLI {
		t.Fatal("AgentOptions() missing Antigravity CLI option")
	}
}
