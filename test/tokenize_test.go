package clang_test

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/frankreh/go-clang-v5.0/clang"
)

// tuRun creates a clang.TranslationUnit and calls run with it.
func tuRun(t *testing.T, options clang.TranslationUnit_Flags, srcCode, hdrCode string,
	run func(idx clang.Index, tu clang.TranslationUnit)) {

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
	assertTrue(t, tu.IsValid())
	defer tu.Dispose()

	run(idx, tu)
}

type run3Stages struct {
	options clang.TranslationUnit_Flags
	hdrCode string
	srcCode string

	// Three functions that can be run on behalf of the translation unit.
	// One that is run on every token, one on every top level cursor, and one for each cursor.
	// They can be set one of two ways:
	// - set directly by the owner of the instance - this was implemented first.
	// - set indirectly by passing a conforming interface (to run2)- this was the second way implemented.
	tokenFn      func(tu clang.TranslationUnit, token clang.Token)
	topCursorFn  func(tu clang.TranslationUnit, cursor, parent clang.Cursor)
	fullCursorFn func(tu clang.TranslationUnit, cursor, parent clang.Cursor, depth int, pre bool)
	tuParseFn    func(tu *clang.TranslationUnit)
}

func (input *run3Stages) run2(t *testing.T, o interface{}) {
	t.Helper()
	atLeastOne := false

	if f, ok := o.(tokenVisiter); ok {
		if input.tokenFn != nil {
			t.Fatal("bug: tokenFn had already been set")
		}
		atLeastOne = true
		input.tokenFn = f.tokenVisit
	}

	if f, ok := o.(topCursorVisiter); ok {
		if input.topCursorFn != nil {
			t.Fatal("bug: topCursorFn had already been set")
		}
		atLeastOne = true
		input.topCursorFn = f.topCursorVisit
	}

	if f, ok := o.(fullCursorVisiter); ok {
		if input.fullCursorFn != nil {
			t.Fatal("bug: fullCursorFn had already been set")
		}
		atLeastOne = true
		input.fullCursorFn = f.fullCursorVisit
	}

	if f, ok := o.(TUParser); ok {
		if input.tuParseFn != nil {
			t.Fatal("bug: tuParseFn had already been set")
		}
		atLeastOne = true
		input.tuParseFn = f.TUParse
	}

	if !atLeastOne {
		t.Fatal("arg o doesn't implement any of the optional callbacks")
	}

	input.run(t)
}

// run creates a clang.TranslationUnit that calls up to three stages of callbacks on it.
// Any of the three callback functions can be left nil.
func (input *run3Stages) run(t *testing.T) {
	tuRun(t, input.options, input.srcCode, input.hdrCode, func(idx clang.Index, tu clang.TranslationUnit) {

		if input.tokenFn != nil {
			sourceRange := tu.TranslationUnitCursor().Extent()
			tokens := tu.Tokenize(sourceRange)
			for _, token := range tokens {
				input.tokenFn(tu, token)
			}
		}

		if input.topCursorFn != nil {

			tu.TranslationUnitCursor().Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {

				input.topCursorFn(tu, cursor, parent)

				return clang.ChildVisit_Continue
			})
		}

		if input.fullCursorFn != nil {

			var f func(cursor, parent clang.Cursor) clang.ChildVisitResult
			depth := 0

			f = func(cursor, parent clang.Cursor) clang.ChildVisitResult {

				input.fullCursorFn(tu, cursor, parent, depth, true) // pre is true

				depth++
				cursor.Visit(f)
				depth--

				input.fullCursorFn(tu, cursor, parent, depth, false) // pre is false

				return clang.ChildVisit_Continue
			}

			tu.TranslationUnitCursor().Visit(f)
		}
		if input.tuParseFn != nil {

			input.tuParseFn(&tu)
		}
	})
}

//--
// Define four interfaces, matching the four forms of translation unit visiting we want.
type tokenVisiter interface {
	tokenVisit(tu clang.TranslationUnit, token clang.Token)
}
type topCursorVisiter interface {
	topCursorVisit(tu clang.TranslationUnit, cursor, parent clang.Cursor)
}

type fullCursorVisiter interface {
	fullCursorVisit(tu clang.TranslationUnit, cursor, parent clang.Cursor, depth int, pre bool)
}

type TUParser interface {
	TUParse(tu *clang.TranslationUnit)
}

//-- 1.
// tokenStrings implements tokenVisiter and collects token spelling and token kind strings.
type tokenStrings struct {
	list []string
}

// tokenVisit implements tokenVisiter, collecting results of calls to tokenDescription.
func (x *tokenStrings) tokenVisit(tu clang.TranslationUnit, token clang.Token) {
	x.list = append(x.list, tokenDescription(tu, token))
}
func init() {
	// assert it implements the desired interface.
	var a interface{} = &tokenStrings{}
	_ = a.(tokenVisiter)
}

// tokenDescription returns a string describing the token.
func tokenDescription(tu clang.TranslationUnit, token clang.Token) string {
	return fmt.Sprintf("%s : %s", tu.TokenSpelling(token), token.Kind().String())
}

func tokenDescriptions(tu clang.TranslationUnit, tokens []clang.Token) []string {
	var r []string
	for _, token := range tokens {
		r = append(r, fmt.Sprintf("%s : %s", tu.TokenSpelling(token), token.Kind().String()))
	}
	return r
}

//-- 2.
// topCursorStrings implements topCursorVisiter and collects cursorString results.
type topCursorStrings struct {
	topLevelNamesToSkip map[string]bool
	hdrs                []string
	list                []string
}

