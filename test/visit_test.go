package clang_test

import (
	"testing"

	"github.com/frankreh/go-clang/clang"
)

func TestVisitBuffer(t *testing.T) {
	tmpfilename := "sample.c"
	buffers := []clang.UnsavedFile{
		clang.NewUnsavedFile(tmpfilename, "int world();"),
	}

	idx := clang.NewIndex(0, 0)
	defer idx.Dispose()

	tu := idx.ParseTranslationUnit(tmpfilename, nil, buffers, 0)
	assertTrue(t, tu.IsValid())
	defer tu.Dispose()

	ok := false
	tu.TranslationUnitCursor().Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.Spelling() == "world" {
			ok = true

			return clang.ChildVisit_Break
		}

		return clang.ChildVisit_Continue
	})
	if !ok {
		t.Error("Expected to find 'world', but didn't")
	}
}
