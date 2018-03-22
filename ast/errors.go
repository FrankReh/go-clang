package ast

import (
	"fmt"
	"strings"
)

func (tu *TranslationUnit) Err(elist ...interface{}) *TuErr {
	err := &TuErr{tu: tu}

	return err.Errors(elist...)
}

type ErrFlag int

const (
	ChildrenCountErr ErrFlag = 1 << iota
)

// Struct TuErr contains accumulated error parts.
type TuErr struct {
	tu       *TranslationUnit
	errflags ErrFlag
	elist    []interface{}

	cursorIds []int
	childIds  []int
	msgs      []string
}

func (t TuErr) Error() string {
	if t.tu == nil {
		return ""
	}
	b := new(strings.Builder)
	fmt.Fprintf(b, "TranslationUnit Error:\n")

	if len(t.cursorIds) > 0 {
		if t.tu == nil {
			fmt.Fprintf(b, " cursorIds %v\n", t.cursorIds)
		} else {
			for _, cursorId := range t.cursorIds {
				if cursorId >= 0 && cursorId < len(t.tu.Cursors) {
					fmt.Fprintf(b, " cursorId %d %#v\n", cursorId, t.tu.Cursors[cursorId])
				} else {
					fmt.Fprintf(b, " cursorId %d\n", cursorId)
				}
			}
		}
	}
	if len(t.childIds) > 0 {
		if t.tu == nil {
			fmt.Fprintf(b, " childIds %v\n", t.childIds)
		} else {
			for _, childId := range t.childIds {
				if childId >= 0 && childId < len(t.tu.Cursors) {
					fmt.Fprintf(b, " childId %d %#v\n", childId, t.tu.Cursors[childId])
				} else {
					fmt.Fprintf(b, " childId %d\n", childId)
				}
			}
		}
	}
	if len(t.elist) > 0 {
		fmt.Fprintf(b, " %v\n", t.elist)
	}
	if len(t.msgs) > 0 {
		fmt.Fprintf(b, " %v\n", t.msgs)
	}
	if t.errflags&ChildrenCountErr != 0 {
		fmt.Fprintf(b, " incorrect children count\n")
	}

	return b.String()
}

func (t *TuErr) Errors(elist ...interface{}) *TuErr {
	for _, e := range elist {
		switch v := e.(type) {
		case *TuErr:
			if t.tu != v.tu {
				panic("Incompatible TuErr passed to TuErr")
			}
			t.elist = append(t.elist, v.elist...)
			t.errflags |= v.errflags
			t.cursorIds = appendUniquely(t.cursorIds, v.cursorIds...)
			t.childIds = appendUniquely(t.childIds, v.childIds...)
			t.msgs = append(t.msgs, v.msgs...)

		case *CursorErr:
			t.Cursor(v.cursor) // Add to our list of cursors
		case CursorErr:
			t.Cursor(v.cursor) // Add to our list of cursors

		case *ChildErr:
			t.Children(v.child) // Add to our list of children
		case ChildErr:
			t.Children(v.child) // Add to our list of children

		case string:
			t.Msg(v)
		default:
			fmt.Println("default")
			t.elist = append(t.elist, e)
		}
	}
	return t
}

func (t *TuErr) Cursor(cursorId int) *TuErr {
	t.cursorIds = appendUniquely(t.cursorIds, cursorId)
	return t
}

func (t *TuErr) Children(ids ...int) *TuErr {
	t.childIds = appendUniquely(t.childIds, ids...)
	return t
}

func (t *TuErr) ChildrenCount() *TuErr {
	t.errflags |= ChildrenCountErr
	return t
}

func (t *TuErr) Msg(msg string) *TuErr {
	t.msgs = append(t.msgs, msg)
	return t
}

// appendUniquely returns list with value appended if not already in it.
func appendUniquely(list []int, values ...int) []int {
	for _, v := range values {
		if isMember(v, list) {
			continue
		}
		list = append(list, v)
	}
	return list
}

// isMember return true if value is found in list.
func isMember(value int, list []int) bool {
	for _, el := range list {
		if value == el {
			return true
		}
	}
	return false
}

// Struct OutOfRangeErr is an error that reports the index and the length
// that was exceeded.
type OutOfRangeErr struct {
	index  int
	length int
}

func (o OutOfRangeErr) Error() string {
	return fmt.Sprintf("index:%d length:%d", o.index, o.length)
}

// Struct CursorErr is an error that reports the cursor index.
type CursorErr struct {
	cursor int
}

func (o CursorErr) Error() string {
	return fmt.Sprintf("cursor:%d", o.cursor)
}

// Struct ChildErr is an error that reports the Child index.
type ChildErr struct {
	child int
}

func (o ChildErr) Error() string {
	return fmt.Sprintf("child:%d", o.child)
}

// Struct NotFoundErr is an error that reports the Child index.
type NotFoundErr struct {
	name      string
	container string
}

func (o NotFoundErr) Error() string {
	return fmt.Sprintf("NotFound:%s in %s", o.name, o.container)
}
