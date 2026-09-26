package agentguidance_test

import (
	"os"
	"testing"

	"github.com/gentleman-programming/gentle-ai/v3/internal/components/agentguidance"
	"github.com/gentleman-programming/gentle-ai/v3/internal/components/reviewassets"
)

// TestMain wires the production review contract source the installer
// registers, so rendered prompts in this package match what users receive.
func TestMain(m *testing.M) {
	agentguidance.SetReviewContractSource(reviewassets.ReviewExecutionContractFor)
	os.Exit(m.Run())
}
