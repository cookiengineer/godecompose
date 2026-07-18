package decompile_crypto

import (
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestCrypto(t *testing.T) {
	b := decompile.CompileAndOpen(t, "crypto")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "crypto")

	for _, m := range r.Matches {
		if decompile.HasMatch([]matcher.Match{m}, "sha256.") || decompile.HasMatch([]matcher.Match{m}, "md5.") ||
			decompile.HasMatch([]matcher.Match{m}, "aes.") || decompile.HasMatch([]matcher.Match{m}, "hmac.") ||
			decompile.HasMatch([]matcher.Match{m}, "rand.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}
