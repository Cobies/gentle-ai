package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gentleman-programming/gentle-ai/v3/internal/components/filemerge"
	"github.com/gentleman-programming/gentle-ai/v3/internal/system"
)

// kiloParityAgentNames is the v3.7.0 non-SDD Kilo agent set (#4471). Kilo
// never received review-validator (it hosts no provider relay) nor the
// gentle-ai-* ODD trio (its orchestrator routes to native delegation).
var kiloParityAgentNames = []string{
	"jd-judge-a", "jd-judge-b", "jd-fix-agent",
	"review-risk", "review-readability", "review-reliability", "review-resilience",
	"review-refuter",
}

func kiloSettingsPath(home string) string {
	return filepath.Join(home, ".config", "kilo", "opencode.json")
}

func kiloAgents(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	root, err := filemerge.UnmarshalJSONObject(raw)
	if err != nil {
		t.Fatalf("decode Kilo settings: %v", err)
	}
	agents, _ := root["agent"].(map[string]any)
	return agents
}

func assertKiloParityAgents(t *testing.T, agents map[string]any) {
	t.Helper()
	for _, name := range kiloParityAgentNames {
		entry, ok := agents[name].(map[string]any)
		if !ok {
			t.Fatalf("Kilo settings missing agent %q; got %v", name, agents)
		}
		if entry["mode"] != "subagent" || entry["hidden"] != true {
			t.Fatalf("Kilo agent %q missing mode/hidden: %#v", name, entry)
		}
		if prompt, _ := entry["prompt"].(string); strings.TrimSpace(prompt) == "" || strings.Contains(prompt, "obsolete") {
			t.Fatalf("Kilo agent %q has no current prompt: %#v", name, entry)
		}
		if _, ok := entry["permission"].(map[string]any); !ok {
			t.Fatalf("Kilo agent %q missing permission: %#v", name, entry)
		}
	}
	for _, name := range []string{"review-validator", "gentle-ai-explore", "gentle-ai-verify", "gentle-ai-worker"} {
		if _, ok := agents[name]; ok {
			t.Fatalf("Kilo must not install %q", name)
		}
	}
	for name, raw := range agents {
		entry, _ := raw.(map[string]any)
		if _, marked := entry["__managed_by"]; marked {
			t.Fatalf("Kilo agent %q still carries __managed_by (#4471): %#v", name, entry)
		}
	}
	orchestrator, _ := agents["gentle-orchestrator"].(map[string]any)
	permission, _ := orchestrator["permission"].(map[string]any)
	task, _ := permission["task"].(map[string]any)
	for _, name := range kiloParityAgentNames {
		if task[name] != "allow" {
			t.Fatalf("Kilo gentle-orchestrator does not allow delegating to %q: %#v", name, task)
		}
	}
}

func TestKiloInstallWritesParityAgentsWithoutManagedByMarker(t *testing.T) {
	home := installTestHome(t)
	args := []string{"--agent", "kilocode", "--preset", "full-gentleman"}
	if _, err := RunInstall(args, system.DetectionResult{}); err != nil {
		t.Fatal(err)
	}
	first, err := os.ReadFile(kiloSettingsPath(home))
	if err != nil {
		t.Fatal(err)
	}
	agents := kiloAgents(t, first)
	assertKiloParityAgents(t, agents)
	for _, name := range kiloParityAgentNames {
		for key := range agents[name].(map[string]any) {
			if !openCodeAgentConfigAllowedKeys[key] {
				t.Fatalf("Kilo agent %q carries non-AgentConfig key %q: %#v", name, key, agents[name])
			}
		}
	}
	if _, err := RunInstall(args, system.DetectionResult{}); err != nil {
		t.Fatal(err)
	}
	if second, _ := os.ReadFile(kiloSettingsPath(home)); !bytes.Equal(first, second) {
		t.Fatalf("second Kilo install changed settings:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

func TestKiloUpgradeRetiresOwnedAgents(t *testing.T) {
	for _, tc := range []struct {
		name string
		run  func(t *testing.T)
	}{
		{"install", func(t *testing.T) {
			t.Helper()
			if _, err := RunInstall([]string{"--agent", "kilocode", "--preset", "full-gentleman"}, system.DetectionResult{}); err != nil {
				t.Fatal(err)
			}
		}},
		{"sync", func(t *testing.T) {
			t.Helper()
			if _, err := RunSync([]string{"--agent", "kilocode"}); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home := installTestHome(t)
			path := kiloSettingsPath(home)
			if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
				t.Fatal(err)
			}
			fixture, err := os.ReadFile("testdata/kilo-v3.7.0-upgrade.json")
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, fixture, 0o600); err != nil {
				t.Fatal(err)
			}
			before, err := filemerge.UnmarshalJSONObject(fixture)
			if err != nil {
				t.Fatal(err)
			}
			tc.run(t)
			first, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Contains(first, []byte(`"__managed_by"`)) {
				t.Fatalf("legacy marker remains:\n%s", first)
			}
			after, err := filemerge.UnmarshalJSONObject(first)
			if err != nil {
				t.Fatal(err)
			}
			agents := after["agent"].(map[string]any)
			original := before["agent"].(map[string]any)
			assertKiloParityAgents(t, agents)
			for name := range agents {
				if strings.HasPrefix(name, "sdd-") {
					t.Errorf("retired owned agent %s remains", name)
				}
			}
			if !reflect.DeepEqual(after["provider"], before["provider"]) {
				t.Error("provider changed")
			}
			if !reflect.DeepEqual(after["mcp"].(map[string]any)["my-server"], before["mcp"].(map[string]any)["my-server"]) {
				t.Error("user MCP server changed")
			}
			if !reflect.DeepEqual(agents["user-owned"], original["user-owned"]) {
				t.Error("user-owned agent changed")
			}
			if !reflect.DeepEqual(agents["other-marked"], map[string]any{"prompt": "keep this"}) {
				t.Error("unknown marked agent changed beyond marker removal")
			}
			for _, name := range append([]string{"gentle-orchestrator"}, kiloParityAgentNames...) {
				entry := agents[name].(map[string]any)
				if _, stale := entry["tools"]; stale {
					t.Errorf("%s retained stale v3.7.0 tools", name)
				}
			}
			for name, want := range map[string]map[string]any{
				"jd-judge-a":          {"model": "user/judge", "variant": "low"},
				"review-risk":         {"variant": "medium"},
				"gentle-orchestrator": {"model": "user/orchestrator"},
			} {
				entry := agents[name].(map[string]any)
				for field, value := range want {
					if entry[field] != value {
						t.Errorf("%s lost user %s: %v", name, field, entry)
					}
				}
			}
			tc.run(t)
			if second, _ := os.ReadFile(path); !bytes.Equal(first, second) {
				t.Errorf("second %s changed settings bytes", tc.name)
			}
		})
	}
}
