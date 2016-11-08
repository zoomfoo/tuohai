package convert

import (
	"testing"
)

func TestRuneAccumulation(t *testing.T) {
	t.Log(RuneAccumulation("b56486af2d149c3b816d593bf0e4f1b5"))
	t.Log(RuneAccumulation("202cb962ac59075b964b07152d234b71"))
}
