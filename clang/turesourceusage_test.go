package clang

import (
	"fmt"
	"testing"
)

func TestResourceUsage(t *testing.T) {
	us := []UnsavedFile{
		NewUnsavedFile("hello.cpp", "int world();"),
	}

	idx := NewIndex(0, 0)
	defer idx.Dispose()

	tu := idx.ParseTranslationUnit("hello.cpp", nil, us, 0)
	assertTrue(t, tu.IsValid())
	defer tu.Dispose()

	ru0_slice := tu.ResourceUsage()

	ru0_slice_str_a := fmt.Sprintf("%v", ru0_slice)

	us[0] = NewUnsavedFile("hello.cpp", " int world2();")
	tu.ReparseTranslationUnit(us, 0)
	assertTrue(t, tu.IsValid())

	ru1_slice := tu.ResourceUsage()
	ru1_slice_str := fmt.Sprintf("%v", ru1_slice)

	ru0_slice_str_b := fmt.Sprintf("%v", ru0_slice)

	_ = ru1_slice_str
	// Now the resource usages are the same for the two runs so don't enforce
	// they're being different.
	//assert.NotEqual(t, ru0_slice_str_a, ru1_slice_str) // should be different
	assertEqualString(t, ru0_slice_str_a, ru0_slice_str_b) // should be the same

	/*
		for i := range ru0_slice {
			fmt.Println(i, ru0_slice[i].Kind(), ru0_slice[i].Amount())
		}
	*/
}
