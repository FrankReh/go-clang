package clang_test

import (
	"testing"
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
