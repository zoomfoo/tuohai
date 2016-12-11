package convert

import (
	"testing"
)

func TestRuneAccumulation(t *testing.T) {
	t.Log(RuneAccumulation("84558b0cf90a4166") % 4)
}
