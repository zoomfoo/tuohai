package convert

import (
	"testing"
)

func TestRuneAccumulation(t *testing.T) {
	t.Log(RuneAccumulation("c444cdf6675ff9c2"))
	t.Log(RuneAccumulation("ba869ffc0b9d11f9"))
}
