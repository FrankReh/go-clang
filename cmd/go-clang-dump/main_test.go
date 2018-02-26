package main

import (
	"testing"
)

func TestGoClangDump(t *testing.T) {
	for _, fname := range []string{
		"../../testdata/basicparsing.c",
	} {
		args := []string{"-fname", fname}
		r := cmd(args)
		if r != 0 {
			t.Errorf("cmd(%v) = %d", args, r)
		}
	}
}
