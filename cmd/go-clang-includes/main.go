// go-clang-includes lists the include files.
//
// This shows two things: how to find the included files and perhaps more
// importantly, how the source filename does not have to be explicitly passed
// to the ParseTranslationUnit method. Libclang finds the source filename from
// the arguments, just as clang would.
//
// So this command takes no flags. All the arguments are passed to
// libclang.
//
// $ go-clang-includes -c -I/include/dir -DDEBUG ../../testdata/hello.c
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

	// Pass the DetailedPreprocessingRecord option so libclang keeps the detailed information about inclusions.
	tu := idx.ParseTranslationUnit("", args, nil, clang.TranslationUnit_DetailedPreprocessingRecord)
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

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case cursorkind.InclusionDirective:
			file := cursor.IncludedFile()
			//uniqueID, _ := file.UniqueID()
			//fmt.Printf("%s: %s\n", cursor.Kind().Spelling(), cursor.Spelling())
			//fmt.Printf("  name: %s, time: %v, uniqueID: %d\n", file.Name(), file.Time(), uniqueID)
			fmt.Printf("%s\n", file.Name())
		}

		// No need to recurse. The inclusion directives will be found at the top level.
		return clang.ChildVisit_Continue
	})

	return 0
}