// tokenVisit implements topCursorVisiter, collecting results of calls to cursorString.
func (x *topCursorStrings) topCursorVisit(tu clang.TranslationUnit, cursor, parent clang.Cursor) {
	if len(x.list) == 0 {
		x.hdrs = append(x.hdrs, cursorStringHeaders()...)
	}

	// Don't process cursors with names found in the skip map.
	if x.topLevelNamesToSkip != nil {
		name := cursor.Spelling()
		if x.topLevelNamesToSkip[name] {
			return
		}
	}
	x.list = append(x.list, cursorString(cursor, parent))
}
func init() {
	// assert it implements the desired interface.
	var a interface{} = &topCursorStrings{}
	_ = a.(topCursorVisiter)
}

//-- 3.
// fullCursorStrings implements fullCursorVisiter and collects cursorString results.
type fullCursorStrings struct {
	topLevelNamesToSkip map[string]bool
	pad                 string
	list                []string
}

// tokenVisit implements fullCursorVisiter, collecting results of calls to cursorString.
func (x *fullCursorStrings) fullCursorVisit(tu clang.TranslationUnit, cursor, parent clang.Cursor, depth int, pre bool) {
	if depth == 0 && x.topLevelNamesToSkip != nil {
		name := cursor.Spelling()
		if x.topLevelNamesToSkip[name] {
			return
		}
	}

	if x.pad == "" {
		x.pad = ". . . . . . . . . ."
	}
	if pre {
		for depth*2 > len(x.pad) {
			x.pad += x.pad
		}

		s := fmt.Sprintf("%s", x.pad[:depth*2])

		kind := cursor.Kind()

		switch {
		case kind.IsUnexposed():
			s += fmt.Sprintf("IsUnexposed(%s)", kind)
		case kind.IsLiteral(),
			kind.IsExpression():
			is := "?is?"
			switch {
			case kind.IsLiteral():
				is = "IsLiteral"
			case kind.IsExpression():
				is = "IsExpression"
			}
			sourceRange := cursor.Extent()
			tokens := tu.Tokenize(sourceRange)
			if len(tokens) == 1 && tokens[0].Kind() == clang.Token_Literal {
				// Hopefully the normal case where a Literal kind of cursor repreents
				// in a single Literal token. Then use its spelling as the cursor name.
				name := tu.TokenSpelling(tokens[0])
				s += fmt.Sprintf("%s/%s:%s", name, cursor.Kind(), is)
			} else {
				// An unexpected kind of token or there was more than one.
				tokdescs := tokenDescriptions(tu, tokens)
				s += fmt.Sprintf("%s:%s:[%v]", kind, is, strings.Join(tokdescs, ", "))
			}
			s += "/" + cursor.DisplayName()
		case kind.IsStatement():
			s += fmt.Sprintf("%s:IsStatement", cursor.Kind())
		default:
			name := cursor.Spelling()
			s += name
			s += fmt.Sprintf("/%s", cursor.Kind())
		}

		x.list = append(x.list, s)
	} else {
		// There is no post child-visit in this case.
	}
}
func init() {
	// assert it implements the desired interface.
	var a interface{} = &fullCursorStrings{}
	_ = a.(fullCursorVisiter)
}

//-- 4.
// tuParser implements TUParser and collects cursorString results.
type tuParser struct {
	TranslationUnit // The non clang version of that gets populated from the clang version.
}

// tokenVisit implements TUParser, collecting results of calls to cursorString.
func (x *tuParser) TUParse(tu *clang.TranslationUnit) {
	ctu := ClangTranslationUnit{}

	ctu.Populate(tu)

	x.TranslationUnit = ctu.GoTu
}

func init() {
	// assert it implements the desired interface.
	var a interface{} = &tuParser{}
	_ = a.(TUParser)
}

// alignOn returns string created by aligning strings by their common substr.
func alignOn(list []string, substr string) string {
	max := 0
	for _, s := range list {
		i := strings.Index(s, substr)
		if i > max {
			max = i
		}
	}

	b := new(strings.Builder)
	for _, s := range list {
		i := strings.Index(s, substr)
		fmt.Fprintf(b, "%*s\n", max+len(s)-i, s)
	}

	return b.String()
}

var (
	front = regexp.MustCompile("^[ \t]*\n+")
	mid   = regexp.MustCompile("[ \t]*\n[ \t]*")
	end   = regexp.MustCompile("\n[ \t]*$")
)

func trimends(s string) string {
	s = front.ReplaceAllString(s, "")
	s = end.ReplaceAllString(s, "")
	return s
}

func trimmiddle(s string) string {
	s = mid.ReplaceAllString(s, "\n")
	return s
}

func trim(s string) string {
	return trimmiddle(trimends(s))
}

func leadingSpaces(s string) int {
	for i, c := range s {
		if c != ' ' {
			return i
		}
	}
	return len(s)
}

// leftJustify returns a string with left spaces and tabs shaved to left justify
// just enough so that one line is either empty or starts with non whitespace.
func leftJustify(s string) string {
	a := strings.Split(s, "\n")
	if len(a) == 0 {
		return s
	}

	// Loop as long as another space or tab can be shaved from front.
	for {
		if len(a[0]) == 0 {
			goto done
		}
		c := a[0][0]
		if c != ' ' && c != '\t' {
			goto done
		}
		// Do all start with same character?
		for i := range a[1:] {
			if len(a[i]) == 0 || a[i][0] != c {
				goto done
			}
		}
		// All started with same space or tab so shave another character.
		for i := range a {
			a[i] = a[i][1:]
		}
	}

done:
	return strings.Join(a, "\n")
}

type predCall struct {
	name string
	pred func(clang.CursorKind) bool
}

