// go-clang-globals lists the global functions and variables or the external references made.
//
// This command takes optionally one flag (-ref). There is no help for it though.
//
// $ go-clang-globals -c ../../testdata/hello.c
// or
// $ go-clang-globals -ref -c ../../testdata/hello.c
package main

import (
	"fmt"
	"os"

	"github.com/frankreh/go-clang/clang"
	"github.com/frankreh/go-clang/clang/cursorkind"
)

func main() {
	os.Exit(cmd(os.Args[1:]))
}

func cmd(args []string) int {
	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()

	// Allow for first argument to be "-ref" which indicates to look for external references.
	references := false
	options := clang.TranslationUnit_DetailedPreprocessingRecord
	if len(args) > 0 && args[0] == "-ref" {
		references = true
		args = args[1:]
		options = 0 // Don't look into include files (not sure about this)
	}

	// Pass the DetailedPreprocessingRecord option so libclang keeps the detailed information about inclusions.
	tu := idx.ParseTranslationUnit("", args, nil, options)
	defer tu.Dispose()

	// The tu.Spelling() shows the source filename that was extracted from the args.
	//fmt.Printf("tu: %s\n", tu.Spelling())

	diagnostics := tu.Diagnostics()
	for _, d := range diagnostics {
		fmt.Println("PROBLEM:", d.Spelling())
	}

	cursor := tu.TranslationUnitCursor()
	if cursor.IsNull() {
		fmt.Println("PROBLEM: TranslationUnitCursor creation failed")
		return 1
	}

	if references {
		seen := make(map[string]bool)

		cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
			switch cursor.Kind() {
			case cursorkind.DeclRefExpr:
				if cursor.Definition().IsNull() {
					name := cursor.Spelling()
					if !seen[name] {
						fmt.Printf("%s\n", name)
						seen[name] = true
					}
				}
			}

			return clang.ChildVisit_Recurse
		})
	} else {
		cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
			switch cursor.Kind() {
			// cursor must be FunctionDecl or VarDecl
			// cursor must be a definition itself
			// cursor storage class must be none
			case cursorkind.FunctionDecl, cursorkind.VarDecl:
				if cursor.Equal(cursor.Definition()) {
					if cursor.StorageClass() == clang.SC_None {
						fmt.Printf("%s: %s (%s)\n", clang.CursorKindSpelling(cursor.Kind()), cursor.Spelling(), cursor.USR())
					}
				}
			}

			// No need to recurse. The globals will be found at the top level.
			return clang.ChildVisit_Continue
		})
	}

	return 0
}
