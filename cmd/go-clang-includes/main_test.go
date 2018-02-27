package main

import (
	"testing"
)

func TestGoClangDump(t *testing.T) {
	for _, args := range [][]string{
		[]string{"-c", "../../testdata/hello.c"},
	} {
		r := cmd(args)
		if r != 0 {
			t.Errorf("cmd(%v) = %d", args, r)
		}
	}
}