var predCalls = []predCall{
	{
		name: "IsDeclaration",
		pred: clang.CursorKind.IsDeclaration,
	},
	{
		name: "IsReference",
		pred: clang.CursorKind.IsReference,
	},
	{
		name: "IsExpression",
		pred: clang.CursorKind.IsExpression,
	},
	{
		name: "IsStatement",
		pred: clang.CursorKind.IsStatement,
	},
	{
		name: "IsAttribute",
		pred: clang.CursorKind.IsAttribute,
	},
	{
		name: "IsInvalid",
		pred: clang.CursorKind.IsInvalid,
	},
	{
		name: "IsTranslationUnit",
		pred: clang.CursorKind.IsTranslationUnit,
	},
	{
		name: "IsPreprocessing",
		pred: clang.CursorKind.IsPreprocessing,
	},
	{
		name: "IsUnexposed",
		pred: clang.CursorKind.IsUnexposed,
	},
}

func cursorStringHeaders() []string {
	var r []string
	// TBD StorageClass for a clang.Cursor_MacroDefinition doesn't seem appropriate.
	r = append(r, fmt.Sprintf("%s\t%s\t%s\t%s\t%s", "cursor.Spelling():", "cursor.Kind():", "Predicates:", "cursor.StorageClass():", "cursor.Linkage():"))
	r = append(r, fmt.Sprintf("%s\t%s\t%s\t%s\t%s", "                 :", "    .String():", "          :", "            .String():", "                :"))
	return r
}

func cursorString(cursor, parent clang.Cursor) string {
	var predstr []string
	cursorKind := cursor.Kind()
	for _, pred := range predCalls {
		if pred.pred(cursorKind) {
			predstr = append(predstr, pred.name)
		}
	}
	// Cursor_FunctionDecl and Cursor_VarDecl:
	// cursor.StorageClass() is valid for function and variable declarations.
	storageClass := ""
	switch cursor.Kind() {
	case clang.Cursor_FunctionDecl, clang.Cursor_VarDecl:
		storageClass = cursor.StorageClass().String()
	}
	//hasAttrs := cursor.HasAttrs()
	linkageKind := cursor.Linkage()
	// For cursor.Kind() Cursor_MacroDefinition, isPreprocessing().
	s := fmt.Sprintf("%s\t%s\t%s\t%s\t%s", cursor.Spelling(), cursor.Kind().String(), strings.Join(predstr, ","), storageClass, linkageKind)
	return s
}

// alignOnTabs splits input into lines, and then each line by tabs,
// and computes maximum width of each field and returns one string
// with each field aligned and left justified.
func alignOnTabs(input string) string {
	lines := strings.Split(input, "\n")
	var widths []int

	var fieldsList [][]string
	// Compute max width of each field
	for _, line := range lines {
		fields := strings.Split(line, "\t")
		for len(widths) < len(fields) {
			widths = append(widths, 0)
		}

		for i, field := range fields {
			if len(field) > widths[i] {
				widths[i] = len(field)
			}
		}
		fieldsList = append(fieldsList, fields)
	}

	var newlines []string
	for _, fields := range fieldsList {
		var newfields []string
		for i, field := range fields {
			newfields = append(newfields, fmt.Sprintf("%-*s", widths[i], field))
		}
		newlines = append(newlines, strings.Join(newfields, " "))
	}
	return strings.Join(newlines, "\n")
}

type testTokenTuple struct {
	name           string
	options        clang.TranslationUnit_Flags
	hdrCode        string
	srcCode        string
	expectedTokens string
}

func TestTokens(t *testing.T) {
	// Shrink all white space down to a single space
	// when comparing got to expected.

	for _, test := range []*testTokenTuple{
		&test1,
		&test2,
		&test3,
		&test4,
	} {
		t.Run(test.name, func(t *testing.T) {
			// TBD left justify expected string if it is printed
			input := run3Stages{
				options: test.options,
				hdrCode: "",
				srcCode: test.srcCode,
			}

			var ts tokenStrings
			input.run2(t, &ts)

			got := alignOn(ts.list, ":")
			expected := test.expectedTokens

			got = leftJustify(trimends(got))
			expected = leftJustify(trimends(expected))

			if got != expected {
				t.Errorf("%s:\n=== got tokens ===\n%s\n=== expected tokens ===\n%s\n",
					test.name, got, expected)
			}
		})
	}
}

var test1 = testTokenTuple{
	name:    "test1",
	srcCode: "int world();",
	expectedTokens: `
    		  int : Token_Keyword
    		world : Token_Identifier
    		    ( : Token_Punctuation
    		    ) : Token_Punctuation
    		    ; : Token_Punctuation
`,
}

var test2 = testTokenTuple{
	name:    "test2",
	srcCode: "#define A(a) (a + 1)",
	expectedTokens: `
    		     # : Token_Punctuation
    		define : Token_Identifier
    		     A : Token_Identifier
    		     ( : Token_Punctuation
    		     a : Token_Identifier
    		     ) : Token_Punctuation
    		     ( : Token_Punctuation
    		     a : Token_Identifier
    		     + : Token_Punctuation
    		     1 : Token_Literal
    		     ) : Token_Punctuation
`,
}

