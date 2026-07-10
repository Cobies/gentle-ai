package sdd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/agents"
	"github.com/gentleman-programming/gentle-ai/internal/model"
)

func TestAntigravitySddAgentsPluginWritesManifestAndHooks(t *testing.T) {
	home := t.TempDir()

	changed, files, err := installAntigravitySddAgentsPlugin(home)
	if err != nil {
		t.Fatalf("installAntigravitySddAgentsPlugin() error = %v", err)
	}
	if !changed {
		t.Fatalf("installAntigravitySddAgentsPlugin() changed = false, want true")
	}
	if len(files) != 2 {
		t.Fatalf("installAntigravitySddAgentsPlugin() files = %d, want 2", len(files))
	}

	pluginPath := filepath.Join(home, ".gemini", "antigravity-cli", "plugins", "gentle-ai-sdd-agents", "plugin.json")
	hooksPath := filepath.Join(home, ".gemini", "antigravity-cli", "plugins", "gentle-ai-sdd-agents", "hooks.json")

	for _, want := range []string{pluginPath, hooksPath} {
		if !containsString(files, want) {
			t.Fatalf("installAntigravitySddAgentsPlugin() files missing %q; got %v", want, files)
		}
		if _, err := os.Stat(want); err != nil {
			t.Fatalf("expected %q to exist on disk: %v", want, err)
		}
	}

	manifest, err := os.ReadFile(pluginPath)
	if err != nil {
		t.Fatalf("ReadFile(plugin.json) error = %v", err)
	}
	for _, want := range []string{`"name": "gentle-ai-sdd-agents"`, `"version": "0.1.0"`} {
		if !strings.Contains(string(manifest), want) {
			t.Fatalf("plugin.json missing %q:\n%s", want, manifest)
		}
	}

	hooks, err := os.ReadFile(hooksPath)
	if err != nil {
		t.Fatalf("ReadFile(hooks.json) error = %v", err)
	}
	for _, want := range []string{
		`"gentle-ai-sdd-agents-hardening"`,
		`"PreInvocation"`,
		"printf",
		"sdd-explore",
		"sdd-apply",
		"review-*",
		"jd-judge-*",
		"jd-fix-agent",
		"fail closed",
	} {
		if !strings.Contains(string(hooks), want) {
			t.Fatalf("hooks.json missing %q:\n%s", want, hooks)
		}
	}
}

func TestAntigravitySddAgentsPluginMergesIntoExistingHooks(t *testing.T) {
	home := t.TempDir()
	hooksPath := filepath.Join(home, ".gemini", "antigravity-cli", "plugins", "gentle-ai-sdd-agents", "hooks.json")
	if err := os.MkdirAll(filepath.Dir(hooksPath), 0o755); err != nil {
		t.Fatal(err)
	}
	// Pre-existing custom hook must be preserved after install.
	preExisting := `{"custom-hook":{"PreInvocation":[{"type":"command","command":"echo custom"}]}}`
	if err := os.WriteFile(hooksPath, []byte(preExisting), 0o644); err != nil {
		t.Fatal(err)
	}

	changed, _, err := installAntigravitySddAgentsPlugin(home)
	if err != nil {
		t.Fatalf("installAntigravitySddAgentsPlugin() error = %v", err)
	}
	if !changed {
		t.Fatal("installAntigravitySddAgentsPlugin() changed = false on first install over existing hooks")
	}

	merged, err := os.ReadFile(hooksPath)
	if err != nil {
		t.Fatalf("ReadFile(hooks.json) error = %v", err)
	}
	for _, want := range []string{`"custom-hook"`, "echo custom", `"gentle-ai-sdd-agents-hardening"`, "sdd-explore"} {
		if !strings.Contains(string(merged), want) {
			t.Fatalf("merged hooks.json missing %q:\n%s", want, merged)
		}
	}
}

func TestAntigravitySddAgentsPluginIdempotent(t *testing.T) {
	home := t.TempDir()

	if _, _, err := installAntigravitySddAgentsPlugin(home); err != nil {
		t.Fatalf("first install error = %v", err)
	}
	second, files, err := installAntigravitySddAgentsPlugin(home)
	if err != nil {
		t.Fatalf("second install error = %v", err)
	}
	if second {
		t.Fatalf("second install changed = true, want false (idempotent)")
	}
	if len(files) != 2 {
		t.Fatalf("second install files = %d, want 2", len(files))
	}
}

func TestAntigravitySddAgentsHardeningContractPhrases(t *testing.T) {
	msg := AntigravitySddAgentsHardeningMessage()
	if strings.Contains(msg, "'") {
		t.Fatalf("hardening message must not contain single quotes because hooks.json embeds it in a single-quoted shell string: %q", msg)
	}
	for _, want := range antigravitySddAgentsHardeningContractPhrases {
		if !strings.Contains(msg, want) {
			t.Errorf("hardening message missing required phrase %q", want)
		}
	}
	for _, forbidden := range antigravitySddAgentsHardeningContractForbids {
		if strings.Contains(msg, forbidden) {
			t.Errorf("hardening message must not contain %q (it would fake a real Antigravity API)", forbidden)
		}
	}
}

