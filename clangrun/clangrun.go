// Package clangrun provides some useful general token and cursor visiting
// routines for use  with the clang package. They could be part of the clang
// package but they don't use an C calls and as such appear to be one level higher
// than it.
//
// Two forms of execute are provided, one that is a low level requiring no other
// setup first. One more a level higher in functionality but requiring some fields
// of a structure be setup first.
//
// Execute()
//
// Callbacks.Execute()
//
package clangrun

import (
	"fmt"

	"github.com/frankreh/go-clang-v5.0/clang"
)

// Define four interfaces, matching the four forms of translation unit visiting
// we support.
type TokenVisiter interface {
	TokenVisit(tu clang.TranslationUnit, token clang.Token)
}
type TopCursorVisiter interface {
	TopCursorVisit(tu clang.TranslationUnit, cursor, parent clang.Cursor)
}

type FullCursorVisiter interface {
	FullCursorVisit(tu clang.TranslationUnit, cursor, parent clang.Cursor)
}

type TUParser interface {
	TUParse(tu *clang.TranslationUnit)
}

type Callbacks struct {
	Options clang.TranslationUnit_Flags
	HdrCode string
	SrcCode string

	// Four slices of callback functions that can be executed on behalf of the
	// translation unit.  One that is run on every token, one on every top
	// level cursor, one for each cursor, and one that is run for the entire
	// TransactionUnit - kind of the fallback.
	// They can be set one of two ways:
	// - set directly with an Append.
	// - set indirectly by layering an interface.
	TokenFn      []func(tu clang.TranslationUnit, token clang.Token)
	TopCursorFn  []func(tu clang.TranslationUnit, cursor, parent clang.Cursor)
	FullCursorFn []func(tu clang.TranslationUnit, cursor, parent clang.Cursor)
	TuParseFn    []func(tu *clang.TranslationUnit)
}

func (c *Callbacks) AppendTokenFn(f func(tu clang.TranslationUnit, token clang.Token)) {
	c.TokenFn = append(c.TokenFn, f)
}

func (c *Callbacks) AppendTopCursorFn(f func(tu clang.TranslationUnit, cursor, parent clang.Cursor)) {
	c.TopCursorFn = append(c.TopCursorFn, f)
}

func (c *Callbacks) AppendFullCursorFn(f func(tu clang.TranslationUnit, cursor, parent clang.Cursor)) {
	c.FullCursorFn = append(c.FullCursorFn, f)
}

func (c *Callbacks) AppendTuParseFn(f func(tu *clang.TranslationUnit)) {
	c.TuParseFn = append(c.TuParseFn, f)
}

func (c *Callbacks) Layer(o interface{}) error {
	atLeastOne := false

	if f, ok := o.(TokenVisiter); ok {
		atLeastOne = true
		c.AppendTokenFn(f.TokenVisit)
	}

	if f, ok := o.(TopCursorVisiter); ok {
		atLeastOne = true
		c.AppendTopCursorFn(f.TopCursorVisit)
	}

	if f, ok := o.(FullCursorVisiter); ok {
		atLeastOne = true
		c.AppendFullCursorFn(f.FullCursorVisit)
	}

	if f, ok := o.(TUParser); ok {
		atLeastOne = true
		c.AppendTuParseFn(f.TUParse)
	}
	if !atLeastOne {
		return fmt.Errorf("The empty interface represents none of the callbacks.")
	}
	return nil
}

func (c *Callbacks) LayerAndExecute(o interface{}) error {

	err := c.Layer(o)
	if err != nil {
		return err
	}

	return c.Execute()
}

// Execute creates a clang.TranslationUnit that calls up to three Callbacks of callbacks on it.
// Any of the three callback functions can be left nil.
func (c *Callbacks) Execute() error {
	err := Execute(c.Options, c.SrcCode, c.HdrCode, func(idx clang.Index, tu clang.TranslationUnit) error {
		atLeastOne := false

		for _, tokenFn := range c.TokenFn {
			atLeastOne = true
			sourceRange := tu.TranslationUnitCursor().Extent()
			tokens := tu.Tokenize(sourceRange)
			for _, token := range tokens {
				tokenFn(tu, token)
			}
		}
		for _, topCursorFn := range c.TopCursorFn {
			atLeastOne = true

			tu.TranslationUnitCursor().Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {

				topCursorFn(tu, cursor, parent)

				return clang.ChildVisit_Continue
			})
		}
		for _, fullCursorFn := range c.FullCursorFn {
			atLeastOne = true

			tu.TranslationUnitCursor().Visit(
				func(cursor, parent clang.Cursor) clang.ChildVisitResult {

					fullCursorFn(tu, cursor, parent)

					return clang.ChildVisit_Recurse
				})
		}
		for _, tuParseFn := range c.TuParseFn {
			atLeastOne = true

			tuParseFn(&tu)
		}
		if !atLeastOne {
			return fmt.Errorf("No callbacks had been set.")
		}

		return nil
	})
	return err
}

// Execute creates a clang.TranslationUnit and calls run with it.
func Execute(options clang.TranslationUnit_Flags, srcCode, hdrCode string,
	run func(idx clang.Index, tu clang.TranslationUnit) error) error {

	srcFilename := "sample.c"

	var buffers []clang.UnsavedFile

	if hdrCode != "" {
		hdrFilename := "hdr.h"

		buffers = append(buffers, clang.NewUnsavedFile(hdrFilename, hdrCode)) // 1. unsaved file for header
		srcCode = fmt.Sprintf("#include \"%s\"\n%s", hdrFilename, srcCode)    // 2. include header in source
	}

	buffers = append(buffers, clang.NewUnsavedFile(srcFilename, srcCode))

	idx := clang.NewIndex(0, 0)
	defer idx.Dispose()

	tu := idx.ParseTranslationUnit(srcFilename, nil, buffers, options)
	if !tu.IsValid() {
		return fmt.Errorf("clang TranslationUnit is not valid")
	}
	defer tu.Dispose()

	return run(idx, tu)
}
