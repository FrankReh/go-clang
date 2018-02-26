package main

import (
	"testing"
)

func TestGoClangCompDB(t *testing.T) {
	for _, path := range []string{
		"../../testdata",
	} {
		r := cmd([]string{path})
		if r != 0 {
			t.Errorf("cmd([]string{%s}) = %d", path, r)
		}
	}
}
