package clang_test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"

	"github.com/frankreh/go-clang/ast"
)

// first_gobtest uses native gob coder and decoder on entire TranslationUnit
// and calls DecodeFinish() explicitly.
func first_gobtest(t *testing.T, testname string, tu *ast.TranslationUnit) int {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf) // Will write to buf.
	dec := gob.NewDecoder(&buf) // Will read from buf.

	// Encode (send) some values.
	err := enc.Encode(tu)
	if err != nil {
		t.Fatalf("encode error: %s %s", testname, err)
	}
	encodeLen := buf.Len()

	// Decode (receive) and print the values.
	var tu2 ast.TranslationUnit
	err = dec.Decode(&tu2)
	if err != nil {
		t.Fatalf("decode error: %s %s", testname, err)
	}

	// In this first gob encoding version, we have to call DecodeFinish() explicitly
	// because there are some non public parts of the object that need setting up.
	tu2.DecodeFinish()
	beforeString := fmt.Sprintf("%v\n", tu)
	afterString := fmt.Sprintf("%v\n", &tu2)

	beforeGoString := fmt.Sprintf("%#v\n", tu)
	afterGoString := fmt.Sprintf("%#v\n", &tu2)

	if eq_err := tu.AssertEqual(&tu2); eq_err != nil {
		t.Errorf("%s: not equal: %s\n=== beforeGoString ===\n%s\n=== afterGoString ===\n%s\n",
			testname, eq_err, beforeString, afterString)
	}

	switch {
	case beforeGoString != afterGoString:
		t.Errorf("%s:\n=== beforeGoString ===\n%s\n=== afterGoString ===\n%s\n",
			testname, beforeGoString, afterGoString)
	default:
		// TBD Printing for debugging.
		//fmt.Printf("debug %s:\n=== afterGoString ===\n%s\n",
		//	testname, afterGoString)
		//fmt.Printf("debug %s:\n=== afterString ===\n%s\n",
		//	testname, afterString)
	}
	return encodeLen
}

// second_gobtest uses TranslationUnit v1 gob coder and decoder on entire TranslationUnit.
func second_gobtest(t *testing.T, testname string, tu *ast.TranslationUnit) int {

	var buf bytes.Buffer

	// Encode (send) some values.
	err := tu.EncodeGobV1(&buf)
	if err != nil {
		t.Fatalf("encode error: %s %s", testname, err)
	}
	encodeLen := buf.Len()

	// Decode (receive) and print the values.
	var tu2 ast.TranslationUnit
	err = tu2.DecodeGobV1(&buf)
	if err != nil {
		t.Fatalf("decode error: %s %s", testname, err)
	}
	beforeString := fmt.Sprintf("%v\n", tu)
	afterString := fmt.Sprintf("%v\n", &tu2)

	beforeGoString := fmt.Sprintf("%#v\n", tu)
	afterGoString := fmt.Sprintf("%#v\n", &tu2)

	if eq_err := tu.AssertEqual(&tu2); eq_err != nil {
		t.Errorf("%s: not equal: %s\n=== beforeGoString ===\n%s\n=== afterGoString ===\n%s\n",
			testname, eq_err, beforeString, afterString)
	}

	switch {
	case beforeGoString != afterGoString:
		t.Errorf("%s:\n=== beforeGoString ===\n%s\n=== afterGoString ===\n%s\n",
			testname, beforeGoString, afterGoString)
	default:
		// TBD Printing for debugging.
		//fmt.Printf("debug %s:\n=== afterGoString ===\n%s\n",
		//	testname, afterGoString)
		//fmt.Printf("debug %s:\n=== afterString ===\n%s\n",
		//	testname, afterString)
	}
	return encodeLen
}