var test3 = testTokenTuple{
	name: "test3",
	srcCode: `
	#define Add(a) (a + 1)
	int foo(int b) {
		if (b & 0x1) {
			return Add(b);
		}
	}
	`,
	expectedTokens: `
    		     # : Token_Punctuation
    		define : Token_Identifier
    		   Add : Token_Identifier
    		     ( : Token_Punctuation
    		     a : Token_Identifier
    		     ) : Token_Punctuation
    		     ( : Token_Punctuation
    		     a : Token_Identifier
    		     + : Token_Punctuation
    		     1 : Token_Literal
    		     ) : Token_Punctuation
    		   int : Token_Keyword
    		   foo : Token_Identifier
    		     ( : Token_Punctuation
    		   int : Token_Keyword
    		     b : Token_Identifier
    		     ) : Token_Punctuation
    		     { : Token_Punctuation
    		    if : Token_Keyword
    		     ( : Token_Punctuation
    		     b : Token_Identifier
    		     & : Token_Punctuation
    		   0x1 : Token_Literal
    		     ) : Token_Punctuation
    		     { : Token_Punctuation
    		return : Token_Keyword
    		   Add : Token_Identifier
    		     ( : Token_Punctuation
    		     b : Token_Identifier
    		     ) : Token_Punctuation
    		     ; : Token_Punctuation
    		     } : Token_Punctuation
    		     } : Token_Punctuation
`,
}

var test4 = testTokenTuple{
	name: "test4",
	srcCode: `
	static int sa = 7;
	static int getsaI() {
		return sa
	}
	int getsa() {
		return getsaI()
	}
	`,
	expectedTokens: `
    		static : Token_Keyword
    		   int : Token_Keyword
    		    sa : Token_Identifier
    		     = : Token_Punctuation
    		     7 : Token_Literal
    		     ; : Token_Punctuation
    		static : Token_Keyword
    		   int : Token_Keyword
    		getsaI : Token_Identifier
    		     ( : Token_Punctuation
    		     ) : Token_Punctuation
    		     { : Token_Punctuation
    		return : Token_Keyword
    		    sa : Token_Identifier
    		     } : Token_Punctuation
    		   int : Token_Keyword
    		 getsa : Token_Identifier
    		     ( : Token_Punctuation
    		     ) : Token_Punctuation
    		     { : Token_Punctuation
    		return : Token_Keyword
    		getsaI : Token_Identifier
    		     ( : Token_Punctuation
    		     ) : Token_Punctuation
    		     } : Token_Punctuation
`,
}

// compilerTopLevelNames returns a map of all top level names that
// the compiler gives us by default. This is later used to ignore
// those same names when having the compiler work on our test code.
func compilerTopLevelNames(t *testing.T) map[string]bool {
	r := make(map[string]bool)

	// Create an essentially blank source code buffer to compile
	// and record the cursor names that are encountered by the
	// top cursor visit routine.
	input := run3Stages{
		options: clang.TranslationUnit_DetailedPreprocessingRecord,
		hdrCode: " ", // Non empty so the cursor for "hdr.h" is also included
		srcCode: "",

		// Only one callback function needed for this collection.
		topCursorFn: func(tu clang.TranslationUnit, cursor, parent clang.Cursor) {
			r[cursor.Spelling()] = true
		},
	}

	input.run(t) // collect the data
	return r
}

type testTuple struct {
	name                string
	options             clang.TranslationUnit_Flags // May be overridden in test code anyway.
	hdrCode             string
	srcCode             string
	expectedTokens      string
	expectedTopCursors  string
	expectedFullCursors string
	expectedTUPopulate  string
}

