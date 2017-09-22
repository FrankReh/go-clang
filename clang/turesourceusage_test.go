package clang

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceUsage(t *testing.T) {
	us := []UnsavedFile{
		NewUnsavedFile("hello.cpp", "int world();"),
	}

	idx := NewIndex(0, 0)
	defer idx.Dispose()

	tu := idx.ParseTranslationUnit("hello.cpp", nil, us, 0)
	assert.True(t, tu.IsValid())
	defer tu.Dispose()

	ru0 := tu.TUResourceUsage()

	entries0 := ru0.Entries()
	entries0a_str := fmt.Sprintf("%v", entries0)

	ru0.Dispose() // This cleanup causes the entries0 slice to have garbage backing store

	us[0] = NewUnsavedFile("hello.cpp", "int world2();")
	tu.ReparseTranslationUnit(us, 0)
	assert.True(t, tu.IsValid())

	ru1 := tu.TUResourceUsage()
	defer ru1.Dispose()

	entries1 := ru1.Entries()
	second_string := fmt.Sprintf("%v", entries1)

	entries0b_str := fmt.Sprintf("%v", entries0)

	assert.NotEqual(t, entries0a_str, second_string) // should be different
	assert.Equal(t, entries0a_str, entries0b_str)    // should be the same
}
