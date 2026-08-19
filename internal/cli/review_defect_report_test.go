package cli

import (
	"strings"
	"testing"
)

// TestReviewDefectReportRendersTemplateShapedFields is the RED-first proof
// for organic-dx Phase 5 task 5.2: reviewDefectReport renders Markdown whose
// section headers match .github/ISSUE_TEMPLATE/bug_report.yml's field labels,
// populated from gentle-ai version + commit, OS/arch, the failing operation
// shape (not raw argv), the verbatim reason code, opaque state identifiers,
// and the stop-code row's terminal precondition.
func TestReviewDefectReportRendersTemplateShapedFields(t *testing.T) {
	report := newReviewDefectReport(reviewDefectReportInput{
		Operation:            "review capture-evidence --lineage --target --expected-revision --outcome --input",
		ReasonCode:           "captured_final_evidence_conflict",
		ErrorMessage:         "captured final evidence already exists with different bytes",
		TerminalPrecondition: "final verification evidence was already captured for this validating target with different bytes",
		StateIdentifiers: map[string]string{
			"state":  "validating",
			"target": "sha256:" + strings.Repeat("a", 64),
		},
	})
	report.Version, report.Commit, report.OS, report.Arch = "1.2.3", "deadbeef", "linux", "amd64"
	rendered := report.render()

	for _, header := range []string{
		"# Bug Description", "## Steps to Reproduce", "## Expected Behavior", "## Actual Behavior",
		"## Gentle AI Version", "## Operating System", "## AI Agent / Client", "## Affected Area", "## Logs / Error Output",
	} {
		if !strings.Contains(rendered, header) {
			t.Errorf("rendered report missing template-shaped header %q:\n%s", header, rendered)
		}
	}
	for _, want := range []string{
		"captured_final_evidence_conflict",
		"captured final evidence already exists with different bytes",
		"final verification evidence was already captured for this validating target with different bytes",
		"1.2.3", "deadbeef", "linux/amd64",
		"review capture-evidence --lineage --target --expected-revision --outcome --input",
		"state: validating",
		"target: sha256:" + strings.Repeat("a", 64),
	} {
		if !strings.Contains(rendered, want) {
			t.Errorf("rendered report missing expected content %q:\n%s", want, rendered)
		}
	}
}

// TestReviewDefectReportScrubsPoisonedInput is the RED-first, hard-privacy
// requirement of task 5.3: driving the builder with content-bearing strings,
// diffs, file contents, absolute paths under $HOME, env-var-shaped values,
// and email addresses must leave none of that raw material in the rendered
// output. A leak here would land on a public issue tracker.
func TestReviewDefectReportScrubsPoisonedInput(t *testing.T) {
	home := "/home/definitely-a-real-user"
	poisoned := reviewDefectReportInput{
		Operation: "review capture-evidence --input " + home + "/secret-project/plan.txt",
		ReasonCode: "captured_final_evidence_conflict\n" +
			"-func leaked() { return \"" + home + "/.ssh/id_rsa\" }",
		ErrorMessage: "open " + home + "/.ssh/id_rsa: permission denied\n" +
			"@@ -1,3 +1,4 @@\n+AWS_SECRET_ACCESS_KEY=super-secret-value-should-never-leak\n" +
			"contact definitely-a-real-user@example.com for access",
		TerminalPrecondition: "diff:\n--- a/file\n+++ b/file\n@@\n-old line with secret token XYZZY-DO-NOT-LEAK\n+new line",
		StateIdentifiers: map[string]string{
			"path":    home + "/.config/gentle-ai/credentials.json",
			"env":     "GITHUB_TOKEN=ghp_shouldneverleakshouldneverleakshouldneverleak",
			"email":   "definitely-a-real-user@example.com",
			"content": "the quick brown fox file contents should never appear\nsecond line of file content",
		},
	}
	rendered := newReviewDefectReport(poisoned).render()

	poisonedSubstrings := []string{
		home,
		"id_rsa",
		"AWS_SECRET_ACCESS_KEY=super-secret-value-should-never-leak",
		"definitely-a-real-user@example.com",
		"XYZZY-DO-NOT-LEAK",
		"ghp_shouldneverleakshouldneverleakshouldneverleak",
		"the quick brown fox file contents should never appear",
		"second line of file content",
		"credentials.json",
	}
	for _, poison := range poisonedSubstrings {
		if strings.Contains(rendered, poison) {
			t.Fatalf("rendered defect report leaked poisoned input %q:\n%s", poison, rendered)
		}
	}
	// A multi-line reason code must never smuggle a second line (e.g. a diff
	// hunk) into the report; only its first, scrubbed line may survive.
	if strings.Contains(rendered, "-func leaked()") {
		t.Fatalf("rendered defect report leaked a multi-line reason code payload:\n%s", rendered)
	}
}