func TestCursors(t *testing.T) {

	options := clang.TranslationUnit_DetailedPreprocessingRecord
	topLevelNames := compilerTopLevelNames(t)

	for _, test := range []testTuple{
		ctest3,
		{
			name: "empty-a1",
		},
		{
			name:                "empty-a2",
			expectedTokens:      ` `,
			expectedTopCursors:  ` `,
			expectedFullCursors: ` `,
			expectedTUPopulate: `
    		Tokens:
    		
    		TokenMap:
    		
    		TokenKindMap:
    		
    		TokenNameMap:
    		
    		Cursors:
    		0:{0 1 0 -1 {0 0} {0 0}}
    		CursorKindMap:
    		0:Cursor_TranslationUnit 1:Cursor_MacroDefinition
    		CursorNameMap:
    		0: 1:sample.c 2:__llvm__
    		
			`,
		},
		{
			name:    "global_var_int_a_is_1",
			srcCode: `int a = 1;`,
			expectedTokens: `
    		int : Token_Keyword
    		  a : Token_Identifier
    		  = : Token_Punctuation
    		  1 : Token_Literal
    		  ; : Token_Punctuation
			`,
			expectedTopCursors: `
    		a Cursor_VarDecl IsDeclaration SC_None Linkage_External
			`,
			expectedFullCursors: ``, // TBD
			expectedTUPopulate: `
    		Tokens:
    		0:0 1:1 2:2 3:3 4:4
    		TokenMap:
    		0:{0 0} 1:{1 1} 2:{2 2} 3:{3 3} 4:{2 4}
    		TokenKindMap:
    		0:Token_Keyword 1:Token_Identifier 2:Token_Punctuation 3:Token_Literal
    		TokenNameMap:
    		0:int 1:a 2:= 3:1 4:;
    		Cursors:
    		0:{0 1 0 -1 {0 0} {0 5}} 1:{2 3 0 0 {0 0} {0 4}} 2:{3 0 0 2 {0 0} {3 1}}
    		CursorKindMap:
    		0:Cursor_TranslationUnit 1:Cursor_MacroDefinition 2:Cursor_VarDecl 3:Cursor_IntegerLiteral
    		CursorNameMap:
    		0: 1:sample.c 2:__llvm__ 3:a
    		
			`,
		},
		{
			name:    "global_var_unsigned_int_a_is_0x1",
			srcCode: `unsigned int a = 0x1;`,
			expectedTokens: `
    		unsigned : Token_Keyword
    		     int : Token_Keyword
    		       a : Token_Identifier
    		       = : Token_Punctuation
    		     0x1 : Token_Literal
    		       ; : Token_Punctuation
			`,
			expectedTopCursors: `
    		a Cursor_VarDecl IsDeclaration SC_None Linkage_External
			`,
			expectedFullCursors: ``, // TBD
		},
		{
			name:    "global_var_a_is_0_001",
			srcCode: `double a = 0.001;`,
			expectedTokens: `
    		double : Token_Keyword
    		     a : Token_Identifier
    		     = : Token_Punctuation
    		 0.001 : Token_Literal
    		     ; : Token_Punctuation
			`,
			expectedTopCursors: `
    		a Cursor_VarDecl IsDeclaration SC_None Linkage_External
			`,
			expectedFullCursors: ``, // TBD
		},
		{
			name: "function_storage_comparison",
			srcCode: `int A() { return 1; }
			  extern int B() {return 2; }
			  static int C() {return 3; }
			  inline int D() {return 4; }
			  `,
			expectedTokens: ``,
			expectedTopCursors: `
    		A Cursor_FunctionDecl IsDeclaration SC_None   Linkage_External
    		B Cursor_FunctionDecl IsDeclaration SC_Extern Linkage_External
    		C Cursor_FunctionDecl IsDeclaration SC_Static Linkage_Internal
    		D Cursor_FunctionDecl IsDeclaration SC_None   Linkage_External
	`,
			expectedFullCursors: ``,
			expectedTUPopulate: `
    		Tokens:
    		0:0 1:1 2:2 3:3 4:4 5:5 6:6 7:7 8:8 9:9 10:0 11:10 12:2 13:3 14:4 15:5 16:11 17:7 18:8 19:12 20:0 21:13 22:2 23:3 24:4
    		25:5 26:14 27:7 28:8 29:15 30:0 31:16 32:2 33:3 34:4 35:5 36:17 37:7 38:8
    		TokenMap:
    		0:{0 0} 1:{1 1} 2:{2 2} 3:{2 3} 4:{2 4} 5:{0 5} 6:{3 6} 7:{2 7} 8:{2 8} 9:{0 9} 10:{1 10} 11:{3 11} 12:{0 12} 13:{1 13}
    		14:{3 14} 15:{0 15} 16:{1 16} 17:{3 17}
    		TokenKindMap:
    		0:Token_Keyword 1:Token_Identifier 2:Token_Punctuation 3:Token_Literal
    		TokenNameMap:
    		0:int 1:A 2:( 3:) 4:{ 5:return 6:1 7:; 8:} 9:extern 10:B 11:2 12:static 13:C 14:3 15:inline 16:D 17:4
    		Cursors:
    		0:{0 1 0 -1 {0 0} {0 39}} 1:{2 3 0 0 {0 0} {0 9}} 2:{2 4 0 0 {0 0} {9 10}} 3:{2 5 0 0 {0 0} {19 10}}
    		4:{2 6 0 0 {0 0} {29 10}} 5:{3 0 0 2 {0 0} {4 5}} 6:{3 0 0 3 {0 0} {14 5}} 7:{3 0 0 4 {0 0} {24 5}}
    		8:{3 0 0 5 {0 0} {34 5}} 9:{4 0 0 6 {0 0} {5 2}} 10:{4 0 0 7 {0 0} {15 2}} 11:{4 0 0 8 {0 0} {25 2}}
    		12:{4 0 0 9 {0 0} {35 2}} 13:{5 0 0 10 {0 0} {6 1}} 14:{5 0 0 11 {0 0} {16 1}} 15:{5 0 0 12 {0 0} {26 1}}
    		16:{5 0 0 13 {0 0} {36 1}}
    		CursorKindMap:
    		0:Cursor_TranslationUnit 1:Cursor_MacroDefinition 2:Cursor_FunctionDecl 3:Cursor_CompoundStmt 4:Cursor_ReturnStmt
    		5:Cursor_IntegerLiteral
    		CursorNameMap:
    		0: 1:sample.c 2:__llvm__ 3:A 4:B 5:C 6:D
    		
			`,
		},
		{
			name:    "void_function_and_return",
			srcCode: `void A() { return; }`,
			expectedTokens: `
    		  void : Token_Keyword
    		     A : Token_Identifier
    		     ( : Token_Punctuation
    		     ) : Token_Punctuation
    		     { : Token_Punctuation
    		return : Token_Keyword
    		     ; : Token_Punctuation
    		     } : Token_Punctuation
			`,
			expectedTopCursors: `
    		A Cursor_FunctionDecl IsDeclaration SC_None Linkage_External
			`,
			expectedFullCursors: ``, // TBD
		},
		{
			name:                "int_function_and_return_1",
			srcCode:             `int A() { return 1; }`,
			expectedTokens:      ``,
			expectedTopCursors:  ``,
			expectedFullCursors: ``, // TBD
		},
		{
			name: "parenthesis_comparison",
			srcCode: `int A() { return 1; }
					  int B() { return (2); }
					  int C() { return ((2)); }
					  `,
			expectedTokens:      ``,
			expectedTopCursors:  ``,
			expectedFullCursors: ``, // TBD
		},
		{
			name: "add_sub_mul_div",
			srcCode: `int A() { return 1 + 2 - 3 * 4 / 5; }
					  `,
			expectedTokens: `
    		   int : Token_Keyword
    		     A : Token_Identifier
    		     ( : Token_Punctuation
    		     ) : Token_Punctuation
    		     { : Token_Punctuation
    		return : Token_Keyword
    		     1 : Token_Literal
    		     + : Token_Punctuation
    		     2 : Token_Literal
    		     - : Token_Punctuation
    		     3 : Token_Literal
    		     * : Token_Punctuation
    		     4 : Token_Literal
    		     / : Token_Punctuation
    		     5 : Token_Literal
    		     ; : Token_Punctuation
    		     } : Token_Punctuation
			`,
			expectedTopCursors:  ``,
			expectedFullCursors: ``, // TBD
		},
		{
			name: "add_and_double_add",
			srcCode: `int A() { return 1 + 2; }
					  int B() { return 1 + 2 + 3; }
					  `,
			expectedTokens:      ``,
			expectedTopCursors:  ``,
			expectedFullCursors: ``, // TBD
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			input := run3Stages{
				options: options,
				hdrCode: test.hdrCode,
				srcCode: test.srcCode,
			}
			// Create an object that conforms to all three stages of data collection.
			tr := struct {
				tokenStrings
				topCursorStrings
				fullCursorStrings
				tuParser
			}{}
			tr.topCursorStrings.topLevelNamesToSkip = topLevelNames
			tr.fullCursorStrings.topLevelNamesToSkip = topLevelNames

			input.run2(t, &tr) // collect the data

			if test.expectedTokens != "" {
				got := alignOn(tr.tokenStrings.list, ":")
				got = leftJustify(trimends(got))

				expected := test.expectedTokens
				expected = leftJustify(trimends(expected))

				if got != expected {
					t.Errorf("%s:\n=== got tokens ===\n%s\n=== expected tokens ===\n%s\n",
						test.name, got, expected)
				}
			}

			if test.expectedTopCursors != "" {

				got := strings.Join(tr.topCursorStrings.list, "\n")
				got = leftJustify(alignOnTabs(trimends(got)))

				expected := test.expectedTopCursors
				//expected = trimends(expectedTopCursorsHeader) + "\n" + trimends(expected)
				expected = leftJustify(trimends(expected))

				if got != expected {
					t.Errorf("%s Top cursors:\n=== got ===\n%s\n=== expected ===\n%s\n",
						test.name, got, expected)
				}
			}

			if test.expectedFullCursors != "" {

				got := strings.Join(tr.fullCursorStrings.list, "\n")
				got = leftJustify(trimends(got))

				expected := test.expectedFullCursors
				expected = leftJustify(trimends(expected))

				if got != expected {
					t.Errorf("%s full cursors:\n=== got ===\n%s\n=== expected ===\n%s\n",
						test.name, got, expected)
				}
			}

			if test.expectedTUPopulate != "" {

				got := fmt.Sprintf("%#v", tr.tuParser.TranslationUnit)
				//got = leftJustify(trimends(got))

				expected := test.expectedTUPopulate
				expected = leftJustify(trimends(expected))

				if got != expected {
					t.Errorf("%s TranslationUnit:\n=== got ===\n%s\n=== expected ===\n%s\n",
						test.name, got, expected)
				}
			}
		})
	}
}

