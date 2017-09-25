package clang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTargetInfo(t *testing.T) {
	idx := NewIndex(0, 0)
	defer idx.Dispose()

	tu := idx.ParseTranslationUnit("cursor.c", nil, nil, 0)
	assert.True(t, tu.IsValid())
	defer tu.Dispose()

	targetinfo := tu.TargetInfo()
	assert.NotEmpty(t, targetinfo.Triple)
	assert.True(t, targetinfo.PointerWidth == 32 || targetinfo.PointerWidth == 64,
		"PointerWidth should be 32 or 64, not %d", targetinfo.PointerWidth)
}
