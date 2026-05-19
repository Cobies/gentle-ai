package model

import "testing"

func TestAgentAntigravityCLI(t *testing.T) {
	if AgentAntigravityCLI != "antigravity-cli" {
		t.Errorf("AgentAntigravityCLI = %q, want %q", AgentAntigravityCLI, "antigravity-cli")
	}
}