func TestAntigravitySddAgentsRoleAllowed(t *testing.T) {
	allowed := []string{
		"sdd-explore", "sdd-spec", "sdd-design", "sdd-tasks", "sdd-apply",
		"sdd-verify", "sdd-archive", "sdd-onboard", "sdd-init", "sdd-propose",
		"review-risk", "review-resilience", "review-readability", "review-reliability",
		"review-refuter",
		"jd-judge-a", "jd-judge-b", "jd-fix-agent",
	}
	for _, role := range allowed {
		if !antigravitySddAgentsRoleAllowed(role) {
			t.Errorf("role %q must be allowed; got false", role)
		}
	}

	// Roles not in the canonical table must NOT be allowed. Locking this
	// prevents accidental widening of the hardening contract.
	for _, role := range []string{"sdd-evil", "review-anything", "jd-super-user", ""} {
		if antigravitySddAgentsRoleAllowed(role) {
			t.Errorf("role %q must NOT be allowed; got true", role)
		}
	}
}

func TestAntigravitySddAgentsRoleScopesSortedAndUnique(t *testing.T) {
	scopes := sortedAntigravitySddAgentsRoleScopes()
	if len(scopes) != len(antigravitySddAgentsRoleScopes) {
		t.Fatalf("sorted length %d != canonical length %d", len(scopes), len(antigravitySddAgentsRoleScopes))
	}
	seen := map[string]bool{}
	for i, entry := range scopes {
		if seen[entry.Role] {
			t.Errorf("duplicate role %q in scope table", entry.Role)
		}
		seen[entry.Role] = true
		if i > 0 && scopes[i-1].Role > entry.Role {
			t.Errorf("scope table not sorted: %q > %q", scopes[i-1].Role, entry.Role)
		}
		if entry.Scope == "" {
			t.Errorf("role %q has empty scope", entry.Role)
		}
	}
}

