package openapi

import (
	"testing"

	"github.com/nasa9084/tracedeq"
)

func unexpected(t *testing.T, name string, got, want interface{}) {
	t.Helper()

	t.Errorf("unexpected %s:\n  got:  %v\n  want: %v", name, got, want)
}

func assertSameError(t *testing.T, got, want error) {
	t.Helper()

	if want == nil && got == nil {
		return
	}
	if got.Error() != want.Error() {
		unexpected(t, "error", got.Error(), want.Error())
		return
	}
}

func assertEqual(t *testing.T, got, want interface{}) {
	t.Helper()

	if result := tracedeq.DeepEqual(got, want); !result.IsEqual {
		unexpected(t, result.Trace.Join("."), result.X, result.Y)
	}
}