const expectedTopCursorsHeader = `
    		cursor.Spelling():	cursor.Kind():	Predicates:	cursor.StorageClass():	cursor.Linkage():
    		                 :	    .String():	          :	            .String():	                :
	`

var ctest1 = testTuple{
	name:           "ctest1",
	options:        clang.TranslationUnit_DetailedPreprocessingRecord,
	hdrCode:        ` `,
	srcCode:        "",
	expectedTokens: ``,
	expectedTopCursors: `
	`,
	expectedFullCursors: `
	`,
}

var ctest2 = testTuple{
	name:    "ctest2",
	options: clang.TranslationUnit_DetailedPreprocessingRecord,
	hdrCode: `
			  extern int hdr_ext_v;
			  static int hdr_sta_v;
			  static int hdr_sta_v = 2;
			  int hdr_glo_v = 3;

			  extern int hdr_ext_f();
			  inline int hdr_inl_f()       { return 11; }
			  static int hdr_sta_f()       { return 12; }
			  int hdr_glo_f()	           { return 13; }
			  static inline int hdr_si_f() { return 14; }
			  `,
	srcCode:        "",
	expectedTokens: ``,
	expectedTopCursors: `
	`,
}

var ctest3 = testTuple{
	name:    "ctest3",
	options: clang.TranslationUnit_DetailedPreprocessingRecord,
	hdrCode: `#define MSG "Hello World!"
			  extern int printf(const char* fmt);
			  `,
	srcCode: `int A() { return 1; }
			  extern int B() {return 2; }
			  static int C() {return 3; }
			  inline int D() {return 4; }
			  `,
	expectedTokens: ``,
	expectedTopCursors: `
    		A Cursor_FunctionDecl IsDeclaration SC_None   Linkage_External
    		B Cursor_FunctionDecl IsDeclaration SC_Extern Linkage_External
    		C Cursor_FunctionDecl IsDeclaration SC_Static Linkage_Internal
    		D Cursor_FunctionDecl IsDeclaration SC_None   Linkage_External
			`,
	expectedFullCursors: ``,
}

// --- 1803 clang ast
// Idea is to create one list of all tokens in tu
// and one list of all cursors in tu. Actually two lists, that mirror each other,
// one of actual clang.cursors and the other of the other of main package Cursors
// that will contain all the computed stuff that can be stored and reloaded without cgo.

// StringId creates indexes for strings. Returning a new index when a new string is
// encountered, and returning the existing index when a string has already been seen.
// The primary API consists of
//
//		Id(string) int
//		ToString(int) string
//
// but helper methods exist for iterating or using with Marshal and Unmarshal
//
//		Len() int
//		Init([]string)
// And publicly accessible
//		Strings []string
//
type StringId struct {
	m       map[string]int
	Strings []string
}

