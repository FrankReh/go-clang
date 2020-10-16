package main

import (
	"flag"
	"strings"
	"testing"
)

func TestGoClangDump(t *testing.T) {
	additional := dropEmpties(strings.Split(*cflags, " "))
	for _, args := range [][]string{
		[]string{"-c", "../../testdata/hello.c"},
		[]string{"-c", "../../testdata/globals.c"},
		[]string{"-ref", "-c", "../../testdata/globals.c"},
	} {
		args = append(additional, args...)
		r := cmd(args)
		if r != 0 {
			t.Errorf("cmd(%v) = %d", args, r)
		}
	}
}

var cflags = flag.String("cflags", "", "space separated flags to pass to clang")

func dropEmpties(ss []string) (r []string) {
	for _, s := range ss {
		if s != "" {
			r = append(r, s)
		}
	}
	return
}
