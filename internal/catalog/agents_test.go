package catalog

import (
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/model"
)

func TestAllAgentsIncludesPi(t *testing.T) {
	agents := AllAgents()

	for _, agent := range agents {
		if agent.ID != model.AgentPi {
			continue
		}

		if agent.Name != "Pi" {
			t.Fatalf("Pi Name = %q, want Pi", agent.Name)
		}

		if agent.Tier != model.TierFull {
			t.Fatalf("Pi Tier = %q, want %q", agent.Tier, model.TierFull)
		}

		if agent.ConfigPath != "~/.pi" {
			t.Fatalf("Pi ConfigPath = %q, want ~/.pi", agent.ConfigPath)
		}

		return
	}

	t.Fatalf("AllAgents() missing %s", model.AgentPi)
}

func TestAllAgentsIncludesAntigravityCLI(t *testing.T) {
	agents := AllAgents()

	for _, agent := range agents {
		if agent.ID != model.AgentAntigravityCLI {
			continue
		}

		if agent.Name != "Google Antigravity CLI" {
			t.Fatalf("Antigravity CLI Name = %q, want Google Antigravity CLI", agent.Name)
		}

		if agent.Tier != model.TierFull {
			t.Fatalf("Antigravity CLI Tier = %q, want %q", agent.Tier, model.TierFull)
		}

		if agent.ConfigPath != "~/.gemini/antigravity-cli" {
			t.Fatalf("Antigravity CLI ConfigPath = %q, want ~/.gemini/antigravity-cli", agent.ConfigPath)
		}

		return
	}

	t.Fatalf("AllAgents() missing %s", model.AgentAntigravityCLI)
}

func TestIsSupportedAgentAcceptsPi(t *testing.T) {
	if !IsSupportedAgent(model.AgentPi) {
		t.Fatalf("IsSupportedAgent(%q) = false, want true", model.AgentPi)
	}
}