func (s StringId) String() string {
	return fmt.Sprintf("%v", s.Strings)
}
func (s StringId) GoString() string {
	return fmt.Sprintf("{Strings: %#v}", s.Strings)
}

// Len returns number of strings mapped. Valid indexes will be 0..Len-1.
func (s *StringId) Len() int {
	return len(s.Strings)
}

func (s *StringId) Init(strings []string) {
	s.m = make(map[string]int)
	s.Strings = make([]string, len(strings))
	for i, str := range strings {
		s.m[str] = i
		s.Strings[i] = str
	}
}

// Id returns the id for this string.
// First id will be zero.
func (s *StringId) Id(str string) int {
	if s.m == nil {
		s.m = make(map[string]int)
	}
	id, ok := s.m[str]
	if ok {
		return id
	}
	id = len(s.Strings)
	s.m[str] = id
	s.Strings = append(s.Strings, str)
	return id
}

// String returns the string for the given id.
// Panic if id is out of range.
func (s *StringId) ToString(id int) string {
	return s.Strings[id]
}

// TokenMap, same structure as StringId
type TokenMap struct {
	m      map[Token]TokenId
	Tokens []Token
}

func (tm TokenMap) String() string {
	return fmt.Sprintf("%v", tm.Tokens)
}
func (tm TokenMap) GoString() string {
	return fmt.Sprintf("{Tokens: %#v}", tm.Tokens)
}

// Len returns number of Tokens mapped. Valid indexes will be 0..Len-1.
func (tm *TokenMap) Len() int {
	return len(tm.Tokens)
}

func (tm *TokenMap) Init(Tokens []Token) {
	tm.m = make(map[Token]TokenId)
	tm.Tokens = make([]Token, len(Tokens))
	for i, t := range Tokens {
		tm.m[t] = TokenId(i) // cast index to TokenId
		tm.Tokens[i] = t
	}
}

// Id returns the id for this Token.
// First id will be zero.
func (tm *TokenMap) Id(t Token) TokenId {
	if tm.m == nil {
		tm.m = make(map[Token]TokenId)
	}
	id, ok := tm.m[t]
	if ok {
		return id
	}
	id = TokenId(len(tm.Tokens)) // cast len() int to next TokenId
	tm.m[t] = id
	tm.Tokens = append(tm.Tokens, t)
	return id
}

// String returns the string for the given id.
// Panic if id is out of range.
func (tm *TokenMap) ToToken(id TokenId) Token {
	return tm.Tokens[id]
}

type IndexPair struct {
	Head int
	Len  int
}

// Cursor is the pure Go version of the clang Cursor.
// It exists as part of the TranslationUniit Cursors list.
type Cursor struct {
	CursorKindId int // Id into CursorKindMap
	CursorNameId int // Id into CursorNameMap
	Index        int // Own index into the TranslationUnit Cursors list.
	ParentIndex  int //- 1 if Cursor is root.
	Children     IndexPair
	Tokens       IndexPair

	// Index and Children don't have to be initialized within Visit and they don't
	// have to be serialized. They can be recomputed from the Cursor position in
	// the overall list, and the ParentIndex.
}

// Token is the pure Go version of the clang Token.
type Token struct {
	TokenKindId int // Id into TokenKindMap
	TokenNameId int // Id into TokenNameMap
}

type TokenId int // TokenMap allows Token <-> TokenId

/*
type Token [2]int

func (t Token) KindId() int {
	return t[0]
}
func (t Token) NameId() int {
	return t[1]
}
func newToken(kindId, nameId int) Token {
	return Token{kindId, nameId}
}
*/

/*
func (t Token) String() string {
	return fmt.Sprintf("%d %d", t.TokenKindId, t.TokenNameId)
}

func (t Token) GoString() string {
	return fmt.Sprintf("{TokenKindId:%d, TokenNameId:%d}", t.TokenKindId, t.TokenNameId)
}
*/

// TranslationUnit is the pure Go version of the clang translation unit.
// Not using pointers so it can be serialized.
type TranslationUnit struct {
	Cursors  []Cursor  // [0] is the Root Cursor.
	TokenIds []TokenId // TokenId is mapped by TokenMap to get a Token.

	// String/Id maps that allow the Cursor and Token to use int IDs properties.
	CursorKindMap StringId
	CursorNameMap StringId
	TokenMap      TokenMap // Token:"{TokenKindId:2 TokenNameId:3}" mapped to ID
	TokenKindMap  StringId
	TokenNameMap  StringId
}

// numberStrings return string of lines with each slice element shown with its
// slice position prefixed. Before letting line get past width, insert a newline.
func numberStrings(list []string, width int) string {
	b := new(strings.Builder)
	line := 0
	for i, s := range list {
		str := fmt.Sprintf("%d:%s", i, s)
		if 1+line+len(str) > width {
			fmt.Fprintf(b, "\n")
			line = 0
		}
		if line > 0 {
			fmt.Fprintf(b, " ")
			line += 1
		}
		fmt.Fprintf(b, "%s", str)
		line += len(str)
	}

	return b.String()
}

