package clang

import (
	"testing"
)

func TestTargetInfo(t *testing.T) {
	idx := NewIndex(0, 0)
	defer idx.Dispose()

	tu := idx.ParseTranslationUnit("cursor.c", nil, nil, 0)
	assertTrue(t, tu.IsValid())
	defer tu.Dispose()

	targetinfo := tu.TargetInfo()
	assertTrue(t, targetinfo.Triple != "")
	if !(targetinfo.PointerWidth == 32 || targetinfo.PointerWidth == 64) {
		t.Fatalf("PointerWidth should be 32 or 64, not %d", targetinfo.PointerWidth)
	}
}
