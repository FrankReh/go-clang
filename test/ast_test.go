package clang_test

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"text/template"

	"github.com/frankreh/go-clang-v5.0/ast"
	"github.com/frankreh/go-clang-v5.0/astbridge"
	"github.com/frankreh/go-clang-v5.0/clang"
	"github.com/frankreh/go-clang-v5.0/clang/cursorkind"
	//"github.com/frankreh/go-clang-v5.0/clang/tokenkind"
	//"github.com/frankreh/go-clang-v5.0/clang/typekind"
	run "github.com/frankreh/go-clang-v5.0/clangrun"
)

var updateFlag = flag.Bool("update", false, "update the table_test.got file.")

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

//-- 2.
// topCursorStrings implements run.TopCursorVisiter and collects cursorString results.
type topCursorStrings struct {
	topLevelNamesToSkip map[string]bool
	sources             ast.Sources

	hdrs []string
	list []string
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

//-- 4.
// tuParser implements run.TUParser and collects cursorString results.
type tuParser struct {
	topLevelNamesToSkip map[string]bool
	sources             ast.Sources

	ast.TranslationUnit // The non clang version of that gets populated from the clang version.
	err                 error
}

// tokenVisit implements TUParser, collecting results of calls to cursorString.
func (x *tuParser) TUParse(tu *clang.TranslationUnit) {
	ctu := astbridge.ClangTranslationUnit{}

	x.err = ctu.Populate(tu, x.topLevelNamesToSkip)
	if x.err != nil {
		return
	}

	// Save the Go version of the TranslationUnit.
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

// compilerTopLevelNames returns a map of all top level names that
// the compiler gives us by default. This is later used to ignore
// those same names when having the compiler work on our test code.
func compilerTopLevelNames() (map[string]bool, error) {
	r := make(map[string]bool)

	// Create an essentially blank source code buffer to compile
	// and record the cursor names that are encountered by the
	// top cursor visit routine.
	callbacks := run.Callbacks{
		Options: clang.TranslationUnit_DetailedPreprocessingRecord,
		// Non empty so the cursor for "hdr.h" is also included
		UnsavedFiles: run.BuildUnsavedFiles(" ", " "),
	}
	// Only one callback function needed for this collection.
	callbacks.AppendTopCursorFn(
		func(tu clang.TranslationUnit, cursor, parent clang.Cursor) {
			r[cursor.Spelling()] = true
		})

	err := callbacks.Execute() // collect the data
	return r, err
}

type testTuple struct {
	Name                string
	Disable             bool
	Comment             string
	Options             clang.TranslationUnit_Flags // May be overridden in test code anyway.
	HdrCode             string
	SrcCode             string
	ExpectedTokens      string
	ExpectedTopCursors  string
	ExpectedFullCursors string
	ExpectedTUPopulate  string
	ExpectedGobSize0    int
	ExpectedGobSize1    int
}

func equal(a, b []testTuple) error {
	if len(a) != len(b) {
		return fmt.Errorf("lengths not equal")
	}
	fmt.Println("lengths", len(a), len(b))
	for i, t := range a {
		if b[i] != t {
			return fmt.Errorf("index %d not equal", i)
		}
	}
	return nil
}

// formatBytes returns the gofmt-ed contents of the buffer.
// (Taken from stringer).
func formatBytes(t *testing.T, b []byte) []byte {
	src, err := format.Source(b)
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		t.Errorf("warning: internal error: invalid Go generated: %s, compile the package to analyze the error", err)
		return b
	}
	return src
}

// prettyprint converts multi-line strings to be indented and to start with a newline
// so that they can be used as raw go strings in generated code and look pretty.
func prettyprint(s string) string {
	assertNoBackquotes(s)

	lines := strings.Split(s, "\n")
	if len(lines) <= 1 {
		return s
	}

	indent := "            "
	// Add spaces to beginning of each nonempty line.
	for i := range lines {
		line := lines[i]
		if line == "" {
			continue
		}
		lines[i] = indent + line
	}

	// Return new string wrapped by extra newlines.
	return "\n" + strings.Join(lines, "\n")
}

// assertNoBackquotes calls panic if a backquote (`) is found in the string.
// Backquotes are used in the generation of the code for multi-line strings
// and is not designed to special case multi-line strings that are to have
// backquotes in them.
func assertNoBackquotes(s string) {
	if strings.Index(s, "`") != -1 {
		panic("This code generation mechanism is not yet designed to handle backquotes inside strings being pretty printed.")
	}
}

// templ is the template used to generate the go code that drives this table test.
var templ = template.Must(template.New("code").Parse(code))

// Getting backquotes into the template so strings can be output as raw strings
// is just a little tricky. Here, the backquotes are put into a template that
// is created using double quotes, and while the template itself needs to use
// double quotes within the string, those can be backslashed. A raw string
// can't backslash a backquote so this two step process is used. There are
// countless other ways this could have been done too.

const code = "{{define \"bq\"}}`{{.}}`{{end}}" + `// Generated by go test -update.
// Edit by hand to add or modify test cases.  Use go test -update to get this
// file regenerated with the expected output of those test cases.
//
// Test cases that you comment out by hand will be lost if the -update flag is
// used.  Rather than comment out the test cases, set their Disable property to
// true if you may want to update the results. The Disable property is
// maintained.
//
// Similarly, any comments you make for yourself while testing will be lost if
// you update this file. But there is a string property named Comment that you
// can use that is retained during update.

package clang_test

import (
	"github.com/frankreh/go-clang-v5.0/clang"
)

var testTupleData = []testTuple{
{{range .}}
    {
        Name: "{{.Name}}",
{{if .Disable                   }}        Disable: {{.Disable}},
{{end}}{{if .Comment            }}        Comment: {{template "bq" .Comment}},
{{end}}{{if .Options            }}        Options: clang.{{.Options}},
{{end}}{{if .HdrCode            }}        HdrCode: {{template "bq" .HdrCode}},
{{end}}{{if .SrcCode            }}        SrcCode: {{template "bq" .SrcCode}},
{{end}}{{if .ExpectedTokens     }}        ExpectedTokens: {{template "bq" .ExpectedTokens}},
{{end}}{{if .ExpectedTopCursors }}        ExpectedTopCursors: {{template "bq" .ExpectedTopCursors}},
{{end}}{{if .ExpectedFullCursors}}        ExpectedFullCursors: {{template "bq" .ExpectedFullCursors}},
{{end}}{{if .ExpectedTUPopulate }}        ExpectedTUPopulate: {{template "bq" .ExpectedTUPopulate}},
{{end}}{{if .ExpectedGobSize0   }}        ExpectedGobSize0: {{.ExpectedGobSize0}},
{{end}}{{if .ExpectedGobSize1   }}        ExpectedGobSize1: {{.ExpectedGobSize1}},
{{end}}    },{{end}}
}
`

func TestAst(t *testing.T) {

	options := clang.TranslationUnit_DetailedPreprocessingRecord
	topLevelNames, err := compilerTopLevelNames()
	if err != nil {
		t.Fatal(err)
	}

	createGot := true

	for index := range testTupleData {
		// Use pointer because we update the entry in place with new got values.
		test := &testTupleData[index]
		if test.Disable {
			continue
		}

		t.Run(test.Name, func(t *testing.T) {

			unsavedFiles := run.BuildUnsavedFiles(test.HdrCode, test.SrcCode)
			sources := &SourcesUnsavedFiles{unsavedFiles}

			input := run.Callbacks{
				Options:      options,
				UnsavedFiles: unsavedFiles,
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
			tr.tuParser.topLevelNamesToSkip = topLevelNames

			tr.topCursorStrings.sources = sources
			tr.fullCursorStrings.sources = sources
			tr.tuParser.sources = sources

			err := input.LayerAndExecute(&tr) // collect the data
			if err != nil {
				t.Fatal(err)
			}
			if tr.tuParser.err != nil {
				t.Fatal(tr.tuParser.err)
			}

			if test.ExpectedTokens != "" {
				got := alignOn(tr.tokenStrings.list, ":")
				got = leftJustify(trimends(got))
				if createGot {
					test.ExpectedTokens = prettyprint(got)
				} else {
					expected := test.ExpectedTokens
					expected = leftJustify(trimends(expected))

					if got != expected {
						t.Errorf("%s:\n=== got tokens ===\n%s\n=== expected tokens ===\n%s\n",
							test.Name, got, expected)
					}
				}
			}

			if test.ExpectedTopCursors != "" {
				got := strings.Join(tr.topCursorStrings.list, "\n")
				got = leftJustify(alignOnTabs(trimends(got)))
				if createGot {
					test.ExpectedTopCursors = prettyprint(got)
				} else {

					expected := test.ExpectedTopCursors
					expected = leftJustify(trimends(expected))

					if got != expected {
						t.Errorf("%s Top cursors:\n=== got ===\n%s\n=== expected ===\n%s\n",
							test.Name, got, expected)
					}
				}
			}

			if test.ExpectedFullCursors != "" {
				got := strings.Join(tr.fullCursorStrings.list, "\n")
				got = leftJustify(trimends(got))
				if createGot {
					test.ExpectedFullCursors = prettyprint(got)
				} else {
					expected := test.ExpectedFullCursors
					expected = leftJustify(trimends(expected))

					if got != expected {
						t.Errorf("%s full cursors:\n=== got ===\n%s\n=== expected ===\n%s\n",
							test.Name, got, expected)
					}
				}
			}

			if test.ExpectedTUPopulate != "" {
				got := fmt.Sprintf("%#v", tr.tuParser.TranslationUnit)
				got = trimends(got)
				if createGot {
					test.ExpectedTUPopulate = prettyprint(got)
				} else {
					expected := test.ExpectedTUPopulate
					expected = leftJustify(trimends(expected))

					if got != expected {
						t.Errorf("%s TranslationUnit:\n=== got ===\n%s\n=== expected ===\n%s\n",
							test.Name, got, expected)
					}
				}
			}

			if test.ExpectedGobSize0 > 0 {
				got := first_gobtest(t, test.Name+" first_gobtest", &tr.tuParser.TranslationUnit)
				if createGot {
					test.ExpectedGobSize0 = got
				} else {

					expected := test.ExpectedGobSize0

					if got != expected {
						t.Errorf("%s Gob test first encoding: got %d, expected %d\n",
							test.Name, got, expected)
					}
				}
			}

			if test.ExpectedGobSize1 > 0 {
				got := second_gobtest(t, test.Name+" second_gobtest", &tr.tuParser.TranslationUnit)
				if createGot {
					test.ExpectedGobSize1 = got
				} else {
					expected := test.ExpectedGobSize1

					if got != expected {
						t.Errorf("%s Gob test second encoding: got %d, expected %d\n",
							test.Name, got, expected)
					}
				}
			}
		})
	}

	if createGot {
		golden := "table_test.go"
		gotfilename := "table_test.got"

		var buf bytes.Buffer
		if err := templ.Execute(&buf, testTupleData); err != nil {
			t.Fatal(err)
		}
		gotbytes := formatBytes(t, buf.Bytes())

		// Don't create the got file unless it is going to be different from the golden file.
		// So read the golden file first. Rather a file read every time, than a file write.
		goldenbytes, err := ioutil.ReadFile(golden)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(goldenbytes, gotbytes) {

			if err := ioutil.WriteFile("table_test.got", gotbytes, 0666); err != nil {
				t.Fatal(err)
			}
			// Compare got file with original.
			var cmd *exec.Cmd
			switch runtime.GOOS {
			case "plan9":
				cmd = exec.Command("/bin/diff", "-c", golden, gotfilename)
			default:
				cmd = exec.Command("/usr/bin/diff", "-u", golden, gotfilename)
			}
			cmdbuf := new(bytes.Buffer)
			cmd.Stdout = cmdbuf
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				if *updateFlag {
					t.Logf("diff test between %s and %s shows: %s\n%s\n",
						golden, gotfilename, err, cmdbuf)
					t.Logf("Updating %s...", golden)
					if err := exec.Command("/bin/cp", gotfilename, golden).Run(); err != nil {
						t.Errorf("Update failed: %s", err)
					}
				} else {
					t.Errorf("diff test between %s and %s failed: %s\n%s\n",
						golden, gotfilename, err, cmdbuf)
				}
			}
		} else {
			t.Logf("golden bytes equal got bytes")
		}
	}
}

// Rule: Do Unexposed cursors always have a single child?
// Rule: Do Back cursors (which I have created) never have a sibling?
// Rule: Does every BinaryOperator cursor have a token with kind:Punctuation?
// expectedFullCursors: could have BinaryOperator names fill in for this. maybe.
// Could check that there are no warnings or errors generated from compiler.
// Start showing the Type and TypeKind for each cursor.