// TestReviewDefectReportScrubbedFieldIsSingleLineAndBounded pins the exact
// scrubbing contract reviewScrubDefectReportField implements, independent of
// the full-report integration test above.
func TestReviewDefectReportScrubbedFieldIsSingleLineAndBounded(t *testing.T) {
	cases := map[string]struct{ input, wantNotContain string }{
		"multiline":     {"first line\nsecond line", "second line"},
		"absolute-path": {"see /home/someone/secret.txt for detail", "/home/someone/secret.txt"},
		"email":         {"reach out to someone@example.com", "someone@example.com"},
		"env-shaped":    {"TOKEN=abc123shouldnotleak", "abc123shouldnotleak"},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			scrubbed := reviewScrubDefectReportField(testCase.input)
			if strings.Contains(scrubbed, testCase.wantNotContain) {
				t.Fatalf("scrubbed field still contains %q: %q", testCase.wantNotContain, scrubbed)
			}
		})
	}
}

// TestReviewDefectReportScrubberPreservesPublicContractsWhileRedactingSecretsAndPaths
// asserts that public contract identifiers and known public contract environment
// assignments survive scrubbing intact, while absolute paths, email addresses,
// and sensitive environment assignments remain cleanly redacted (#3443).
func TestReviewDefectReportScrubberPreservesPublicContractsWhileRedactingSecretsAndPaths(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "contract-integration-v1",
			input: "gentle-ai.review-integration/v1",
			want:  "gentle-ai.review-integration/v1",
		},
		{
			name:  "contract-integration-v2",
			input: "gentle-ai.review-integration/v2",
			want:  "gentle-ai.review-integration/v2",
		},
		{
			name:  "contract-pi-relay-v1",
			input: "gentle-pi.review-relay/v1",
			want:  "gentle-pi.review-relay/v1",
		},
		{
			name:  "contract-failure-v2",
			input: "gentle-ai.review-failure/v2",
			want:  "gentle-ai.review-failure/v2",
		},
		{
			name:  "contract-sdd-binding-v1",
			input: "gentle-ai.sdd-review-binding/v1",
			want:  "gentle-ai.sdd-review-binding/v1",
		},
		{
			name:  "command-with-contract",
			input: "gentle-ai review status --contract gentle-ai.review-integration/v2 --next-transition",
			want:  "gentle-ai review status --contract gentle-ai.review-integration/v2 --next-transition",
		},
		{
			name:  "public-contract-env-assignment",
			input: "export GENTLE_PI_REVIEW_RELAY_CONTRACT=gentle-pi.review-relay/v1 and re-run",
			want:  "export GENTLE_PI_REVIEW_RELAY_CONTRACT=gentle-pi.review-relay/v1 and re-run",
		},
		{
			name:  "posix-absolute-path",
			input: "see /home/someone/secret.txt for detail",
			want:  "see <redacted> for detail",
		},
		{
			name:  "windows-backslash-absolute-path",
			input: `path C:\Users\someone\secret.txt is sensitive`,
			want:  "path <redacted> is sensitive",
		},
		{
			name:  "windows-slash-absolute-path",
			input: "path C:/Users/someone/secret.txt is sensitive",
			want:  "path <redacted> is sensitive",
		},
		{
			name:  "windows-root-path",
			input: `root \secret\path is sensitive`,
			want:  "root <redacted> is sensitive",
		},
		{
			name:  "secret-token-env",
			input: "TOKEN=supersecret123",
			want:  "<redacted>",
		},
		{
			name:  "aws-secret-env",
			input: "AWS_SECRET_ACCESS_KEY=supersecretkey",
			want:  "<redacted>",
		},
		{
			name:  "github-token-env",
			input: "GITHUB_TOKEN=ghp_secrettoken123",
			want:  "<redacted>",
		},
		{
			name:  "email-address",
			input: "contact definitely-a-real-user@example.com for access",
			want:  "contact <redacted> for access",
		},
		{
			name:  "multiline-truncated",
			input: "first line\nsecond line",
			want:  "first line",
		},
		{
			name:  "colon-delimited-posix-path",
			input: "open:/home/someone/secret.txt: permission denied",
			want:  "open:<redacted>: permission denied",
		},
		{
			name:  "url-preservation",
			input: "see http://example.com/api/v1 and https://gentle.ai/docs and file:///var/log",
			want:  "see http://example.com/api/v1 and https://gentle.ai/docs and file:///var/log",
		},
		{
			name:  "public-contract-env-semicolon",
			input: "export GENTLE_PI_REVIEW_RELAY_CONTRACT=gentle-pi.review-relay/v1;",
			want:  "export GENTLE_PI_REVIEW_RELAY_CONTRACT=gentle-pi.review-relay/v1;",
		},
		{
			name:  "public-contract-env-parens-period",
			input: "(see GENTLE_PI_REVIEW_RELAY_CONTRACT=gentle-pi.review-relay/v1).",
			want:  "(see GENTLE_PI_REVIEW_RELAY_CONTRACT=gentle-pi.review-relay/v1).",
		},
		{
			name:  "sensitive-env-semicolon",
			input: "AWS_SECRET_ACCESS_KEY=supersecretkey;",
			want:  "<redacted>;",
		},
		{
			name:  "curly-brace-delimited-path",
			input: "path {/var/log/secret.log} enclosed",
			want:  "path {<redacted>} enclosed",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := reviewScrubDefectReportField(tc.input)
			if got != tc.want {
				t.Fatalf("reviewScrubDefectReportField(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