func (t TranslationUnit) GoString() string {
	// Convert two slices to slices of strings
	tokens := make([]string, len(t.TokenIds))
	for i := range tokens {
		tokens[i] = fmt.Sprintf("%v", t.TokenIds[i])
	}
	tm_tokens := make([]string, len(t.TokenMap.Tokens))
	for i := range tm_tokens {
		tm_tokens[i] = fmt.Sprintf("%v", t.TokenMap.Tokens[i])
	}
	cursors := make([]string, len(t.Cursors))
	for i := range cursors {
		cursors[i] = fmt.Sprintf("%v", t.Cursors[i])
	}
	width := 120
	b := new(strings.Builder)
	fmt.Fprintf(b, "Tokens:\n%v\n", numberStrings(tokens, width))
	fmt.Fprintf(b, "TokenMap:\n%v\n", numberStrings(tm_tokens, width))
	fmt.Fprintf(b, "TokenKindMap:\n%v\n", numberStrings(t.TokenKindMap.Strings, width))
	fmt.Fprintf(b, "TokenNameMap:\n%v\n", numberStrings(t.TokenNameMap.Strings, width))
	fmt.Fprintf(b, "Cursors:\n%v\n", numberStrings(cursors, width))
	fmt.Fprintf(b, "CursorKindMap:\n%v\n", numberStrings(t.CursorKindMap.Strings, width))
	fmt.Fprintf(b, "CursorNameMap:\n%v\n", numberStrings(t.CursorNameMap.Strings, width))
	return b.String()
}

// ClangTranslationUnit references the clang package components.
type ClangTranslationUnit struct {
	GoTu            TranslationUnit
	ClangTu         *clang.TranslationUnit
	ClangTokens     []clang.Token
	ClangRootCursor clang.Cursor
	ClangCursors    []clang.Cursor
}

func (ctu *ClangTranslationUnit) convertClangTokens(clangTokens []clang.Token) []TokenId {
	r := make([]TokenId, len(clangTokens))

	for i := range r {
		tokenKindString := clangTokens[i].Kind().String()
		tokenSpelling := ctu.ClangTu.TokenSpelling(clangTokens[i])
		token := Token{
			TokenKindId: ctu.GoTu.TokenKindMap.Id(tokenKindString),
			TokenNameId: ctu.GoTu.TokenNameMap.Id(tokenSpelling),
		}
		tokenId := ctu.GoTu.TokenMap.Id(token)
		r[i] = tokenId
	}
	return r
}

func mapSourceLocationToIndex(tu *clang.TranslationUnit, clangTokens []clang.Token) map[clang.SourceLocation]int {
	r := make(map[clang.SourceLocation]int)

	for i := range clangTokens {
		r[tu.TokenLocation(clangTokens[i])] = i
	}
	return r
}

func (ctu *ClangTranslationUnit) Populate(tu *clang.TranslationUnit) error {
	if ctu.ClangTu != nil {
		return errors.New("Already populated")
	}

	// For some tidyness, have the "" string map to the 0 id.
	_ = ctu.GoTu.CursorNameMap.Id("")

	ctu.ClangTu = tu
	ctu.ClangRootCursor = tu.TranslationUnitCursor()

	ctu.ClangTokens = tu.Tokenize(ctu.ClangRootCursor.Extent())

	ctu.GoTu.TokenIds = ctu.convertClangTokens(ctu.ClangTokens)

	mapTokenIndex := mapSourceLocationToIndex(tu, ctu.ClangTokens)

	// Layer children to end of list, one set of children at a time.
	// Seed list with the root.

	ctu.ClangCursors = append(ctu.ClangCursors, ctu.ClangRootCursor)
	ctu.GoTu.Cursors = append(ctu.GoTu.Cursors, Cursor{
		CursorKindId: ctu.GoTu.CursorKindMap.Id(ctu.ClangRootCursor.Kind().String()),
		CursorNameId: ctu.GoTu.CursorNameMap.Id(ctu.ClangRootCursor.Spelling()),
		ParentIndex:  -1,
		Tokens: IndexPair{ // By definition, all the clang tokens.
			Head: 0,
			Len:  len(ctu.ClangTokens),
		},

		// Index can be set manually later.
		//Index:        0,
	})

	var outsideErr error

	// Grow the list of clang cursors by visiting the list of clang cursors.
	for parentIndex := 0; parentIndex < len(ctu.ClangCursors); parentIndex++ {

		ctu.ClangCursors[parentIndex].Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {

			ctu.ClangCursors = append(ctu.ClangCursors, cursor)

			newCursor := Cursor{
				CursorKindId: ctu.GoTu.CursorKindMap.Id(cursor.Kind().String()),
				CursorNameId: ctu.GoTu.CursorNameMap.Id(cursor.Spelling()),
				ParentIndex:  parentIndex,

				// Index can be set manually later.
				//Index:        len(ctu.GoTu.Cursors),
			}

			// Rest of this is to set newCursor.Tokens.

			// Get clang tokens for this cursor long enough to find them in
			// the global tu list of tokens. It should be enough to get just
			// the first token from the cursor, but there is no libclang call
			// for that.
			tokens := tu.Tokenize(cursor.Extent())

			if len(tokens) > 0 {
				// TBD work against the parent's list first to reduce the search times.
				index, ok := mapTokenIndex[tu.TokenLocation(tokens[0])]
				if !ok {
					outsideErr = errors.New("bug: token location not found in map")
					return clang.ChildVisit_Break
				}

				newCursor.Tokens.Head = index
				newCursor.Tokens.Len = len(tokens)
			}

			ctu.GoTu.Cursors = append(ctu.GoTu.Cursors, newCursor)

			/*
				 * Determining the children doesn't have to be done here.
				 * It is enough that the ParentIndex was set in this visit.
				// Update parent's notion of children.
				//
				// Effectively append own index to parent's list of children.
				// First one through becomes the head, and each one through
				// increases the length by one.
				c := &ctu.GoTu.Cursors[parentIndex].Children
				if c.Head == 0 {
					// ownIndex will never by zero because list starts off with root in it,
					// so length of that list starts at 1.
					c.Head = ownIndex
				}
				c.Len++
			*/

			return clang.ChildVisit_Continue // Continue to next sibling
		})
	}
	if outsideErr != nil {
		return outsideErr
	}

	return nil
}
