package clang_test

import (
	"testing"

	"github.com/frankreh/go-clang-v5.0/clang"
)

func assertTrue(t *testing.T, b bool) {
	t.Helper()
	if !b {
		t.Fatal("not true")
	}
}

func assertEqualString(t *testing.T, s1, s2 string) {
	t.Helper()
	if s1 != s2 {
		t.Fatalf("%s != %s", s1, s2)
	}
}

func assertEqualInt(t *testing.T, i1, i2 int) {
	t.Helper()
	if i1 != i2 {
		t.Fatalf("%d != %d", i1, i2)
	}
}

func assertStringEmpty(t *testing.T, s string) {
	t.Helper()
	if s != "" {
		t.Fatalf("string not empty: %s", s)
	}
}

func assertStringNotEmpty(t *testing.T, s string) {
	t.Helper()
	if s == "" {
		t.Fatalf("string is empty")
	}
}

func TestBasicParsing(t *testing.T) {
	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()

	tu := idx.ParseTranslationUnit("../testdata/basicparsing.c", nil, nil, 0)
	assertTrue(t, tu.IsValid())
	defer tu.Dispose()

	found := 0

	tu.TranslationUnitCursor().Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.IsNull() {
			return clang.ChildVisit_Continue
		}

		switch cursor.Kind() {
		case clang.Cursor_FunctionDecl:
			assertEqualString(t, "foo", cursor.Spelling())

			found++
		case clang.Cursor_ParmDecl:
			assertEqualString(t, "bar", cursor.Spelling())

			found++
		}

		return clang.ChildVisit_Recurse
	})

	assertEqualInt(t, 2, found)
}

func TestReparse(t *testing.T) {
	us := []clang.UnsavedFile{
		clang.NewUnsavedFile("hello.cpp", "int world();"),
	}

	idx := clang.NewIndex(0, 0)
	defer idx.Dispose()

	tu := idx.ParseTranslationUnit("hello.cpp", nil, us, 0)
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

	us[0] = clang.NewUnsavedFile("hello.cpp", "int world2();")
	tu.ReparseTranslationUnit(us, 0)

	ok = false
	tu.TranslationUnitCursor().Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if s := cursor.Spelling(); s == "world2" {
			ok = true

			return clang.ChildVisit_Break
		} else if s == "world" {
			t.Errorf("'world' should no longer be part of the translationunit, but it still is")
		}

		return clang.ChildVisit_Continue
	})
	if !ok {
		t.Error("Expected to find 'world2', but didn't")
	}
}
