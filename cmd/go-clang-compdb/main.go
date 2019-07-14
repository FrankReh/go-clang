// go-clang-compdb dumps the content of a clang compilation database
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/frankreh/go-clang/clang"
)

func main() {
	os.Exit(cmd(os.Args[1:]))
}

func cmd(args []string) int {
	if len(args) == 0 {
		fmt.Printf("**error: you need to give a directory containing a 'compile_commands.json' file\n")

		return 1
	}

	dir := os.ExpandEnv(args[0])
	fmt.Printf(":: inspecting [%s]...\n", dir)

	fname := filepath.Join(dir, "compile_commands.json")
	f, err := os.Open(fname)
	if err != nil {
		fmt.Printf("**error: could not open file [%s]: %v\n", fname, err)

		return 1
	}
	f.Close()

	db, err := clang.FromDirectory(dir)
	if err != nil {
		fmt.Printf("**error: could not open compilation database at [%s]: %v\n", dir, err)

		os.Exit(1)
	}
	defer db.Dispose()

	cmds := db.AllCompileCommands()

	// To skip all the rest, just use this one liner.
	// fmt.Printf(":: cmds = %q\n", cmds)

	fmt.Printf(":: got %d compile commands\n", len(cmds))

	for i, cmd := range cmds {

		fmt.Printf("::  --- cmd=%d ---\n", i)
		fmt.Printf("::  dir  = %q\n", cmd.Directory)
		fmt.Printf("::  file = %q\n", cmd.Filename)

		fmt.Printf("::  nargs= %d\n", len(cmd.Args))

		fmt.Printf("::  args =%q\n", cmd.Args)
		if i+1 != len(cmds) {
			fmt.Printf("::\n")
		}
	}
	fmt.Printf(":: inspecting [%s]... [done]\n", dir)

	return 0
}
