package clang_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/frankreh/go-clang-v5.0/ast"
	"github.com/frankreh/go-clang-v5.0/astbridge"
	"github.com/frankreh/go-clang-v5.0/clang"
	"github.com/frankreh/go-clang-v5.0/clang/cursorkind"
	"github.com/frankreh/go-clang-v5.0/clang/tokenkind"
	run "github.com/frankreh/go-clang-v5.0/clangrun"
)

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

//-- 1.
// tokenStrings implements run.TokenVisiter and collects token spelling and token kind strings.
type tokenStrings struct {
	list []string
}

// TokenVisit implements run.TokenVisiter, collecting results of calls to tokenDescription.
func (x *tokenStrings) TokenVisit(tu clang.TranslationUnit, token clang.Token) {
	x.list = append(x.list, tokenDescription(tu, token))
}
func init() {
	// assert it implements the desired interface.
	var a interface{} = &tokenStrings{}
	_ = a.(run.TokenVisiter)
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
// topCursorStrings implements run.TopCursorVisiter and collects cursorString results.
type topCursorStrings struct {
	topLevelNamesToSkip map[string]bool
	hdrs                []string
	list                []string
}

// TopCursorVisit implements run.TopCursorVisiter, collecting results of calls to cursorString.
func (x *topCursorStrings) TopCursorVisit(tu clang.TranslationUnit, cursor, parent clang.Cursor) {
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
	_ = a.(run.TopCursorVisiter)
}

//-- 3.
// fullCursorStrings implements run.FullCursorVisiter and collects cursorString results.
type fullCursorStrings struct {
	topLevelNamesToSkip map[string]bool
	depthMap            map[clang.Cursor]int // 0 When parent is root, -1 when parent is one we are skipping.
	pad                 string
	list                []string
}

// FullCursorVisit implements run.FullCursorVisiter, collecting results of calls to cursorString.
func (x *fullCursorStrings) FullCursorVisit(tu clang.TranslationUnit, cursor, parent clang.Cursor) {
	// parent hash is not expected to be in the map so fact zero is returned is useful.
	if x.depthMap == nil {
		x.depthMap = make(map[clang.Cursor]int)
	}

	depth := x.depthMap[parent]

	if depth == 0 && x.topLevelNamesToSkip != nil {
		name := cursor.Spelling()
		if x.topLevelNamesToSkip[name] {
			x.depthMap[cursor] = -1
			return
		}
	}
	if depth == -1 {
		x.depthMap[cursor] = -1
		return
	}
	x.depthMap[cursor] = depth + 1

	if x.pad == "" {
		x.pad = ". . . . . . . . . ."
	}
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
		if len(tokens) == 1 && tokens[0].Kind() == tokenkind.Literal {
			// Hopefully the normal case where a Literal kind of cursor represents
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
}
func init() {
	// assert it implements the desired interface.
	var a interface{} = &fullCursorStrings{}
	_ = a.(run.FullCursorVisiter)
}

//-- 4.
// tuParser implements run.TUParser and collects cursorString results.
type tuParser struct {
	ast.TranslationUnit // The non clang version of that gets populated from the clang version.
}

// tokenVisit implements TUParser, collecting results of calls to cursorString.
func (x *tuParser) TUParse(tu *clang.TranslationUnit) {
	ctu := astbridge.ClangTranslationUnit{}

	ctu.Populate(tu)

	x.TranslationUnit = ctu.GoTu
}

func init() {
	// assert it implements the desired interface.
	var a interface{} = &tuParser{}
	_ = a.(run.TUParser)
}

// Utilities

func cursorStringHeaders() []string {
	var r []string
	// TBD StorageClass for a cursorkind.MacroDefinition doesn't seem appropriate.
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
	// FunctionDecl and VarDecl:
	// cursor.StorageClass() is valid for function and variable declarations.
	storageClass := ""
	switch cursor.Kind() {
	case cursorkind.FunctionDecl, cursorkind.VarDecl:
		storageClass = cursor.StorageClass().String()
	}
	//hasAttrs := cursor.HasAttrs()
	linkageKind := cursor.Linkage()
	// For cursor.Kind() MacroDefinition, isPreprocessing().
	s := fmt.Sprintf("%s\t%s\t%s\t%s\t%s", cursor.Spelling(), cursor.Kind().String(), strings.Join(predstr, ","), storageClass, linkageKind)
	return s
}

type predCall struct {
	name string
	pred func(cursorkind.Kind) bool
}

var predCalls = []predCall{
	{
		name: "IsDeclaration",
		pred: cursorkind.Kind.IsDeclaration,
	},
	{
		name: "IsReference",
		pred: cursorkind.Kind.IsReference,
	},
	{
		name: "IsExpression",
		pred: cursorkind.Kind.IsExpression,
	},
	{
		name: "IsStatement",
		pred: cursorkind.Kind.IsStatement,
	},
	{
		name: "IsAttribute",
		pred: cursorkind.Kind.IsAttribute,
	},
	{
		name: "IsInvalid",
		pred: cursorkind.Kind.IsInvalid,
	},
	{
		name: "IsTranslationUnit",
		pred: cursorkind.Kind.IsTranslationUnit,
	},
	{
		name: "IsPreprocessing",
		pred: cursorkind.Kind.IsPreprocessing,
	},
	{
		name: "IsUnexposed",
		pred: cursorkind.Kind.IsUnexposed,
	},
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
			input := run.Callbacks{
				Options: test.options,
				HdrCode: "",
				SrcCode: test.srcCode,
			}

			var ts tokenStrings
			err := input.LayerAndExecute(&ts)
			if err != nil {
				t.Fatal(err)
			}

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
    		  int : Keyword
    		world : Identifier
    		    ( : Punctuation
    		    ) : Punctuation
    		    ; : Punctuation
`,
}

var test2 = testTokenTuple{
	name:    "test2",
	srcCode: "#define A(a) (a + 1)",
	expectedTokens: `
    		     # : Punctuation
    		define : Identifier
    		     A : Identifier
    		     ( : Punctuation
    		     a : Identifier
    		     ) : Punctuation
    		     ( : Punctuation
    		     a : Identifier
    		     + : Punctuation
    		     1 : Literal
    		     ) : Punctuation
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
    		     # : Punctuation
    		define : Identifier
    		   Add : Identifier
    		     ( : Punctuation
    		     a : Identifier
    		     ) : Punctuation
    		     ( : Punctuation
    		     a : Identifier
    		     + : Punctuation
    		     1 : Literal
    		     ) : Punctuation
    		   int : Keyword
    		   foo : Identifier
    		     ( : Punctuation
    		   int : Keyword
    		     b : Identifier
    		     ) : Punctuation
    		     { : Punctuation
    		    if : Keyword
    		     ( : Punctuation
    		     b : Identifier
    		     & : Punctuation
    		   0x1 : Literal
    		     ) : Punctuation
    		     { : Punctuation
    		return : Keyword
    		   Add : Identifier
    		     ( : Punctuation
    		     b : Identifier
    		     ) : Punctuation
    		     ; : Punctuation
    		     } : Punctuation
    		     } : Punctuation
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
    		static : Keyword
    		   int : Keyword
    		    sa : Identifier
    		     = : Punctuation
    		     7 : Literal
    		     ; : Punctuation
    		static : Keyword
    		   int : Keyword
    		getsaI : Identifier
    		     ( : Punctuation
    		     ) : Punctuation
    		     { : Punctuation
    		return : Keyword
    		    sa : Identifier
    		     } : Punctuation
    		   int : Keyword
    		 getsa : Identifier
    		     ( : Punctuation
    		     ) : Punctuation
    		     { : Punctuation
    		return : Keyword
    		getsaI : Identifier
    		     ( : Punctuation
    		     ) : Punctuation
    		     } : Punctuation
`,
}

// compilerTopLevelNames returns a map of all top level names that
// the compiler gives us by default. This is later used to ignore
// those same names when having the compiler work on our test code.
func compilerTopLevelNames() (map[string]bool, error) {
	r := make(map[string]bool)

	// Create an essentially blank source code buffer to compile
	// and record the cursor names that are encountered by the
	// top cursor visit routine.
	input := run.Callbacks{
		Options: clang.TranslationUnit_DetailedPreprocessingRecord,
		HdrCode: " ", // Non empty so the cursor for "hdr.h" is also included
		SrcCode: "",
	}
	// Only one callback function needed for this collection.
	input.AppendTopCursorFn(
		func(tu clang.TranslationUnit, cursor, parent clang.Cursor) {
			r[cursor.Spelling()] = true
		})

	err := input.Execute() // collect the data
	return r, err
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
	expectedGobSize0    int
	expectedGobSize1    int
}

func TestCursors(t *testing.T) {

	options := clang.TranslationUnit_DetailedPreprocessingRecord
	topLevelNames, err := compilerTopLevelNames()
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []testTuple{
		{
			name:             "empty_a1",
			expectedGobSize0: 513,
			expectedGobSize1: 276,
		},
		{
			name:                "empty_a2",
			expectedTokens:      ` `,
			expectedTopCursors:  ` `,
			expectedFullCursors: ` `,
			expectedTUPopulate: `
    		Tokens:
    		
    		TokenMap:
    		
    		TokenNameMap:
    		
    		Cursors:
    		0:{TranslationUnit 1 -1 {0 0} {0 0}}
    		CursorNameMap:
    		0: 1:sample.c
			`,
			expectedGobSize0: 513,
			expectedGobSize1: 276,
		},
		{
			name:             "gob_empty",
			expectedGobSize0: 513,
			expectedGobSize1: 276,
		},
		{
			name: "gob_foo",
			srcCode: `
			void foo() { }
			`,
			expectedGobSize0: 602,
			expectedGobSize1: 339,
		},
		{
			name: "gob_foo_bar",
			srcCode: `
				void foo() { }
				void bar() { }
				`,
			expectedGobSize0: 654,
			expectedGobSize1: 370,
		},
		{
			name:    "global_var_int_a_is_1",
			srcCode: `int a = 1;`,
			expectedTokens: `
    		int : Keyword
    		  a : Identifier
    		  = : Punctuation
    		  1 : Literal
    		  ; : Punctuation
			`,
			expectedTopCursors: `
    		a VarDecl IsDeclaration SC_None Linkage_External
			`,
			expectedFullCursors: `
    		a/VarDecl
    		. 1/IntegerLiteral:IsLiteral/
			`,
			expectedTUPopulate: `
    		Tokens:
    		0:0 1:1 2:2 3:3 4:4
    		TokenMap:
    		0:{Keyword 0} 1:{Identifier 1} 2:{Punctuation 2} 3:{Literal 3} 4:{Punctuation 4}
    		TokenNameMap:
    		0:int 1:a 2:= 3:1 4:;
    		Cursors:
    		0:{TranslationUnit 1 -1 {1 1} {0 5}} 1:{VarDecl 2 0 {2 1} {0 4}} 2:{IntegerLiteral 0 1 {0 0} {3 1}}
    		CursorNameMap:
    		0: 1:sample.c 2:a
			`,
			expectedGobSize0: 592,
			expectedGobSize1: 329,
		},
		{
			name:    "global_var_unsigned_int_a_is_0x1",
			srcCode: `unsigned int a = 0x1;`,
			expectedTokens: `
    		unsigned : Keyword
    		     int : Keyword
    		       a : Identifier
    		       = : Punctuation
    		     0x1 : Literal
    		       ; : Punctuation
			`,
			expectedTopCursors: `
    		a VarDecl IsDeclaration SC_None Linkage_External
			`,
			expectedFullCursors: `
    		a/VarDecl
    		. IsUnexposed(UnexposedExpr)
    		. . 0x1/IntegerLiteral:IsLiteral/
			`,
			expectedGobSize0: 628,
			expectedGobSize1: 352,
		},
		{
			name:    "global_var_a_is_0_001",
			srcCode: `double a = 0.001;`,
			expectedTokens: `
    		double : Keyword
    		     a : Identifier
    		     = : Punctuation
    		 0.001 : Literal
    		     ; : Punctuation
			`,
			expectedTopCursors: `
    		a VarDecl IsDeclaration SC_None Linkage_External
			`,
			expectedFullCursors: `
    		a/VarDecl
    		. 0.001/FloatingLiteral:IsLiteral/
			`,
			expectedGobSize0: 599,
			expectedGobSize1: 336,
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
    		A FunctionDecl IsDeclaration SC_None   Linkage_External
    		B FunctionDecl IsDeclaration SC_Extern Linkage_External
    		C FunctionDecl IsDeclaration SC_Static Linkage_Internal
    		D FunctionDecl IsDeclaration SC_None   Linkage_External
			`,
			expectedFullCursors: ``,
			expectedTUPopulate: `
    		Tokens:
    		0:0 1:1 2:2 3:3 4:4 5:5 6:6 7:7 8:8 9:9 10:0 11:10 12:2 13:3 14:4 15:5 16:11 17:7 18:8 19:12 20:0 21:13 22:2 23:3 24:4
    		25:5 26:14 27:7 28:8 29:15 30:0 31:16 32:2 33:3 34:4 35:5 36:17 37:7 38:8
    		TokenMap:
    		0:{Keyword 0} 1:{Identifier 1} 2:{Punctuation 2} 3:{Punctuation 3} 4:{Punctuation 4} 5:{Keyword 5} 6:{Literal 6}
    		7:{Punctuation 7} 8:{Punctuation 8} 9:{Keyword 9} 10:{Identifier 10} 11:{Literal 11} 12:{Keyword 12} 13:{Identifier 13}
    		14:{Literal 14} 15:{Keyword 15} 16:{Identifier 16} 17:{Literal 17}
    		TokenNameMap:
    		0:int 1:A 2:( 3:) 4:{ 5:return 6:1 7:; 8:} 9:extern 10:B 11:2 12:static 13:C 14:3 15:inline 16:D 17:4
    		Cursors:
    		0:{TranslationUnit 1 -1 {1 4} {0 39}} 1:{FunctionDecl 2 0 {5 1} {0 9}} 2:{FunctionDecl 3 0 {6 1} {9 10}}
    		3:{FunctionDecl 4 0 {7 1} {19 10}} 4:{FunctionDecl 5 0 {8 1} {29 10}} 5:{CompoundStmt 0 1 {9 1} {4 5}}
    		6:{CompoundStmt 0 2 {10 1} {14 5}} 7:{CompoundStmt 0 3 {11 1} {24 5}} 8:{CompoundStmt 0 4 {12 1} {34 5}}
    		9:{ReturnStmt 0 5 {13 1} {5 2}} 10:{ReturnStmt 0 6 {14 1} {15 2}} 11:{ReturnStmt 0 7 {15 1} {25 2}}
    		12:{ReturnStmt 0 8 {16 1} {35 2}} 13:{IntegerLiteral 0 9 {0 0} {6 1}} 14:{IntegerLiteral 0 10 {0 0} {16 1}}
    		15:{IntegerLiteral 0 11 {0 0} {26 1}} 16:{IntegerLiteral 0 12 {0 0} {36 1}}
    		CursorNameMap:
    		0: 1:sample.c 2:A 3:B 4:C 5:D
			`,
			expectedGobSize0: 984,
			expectedGobSize1: 563,
		},
		{
			name:    "void_function_and_return",
			srcCode: `void A() { return; }`,
			expectedTokens: `
    		  void : Keyword
    		     A : Identifier
    		     ( : Punctuation
    		     ) : Punctuation
    		     { : Punctuation
    		return : Keyword
    		     ; : Punctuation
    		     } : Punctuation
			`,
			expectedTopCursors: `
    		A FunctionDecl IsDeclaration SC_None Linkage_External
			`,
			expectedFullCursors: `
    		A/FunctionDecl
    		. CompoundStmt:IsStatement
    		. . ReturnStmt:IsStatement
			`,
			expectedGobSize0: 637,
			expectedGobSize1: 361,
		},
		{
			name:               "int_function_and_return_1",
			srcCode:            `int A() { return 1; }`,
			expectedTokens:     ``,
			expectedTopCursors: ``,
			expectedFullCursors: `
    		A/FunctionDecl
    		. CompoundStmt:IsStatement
    		. . ReturnStmt:IsStatement
    		. . . 1/IntegerLiteral:IsLiteral/
			`,
			expectedGobSize0: 662,
			expectedGobSize1: 374,
		},
		{
			name: "parenthesis_comparison",
			srcCode: `int A() { return 1; }
					  int B() { return (2); }
					  int C() { return ((2)); }
					  `,
			expectedTokens:     ``,
			expectedTopCursors: ``,
			expectedFullCursors: `
    		A/FunctionDecl
    		. CompoundStmt:IsStatement
    		. . ReturnStmt:IsStatement
    		. . . 1/IntegerLiteral:IsLiteral/
    		B/FunctionDecl
    		. CompoundStmt:IsStatement
    		. . ReturnStmt:IsStatement
    		. . . ParenExpr:IsExpression:[( : Punctuation, 2 : Literal, ) : Punctuation]/
    		. . . . 2/IntegerLiteral:IsLiteral/
    		C/FunctionDecl
    		. CompoundStmt:IsStatement
    		. . ReturnStmt:IsStatement
    		. . . ParenExpr:IsExpression:[( : Punctuation, ( : Punctuation, 2 : Literal, ) : Punctuation, ) : Punctuation]/
    		. . . . ParenExpr:IsExpression:[( : Punctuation, 2 : Literal, ) : Punctuation]/
    		. . . . . 2/IntegerLiteral:IsLiteral/
			`,
			expectedGobSize0: 904,
			expectedGobSize1: 491,
		},
		{
			name: "func_sub_77_and_78",
			srcCode: `int A() { return 77 - 78; }
					  `,
			expectedTokens: `
    		   int : Keyword
    		     A : Identifier
    		     ( : Punctuation
    		     ) : Punctuation
    		     { : Punctuation
    		return : Keyword
    		    77 : Literal
    		     - : Punctuation
    		    78 : Literal
    		     ; : Punctuation
    		     } : Punctuation
			`,
			expectedTopCursors: ``,
			expectedFullCursors: `
    		A/FunctionDecl
    		. CompoundStmt:IsStatement
    		. . ReturnStmt:IsStatement
    		. . . BinaryOperator:IsExpression:[77 : Literal, - : Punctuation, 78 : Literal]/
    		. . . . 77/IntegerLiteral:IsLiteral/
    		. . . . 78/IntegerLiteral:IsLiteral/
			`,
			expectedTUPopulate: `
    		Tokens:
    		0:0 1:1 2:2 3:3 4:4 5:5 6:6 7:7 8:8 9:9 10:10
    		TokenMap:
    		0:{Keyword 0} 1:{Identifier 1} 2:{Punctuation 2} 3:{Punctuation 3} 4:{Punctuation 4} 5:{Keyword 5} 6:{Literal 6}
    		7:{Punctuation 7} 8:{Literal 8} 9:{Punctuation 9} 10:{Punctuation 10}
    		TokenNameMap:
    		0:int 1:A 2:( 3:) 4:{ 5:return 6:77 7:- 8:78 9:; 10:}
    		Cursors:
    		0:{TranslationUnit 1 -1 {1 1} {0 11}} 1:{FunctionDecl 2 0 {2 1} {0 11}} 2:{CompoundStmt 0 1 {3 1} {4 7}}
    		3:{ReturnStmt 0 2 {4 1} {5 4}} 4:{BinaryOperator 3 3 {5 2} {6 3}} 5:{IntegerLiteral 0 4 {0 0} {6 1}}
    		6:{IntegerLiteral 0 4 {0 0} {8 1}}
    		CursorNameMap:
    		0: 1:sample.c 2:A 3:-
			`,
			expectedGobSize0: 714,
			expectedGobSize1: 404,
		},
		{
			name:    "var_sub_77_and_78",
			srcCode: `int a = 77 - 78;`,
			expectedTokens: `
    		int : Keyword
    		  a : Identifier
    		  = : Punctuation
    		 77 : Literal
    		  - : Punctuation
    		 78 : Literal
    		  ; : Punctuation
			`,
			expectedTopCursors: ``,
			expectedFullCursors: `
    		a/VarDecl
    		. BinaryOperator:IsExpression:[77 : Literal, - : Punctuation, 78 : Literal]/
    		. . 77/IntegerLiteral:IsLiteral/
    		. . 78/IntegerLiteral:IsLiteral/
			`,
			expectedTUPopulate: `
    		Tokens:
    		0:0 1:1 2:2 3:3 4:4 5:5 6:6
    		TokenMap:
    		0:{Keyword 0} 1:{Identifier 1} 2:{Punctuation 2} 3:{Literal 3} 4:{Punctuation 4} 5:{Literal 5} 6:{Punctuation 6}
    		TokenNameMap:
    		0:int 1:a 2:= 3:77 4:- 5:78 6:;
    		Cursors:
    		0:{TranslationUnit 1 -1 {1 1} {0 7}} 1:{VarDecl 2 0 {2 1} {0 6}} 2:{BinaryOperator 3 1 {3 2} {3 3}}
    		3:{IntegerLiteral 0 2 {0 0} {3 1}} 4:{IntegerLiteral 0 2 {0 0} {5 1}}
    		CursorNameMap:
    		0: 1:sample.c 2:a 3:-
			`,
			expectedGobSize0: 645,
			expectedGobSize1: 359,
		},
		{
			name:    "var_sub_p77p_and_78",
			srcCode: `int a = (77) - 78;`,
			// BinaryOperator - has 5 tokens and 2 children. First child has 3 tokens. Means the 4th is the binary operator symbol.
			expectedTokens: `
    		int : Keyword
    		  a : Identifier
    		  = : Punctuation
    		  ( : Punctuation
    		 77 : Literal
    		  ) : Punctuation
    		  - : Punctuation
    		 78 : Literal
    		  ; : Punctuation
			`,
			expectedTopCursors: ``,
			expectedFullCursors: `
    		a/VarDecl
    		. BinaryOperator:IsExpression:[( : Punctuation, 77 : Literal, ) : Punctuation, - : Punctuation, 78 : Literal]/
    		. . ParenExpr:IsExpression:[( : Punctuation, 77 : Literal, ) : Punctuation]/
    		. . . 77/IntegerLiteral:IsLiteral/
    		. . 78/IntegerLiteral:IsLiteral/
			`,
			expectedTUPopulate: `
    		Tokens:
    		0:0 1:1 2:2 3:3 4:4 5:5 6:6 7:7 8:8
    		TokenMap:
    		0:{Keyword 0} 1:{Identifier 1} 2:{Punctuation 2} 3:{Punctuation 3} 4:{Literal 4} 5:{Punctuation 5} 6:{Punctuation 6}
    		7:{Literal 7} 8:{Punctuation 8}
    		TokenNameMap:
    		0:int 1:a 2:= 3:( 4:77 5:) 6:- 7:78 8:;
    		Cursors:
    		0:{TranslationUnit 1 -1 {1 1} {0 9}} 1:{VarDecl 2 0 {2 1} {0 8}} 2:{BinaryOperator 3 1 {3 2} {3 5}}
    		3:{ParenExpr 0 2 {5 1} {3 3}} 4:{IntegerLiteral 0 2 {0 0} {7 1}} 5:{IntegerLiteral 0 3 {0 0} {4 1}}
    		CursorNameMap:
    		0: 1:sample.c 2:a 3:-
			`,
			expectedGobSize0: 675,
			expectedGobSize1: 377,
		},
		{
			name: "var_a_is_minus_b",
			srcCode: `int a = 99;
					  int b = -a;`,
			// UnaryOperator - has 2 tokens and 1 children. First child has 1 token. Means the 1st is the unary operator symbol.
			expectedTokens: `
    		int : Keyword
    		  a : Identifier
    		  = : Punctuation
    		 99 : Literal
    		  ; : Punctuation
    		int : Keyword
    		  b : Identifier
    		  = : Punctuation
    		  - : Punctuation
    		  a : Identifier
    		  ; : Punctuation
			`,
			expectedTopCursors: ``,
			expectedFullCursors: `
    		a/VarDecl
    		. 99/IntegerLiteral:IsLiteral/
    		b/VarDecl
    		. UnaryOperator:IsExpression:[- : Punctuation, a : Identifier]/
    		. . IsUnexposed(UnexposedExpr)
    		. . . DeclRefExpr:IsExpression:[a : Identifier]/a
			`,
			expectedTUPopulate: `
    		Tokens:
    		0:0 1:1 2:2 3:3 4:4 5:0 6:5 7:2 8:6 9:1 10:4
    		TokenMap:
    		0:{Keyword 0} 1:{Identifier 1} 2:{Punctuation 2} 3:{Literal 3} 4:{Punctuation 4} 5:{Identifier 5} 6:{Punctuation 6}
    		TokenNameMap:
    		0:int 1:a 2:= 3:99 4:; 5:b 6:-
    		Cursors:
    		0:{TranslationUnit 1 -1 {1 2} {0 11}} 1:{VarDecl 2 0 {3 1} {0 4}} 2:{VarDecl 3 0 {4 1} {5 5}}
    		3:{IntegerLiteral 0 1 {0 0} {3 1}} 4:{UnaryOperator 0 2 {5 1} {8 2}} 5:{UnexposedExpr 2 4 {6 1} {9 1}}
    		6:{DeclRefExpr 2 5 {0 0} {9 1}}
    		CursorNameMap:
    		0: 1:sample.c 2:a 3:b
			`,
			expectedGobSize0: 685,
			expectedGobSize1: 373,
		},
		{
			name: "add_sub_mul_div",
			srcCode: `int A() { return 1 + 2 - 3 * 4 / 5; }
					  `,
			expectedTokens: `
    		   int : Keyword
    		     A : Identifier
    		     ( : Punctuation
    		     ) : Punctuation
    		     { : Punctuation
    		return : Keyword
    		     1 : Literal
    		     + : Punctuation
    		     2 : Literal
    		     - : Punctuation
    		     3 : Literal
    		     * : Punctuation
    		     4 : Literal
    		     / : Punctuation
    		     5 : Literal
    		     ; : Punctuation
    		     } : Punctuation
			`,
			expectedTopCursors: ``,
			expectedFullCursors: `
    		A/FunctionDecl
    		. CompoundStmt:IsStatement
    		. . ReturnStmt:IsStatement
    		. . . BinaryOperator:IsExpression:[1 : Literal, + : Punctuation, 2 : Literal, - : Punctuation, 3 : Literal, * : Punctuation, 4 : Literal, / : Punctuation, 5 : Literal]/
    		. . . . BinaryOperator:IsExpression:[1 : Literal, + : Punctuation, 2 : Literal]/
    		. . . . . 1/IntegerLiteral:IsLiteral/
    		. . . . . 2/IntegerLiteral:IsLiteral/
    		. . . . BinaryOperator:IsExpression:[3 : Literal, * : Punctuation, 4 : Literal, / : Punctuation, 5 : Literal]/
    		. . . . . BinaryOperator:IsExpression:[3 : Literal, * : Punctuation, 4 : Literal]/
    		. . . . . . 3/IntegerLiteral:IsLiteral/
    		. . . . . . 4/IntegerLiteral:IsLiteral/
    		. . . . . 5/IntegerLiteral:IsLiteral/
			`,
			expectedTUPopulate: `
    		Tokens:
    		0:0 1:1 2:2 3:3 4:4 5:5 6:6 7:7 8:8 9:9 10:10 11:11 12:12 13:13 14:14 15:15 16:16
    		TokenMap:
    		0:{Keyword 0} 1:{Identifier 1} 2:{Punctuation 2} 3:{Punctuation 3} 4:{Punctuation 4} 5:{Keyword 5} 6:{Literal 6}
    		7:{Punctuation 7} 8:{Literal 8} 9:{Punctuation 9} 10:{Literal 10} 11:{Punctuation 11} 12:{Literal 12}
    		13:{Punctuation 13} 14:{Literal 14} 15:{Punctuation 15} 16:{Punctuation 16}
    		TokenNameMap:
    		0:int 1:A 2:( 3:) 4:{ 5:return 6:1 7:+ 8:2 9:- 10:3 11:* 12:4 13:/ 14:5 15:; 16:}
    		Cursors:
    		0:{TranslationUnit 1 -1 {1 1} {0 17}} 1:{FunctionDecl 2 0 {2 1} {0 17}} 2:{CompoundStmt 0 1 {3 1} {4 13}}
    		3:{ReturnStmt 0 2 {4 1} {5 10}} 4:{BinaryOperator 3 3 {5 2} {6 9}} 5:{BinaryOperator 4 4 {7 2} {6 3}}
    		6:{BinaryOperator 5 4 {9 2} {10 5}} 7:{IntegerLiteral 0 5 {0 0} {6 1}} 8:{IntegerLiteral 0 5 {0 0} {8 1}}
    		9:{BinaryOperator 6 6 {11 2} {10 3}} 10:{IntegerLiteral 0 6 {0 0} {14 1}} 11:{IntegerLiteral 0 9 {0 0} {10 1}}
    		12:{IntegerLiteral 0 9 {0 0} {12 1}}
    		CursorNameMap:
    		0: 1:sample.c 2:A 3:- 4:+ 5:/ 6:*
			`,
			expectedGobSize0: 863,
			expectedGobSize1: 486,
		},
		{
			name: "add_and_double_add",
			srcCode: `int A() { return 1 + 2; }
					  int B() { return 1 + 2 + 3; }
					  `,
			expectedTokens:     ``,
			expectedTopCursors: ``,
			expectedFullCursors: `
    		A/FunctionDecl
    		. CompoundStmt:IsStatement
    		. . ReturnStmt:IsStatement
    		. . . BinaryOperator:IsExpression:[1 : Literal, + : Punctuation, 2 : Literal]/
    		. . . . 1/IntegerLiteral:IsLiteral/
    		. . . . 2/IntegerLiteral:IsLiteral/
    		B/FunctionDecl
    		. CompoundStmt:IsStatement
    		. . ReturnStmt:IsStatement
    		. . . BinaryOperator:IsExpression:[1 : Literal, + : Punctuation, 2 : Literal, + : Punctuation, 3 : Literal]/
    		. . . . BinaryOperator:IsExpression:[1 : Literal, + : Punctuation, 2 : Literal]/
    		. . . . . 1/IntegerLiteral:IsLiteral/
    		. . . . . 2/IntegerLiteral:IsLiteral/
    		. . . . 3/IntegerLiteral:IsLiteral/
			`,
			expectedGobSize0: 879,
			expectedGobSize1: 480,
		},
		{
			name:    "hdr_processing_only",
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
			srcCode:            "",
			expectedTokens:     ``,
			expectedTopCursors: ` `,
			expectedGobSize0:   576,
			expectedGobSize1:   325,
		},
		{
			name:    "hdr_and_source",
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
    		A FunctionDecl IsDeclaration SC_None   Linkage_External
    		B FunctionDecl IsDeclaration SC_Extern Linkage_External
    		C FunctionDecl IsDeclaration SC_Static Linkage_Internal
    		D FunctionDecl IsDeclaration SC_None   Linkage_External
			`,
			expectedFullCursors: `
    		A/FunctionDecl
    		. CompoundStmt:IsStatement
    		. . ReturnStmt:IsStatement
    		. . . 1/IntegerLiteral:IsLiteral/
    		B/FunctionDecl
    		. CompoundStmt:IsStatement
    		. . ReturnStmt:IsStatement
    		. . . 2/IntegerLiteral:IsLiteral/
    		C/FunctionDecl
    		. CompoundStmt:IsStatement
    		. . ReturnStmt:IsStatement
    		. . . 3/IntegerLiteral:IsLiteral/
    		D/FunctionDecl
    		. CompoundStmt:IsStatement
    		. . ReturnStmt:IsStatement
    		. . . 4/IntegerLiteral:IsLiteral/
			`,
			expectedGobSize0: 1039,
			expectedGobSize1: 610,
		},
		{
			name:    "struct_a_b",
			options: clang.TranslationUnit_DetailedPreprocessingRecord,
			hdrCode: ``,
			srcCode: `
				  struct StructS {
					  int a;
					  char b;
				  } VarS;
				  `,
			expectedTokens: `
    		 struct : Keyword
    		StructS : Identifier
    		      { : Punctuation
    		    int : Keyword
    		      a : Identifier
    		      ; : Punctuation
    		   char : Keyword
    		      b : Identifier
    		      ; : Punctuation
    		      } : Punctuation
    		   VarS : Identifier
    		      ; : Punctuation
				`,
			expectedTopCursors: `
    		StructS StructDecl IsDeclaration         Linkage_External
    		VarS    VarDecl    IsDeclaration SC_None Linkage_External
				`,
			expectedFullCursors: `
    		StructS/StructDecl
    		. a/FieldDecl
    		. b/FieldDecl
    		VarS/VarDecl
    		. StructS/StructDecl
    		. . a/FieldDecl
    		. . b/FieldDecl
				`,
			expectedTUPopulate: `
    		Tokens:
    		0:0 1:1 2:2 3:3 4:4 5:5 6:6 7:7 8:5 9:8 10:9 11:5
    		TokenMap:
    		0:{Keyword 0} 1:{Identifier 1} 2:{Punctuation 2} 3:{Keyword 3} 4:{Identifier 4} 5:{Punctuation 5} 6:{Keyword 6}
    		7:{Identifier 7} 8:{Punctuation 8} 9:{Identifier 9}
    		TokenNameMap:
    		0:struct 1:StructS 2:{ 3:int 4:a 5:; 6:char 7:b 8:} 9:VarS
    		Cursors:
    		0:{TranslationUnit 1 -1 {1 2} {0 12}} 1:{StructDecl 2 0 {3 2} {0 10}} 2:{VarDecl 3 0 {5 1} {0 11}}
    		3:{FieldDecl 4 1 {0 0} {3 2}} 4:{FieldDecl 5 1 {0 0} {6 2}} 5:{StructDecl 2 2 {6 2} {0 10}}
    		6:{FieldDecl 4 5 {0 0} {3 2}} 7:{FieldDecl 5 5 {0 0} {6 2}}
    		CursorNameMap:
    		0: 1:sample.c 2:StructS 3:VarS 4:a 5:b
				`,
			expectedGobSize0: 744,
			expectedGobSize1: 425,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			input := run.Callbacks{
				Options: options,
				HdrCode: test.hdrCode,
				SrcCode: test.srcCode,
			}
			// Create an object that conforms to all four Callbacks of data collection.
			tr := struct {
				tokenStrings
				topCursorStrings
				fullCursorStrings
				tuParser
			}{}
			tr.topCursorStrings.topLevelNamesToSkip = topLevelNames
			tr.fullCursorStrings.topLevelNamesToSkip = topLevelNames

			err := input.LayerAndExecute(&tr) // collect the data
			if err != nil {
				t.Fatal(err)
			}

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
				got = trimends(got)

				expected := test.expectedTUPopulate
				expected = leftJustify(trimends(expected))

				if got != expected {
					t.Errorf("%s TranslationUnit:\n=== got ===\n%s\n=== expected ===\n%s\n",
						test.name, got, expected)
				}
			}

			if test.expectedGobSize0 >= 0 {

				got := first_gobtest(t, test.name+" first_gobtest", &tr.tuParser.TranslationUnit)

				_ = got
				/* Disable gob size check while reworking ast
				 */
				expected := test.expectedGobSize0

				if got != expected {
					t.Errorf("%s Gob test first encoding: got %d, expected %d\n",
						test.name, got, expected)
				}
			}

			if test.expectedGobSize1 >= 0 {

				got := second_gobtest(t, test.name+" second_gobtest", &tr.tuParser.TranslationUnit)

				_ = got
				/* Disable gob size check while reworking ast
				 */
				expected := test.expectedGobSize1

				if got != expected {
					t.Errorf("%s Gob test second encoding: got %d, expected %d\n",
						test.name, got, expected)
				}
			}
		})
	}
}