func TestHasAntigravitySddAgentsHardeningContractRejectsPhrasesOutsideManagedBlock(t *testing.T) {
	home := t.TempDir()
	hooksPath := filepath.Join(antigravitySddAgentsPluginDir(home), "hooks.json")
	if err := os.MkdirAll(filepath.Dir(hooksPath), 0o755); err != nil {
		t.Fatal(err)
	}
	body := `{
  "gentle-ai-sdd-agents-hardening": {
    "PreInvocation": [
      {"type":"command","command":"printf empty"}
    ]
  },
  "unrelated": {
    "PreInvocation": [
      {"type":"command","command":"sdd-explore sdd-propose sdd-apply sdd-verify review-* jd-judge-* jd-fix-agent fail closed CodeGraph"}
    ]
  }
}`
	if err := os.WriteFile(hooksPath, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if HasAntigravitySddAgentsHardeningContract(home) {
		t.Fatal("HasAntigravitySddAgentsHardeningContract() = true when required phrases only appear outside managed block")
	}
}

func TestHasAntigravitySddAgentsHardeningContract(t *testing.T) {
	home := t.TempDir()
	if HasAntigravitySddAgentsHardeningContract(home) {
		t.Fatal("HasAntigravitySddAgentsHardeningContract() = true before install, want false")
	}

	if _, _, err := installAntigravitySddAgentsPlugin(home); err != nil {
		t.Fatalf("installAntigravitySddAgentsPlugin() error = %v", err)
	}
	if !HasAntigravitySddAgentsHardeningContract(home) {
		t.Fatal("HasAntigravitySddAgentsHardeningContract() = false after install, want true")
	}
}

func TestAntigravitySddAgentsPluginDoesNotAffectOtherAgents(t *testing.T) {
	// Sanity: installAntigravitySddAgentsPlugin must NEVER touch paths
	// outside ~/.gemini/antigravity-cli/plugins/gentle-ai-sdd-agents.
	home := t.TempDir()
	_, files, err := installAntigravitySddAgentsPlugin(home)
	if err != nil {
		t.Fatalf("installAntigravitySddAgentsPlugin() error = %v", err)
	}
	for _, f := range files {
		clean := filepath.Clean(f)
		prefix := filepath.Clean(filepath.Join(home, ".gemini", "antigravity-cli", "plugins", "gentle-ai-sdd-agents")) + string(os.PathSeparator)
		if !strings.HasPrefix(clean, prefix) {
			t.Fatalf("installAntigravitySddAgentsPlugin() wrote %q outside the gentle-ai-sdd-agents plugin dir", clean)
		}
		if strings.Contains(clean, "settings.json") {
			t.Fatalf("installAntigravitySddAgentsPlugin() must not touch user settings.json; got %q", clean)
		}
		if strings.Contains(clean, "mcp_config.json") {
			t.Fatalf("installAntigravitySddAgentsPlugin() must not touch user mcp_config.json; got %q", clean)
		}
	}
}

func TestAntigravitySddAgentsPluginScopingToCLIOnly(t *testing.T) {
	// The antigravity-desktop variant does not currently consume
	// plugin/hooks.json; we deliberately do NOT install gentle-ai-sdd-agents
	// under ~/.gemini/antigravity-desktop. This test locks that scope.
	home := t.TempDir()
	_, _, err := installAntigravitySddAgentsPlugin(home)
	if err != nil {
		t.Fatalf("installAntigravitySddAgentsPlugin() error = %v", err)
	}
	desktopPlugin := filepath.Join(home, ".gemini", "antigravity-desktop", "plugins", "gentle-ai-sdd-agents", "plugin.json")
	if _, err := os.Stat(desktopPlugin); !os.IsNotExist(err) {
		t.Fatalf("installAntigravitySddAgentsPlugin() must NOT write to antigravity-desktop; stat err = %v", err)
	}
}

func TestInjectAntigravityInstallsSddAgentsHardeningPlugin(t *testing.T) {
	home := t.TempDir()
	adapter, err := agents.NewAdapter(model.AgentAntigravity)
	if err != nil {
		t.Fatalf("NewAdapter(antigravity) error = %v", err)
	}

	result, err := Inject(home, adapter, "")
	if err != nil {
		t.Fatalf("Inject(antigravity) error = %v", err)
	}
	if !result.Changed {
		t.Fatal("Inject(antigravity) changed = false, want true")
	}

	pluginPath := filepath.Join(home, ".gemini", "antigravity-cli", "plugins", "gentle-ai-sdd-agents", "plugin.json")
	hooksPath := filepath.Join(home, ".gemini", "antigravity-cli", "plugins", "gentle-ai-sdd-agents", "hooks.json")
	for _, want := range []string{pluginPath, hooksPath} {
		if !containsString(result.Files, want) {
			t.Fatalf("Inject(antigravity).Files missing %q; got %v", want, result.Files)
		}
		if _, err := os.Stat(want); err != nil {
			t.Fatalf("expected %q to exist after Inject: %v", want, err)
		}
	}

	if !HasAntigravitySddAgentsHardeningContract(home) {
		t.Fatal("HasAntigravitySddAgentsHardeningContract() = false after Inject, want true")
	}

	// Sanity: the hooks.json must be valid JSON with the expected schema.
	hooksBody, err := os.ReadFile(hooksPath)
	if err != nil {
		t.Fatalf("ReadFile(hooks.json) error = %v", err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(hooksBody, &parsed); err != nil {
		t.Fatalf("hooks.json is not valid JSON: %v\n%s", err, hooksBody)
	}
	block, ok := parsed["gentle-ai-sdd-agents-hardening"].(map[string]any)
	if !ok {
		t.Fatalf("hooks.json missing gentle-ai-sdd-agents-hardening block: %s", hooksBody)
	}
	pre, ok := block["PreInvocation"].([]any)
	if !ok || len(pre) == 0 {
		t.Fatalf("hooks.json missing PreInvocation list: %s", hooksBody)
	}
}

func TestInjectDoesNotInstallSddAgentsPluginForNonAntigravity(t *testing.T) {
	// Lock the scope: only Antigravity gets the gentle-ai-sdd-agents plugin.
	// Other agents must not see plugin/hooks.json appear in their
	// injection result.
	cases := []struct {
		name    string
		agentID string
		invoke  func(home string) error
	}{
		{
			name:    "opencode",
			agentID: "opencode",
			invoke: func(home string) error {
				adapter, err := agents.NewAdapter("opencode")
				if err != nil {
					return err
				}
				_, err = Inject(home, adapter, "single")
				return err
			},
		},
		{
			name:    "claude",
			agentID: "claude-code",
			invoke: func(home string) error {
				adapter, err := agents.NewAdapter("claude-code")
				if err != nil {
					return err
				}
				_, err = Inject(home, adapter, "")
				return err
			},
		},
		{
			name:    "gemini-cli",
			agentID: "gemini-cli",
			invoke: func(home string) error {
				adapter, err := agents.NewAdapter("gemini-cli")
				if err != nil {
					return err
				}
				_, err = Inject(home, adapter, "")
				return err
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			if err := tc.invoke(home); err != nil {
				t.Fatalf("Inject(%s) error = %v", tc.agentID, err)
			}
			// The antigravity-cli plugin path must NOT be created for other
			// agents — it is Antigravity-only.
			pluginPath := filepath.Join(home, ".gemini", "antigravity-cli", "plugins", "gentle-ai-sdd-agents", "plugin.json")
			if _, err := os.Stat(pluginPath); !os.IsNotExist(err) {
				t.Fatalf("Inject(%s) wrote Antigravity-only plugin path %q; stat err = %v", tc.agentID, pluginPath, err)
			}
		})
	}
}
