package main

import (
	"flag"
	"strings"
	"testing"
)

var cflags = flag.String("cflags", "", "space separated flags to pass to clang")

func dropEmpties(ss []string) (r []string) {
	for _, s := range ss {
		if s != "" {
			r = append(r, s)
		}
	}
	return
}

func TestGoClangDump(t *testing.T) {
	additional := dropEmpties(strings.Split(*cflags, " "))
	for _, args := range [][]string{
		[]string{"-c", "../../testdata/hello.c"},
	} {
		args = append(additional, args...)
		r := cmd(args)
		if r != 0 {
			t.Errorf("cmd(%v) = %d", args, r)
		}
	}
}
