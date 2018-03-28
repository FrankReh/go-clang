package clang_test

import (
	"fmt"
	//"io/ioutil"
	"strings"

	"github.com/frankreh/go-clang-v5.0/ast"
	"github.com/frankreh/go-clang-v5.0/clang"
	"github.com/frankreh/go-clang-v5.0/clang/cursorkind"
	//"github.com/frankreh/go-clang-v5.0/clang/tokenkind"
	"github.com/frankreh/go-clang-v5.0/clang/typekind"
)

// TBD move SourcesUnsavedFiles to clangrun package perhaps, or clangbridge.
type SourcesUnsavedFiles struct {
	unsavedFiles []clang.UnsavedFile
}

func (s *SourcesUnsavedFiles) Extract(file string, soffset, eoffset int) (string, error) {
	for i := range s.unsavedFiles {
		if s.unsavedFiles[i].Filename() != file {
			continue
		}
		contents := s.unsavedFiles[i].Contents()

		if int(eoffset) > len(contents) {
			return "", fmt.Errorf("Extract: eoffset greater than file size, %d, %d", eoffset, len(contents))
		}
		return contents[soffset:eoffset], nil
	}
	// Since the file wasn't found, build the list of filenames for the error.
	var f []string
	for i := range s.unsavedFiles {
		f = append(f,  s.unsavedFiles[i].Filename())
	}
	return "", fmt.Errorf("file not found in buffers: %s, %v", file, f)
}

func cursorPrint(depth int, cursor, parent clang.Cursor, sources ast.Sources) string {

	var printer Printer

	ismember := func(str string, list []string) bool {
		for _, s := range list {
			if s == str {
				return true
			}
		}
		return false
	}
	printbools := func(b bool, hdr string, strings []string) {
		if len(strings) > 0 {
			printer.printf(b, "%*s%s:", depth, "", hdr)
			for _, s := range strings {
				printer.printf(b, " %s", s)
			}

			printer.printf(b, "\n")
		}
	}
	printline := func(b bool, hdr, m string, nostrings ...string) {
		if len(m) > 0 && !ismember(m, nostrings) {
			printer.printf(b, "%*s%s %s\n", depth, "", hdr, m)
		}
	}
	/*
		printlineuint32 := func(hdr string, m uint32) {
			printline(hdr, fmt.Sprintf("%X", m))
		}
	*/
	printlineuint32filter0 := func(b bool, hdr string, m uint32) {
		if m > 0 {
			printline(b, hdr, fmt.Sprintf("%d", m))
		}
	}
	printlineint32 := func(b bool, hdr string, m int32) {
		if m >= 0 {
			printline(b, hdr, fmt.Sprintf("%d", m))
		}
	}
	printlineint64 := func(b bool, hdr string, m int64) {
		if m >= 0 {
			printline(b, hdr, fmt.Sprintf("%d", m))
		}
	}
	/*
		storageClassFilter := func(s string) string {
			switch s {
			case "SC=Invalid", "SC=None":
				return ""
			}
			return s
		}
	*/
	extratype := func(b bool, hdr string, type1, type2 clang.Type) {
		k2 := type2.Kind()
		if k2 != typekind.Invalid && !type1.Equal(type2) {
			printline(b, hdr, type2.Spelling())
			printline(b, hdr+".Kind()", k2.String())
		}
	}

	// Given a cursor that represents a documentable entity (e.g., declaration), return the associated parsed comment as a CXComment_FullComment AST node.
	// func (c Cursor) ParsedComment() Comment
	// func commentBools(c clang.Comment) []string

	ct := cursor.Type()
	kind := cursor.Kind()
	cursor_spelling := cursor.Spelling()
	cursor_kind_spelling := kind.String()

	printline(true, "cursor", cursor_spelling)
	printline(printk, "Kind()", cursor_kind_spelling)
	printline(printe, "Extent:", extentdescription(cursor))
	switch cursor_kind_spelling {
	case "FunctionDecl":
		// To see if function declaration is an extern, see if it has a child of type CompoundStmt.
		hasCompoundStmt := hasChildOfKind(cursor, cursorkind.CompoundStmt)
		if !hasCompoundStmt {
			printline(true, "cursor (function)", "has no CompoundStmt")
		}
	case "macro expansion":
		// cursor __clang_version__
		// Kind() macro expansion
		// cursor.Extent() __clang_version__
		if macroexpansions.Add(cursor, sources) == nil {
			if printk || printa {
				printline(true, "Macroexpansion added to position:", fmt.Sprintf("%d", macroexpansions.Len()-1))
			}
			cref := cursor.Referenced()
			cdef := cursor.Definition()
			var msg string
			if cref.Equal(cdef) {
				msg = "the same"
			} else {
				msg = "not the same"
			}
			printline(true, "Reference and Definition are", msg)
			printline(true, "Definition location:", extentdescription(cdef))
		}
	}

	if extentlimit > 0 {
		if s, err := extentsource(cursor, sources); err != nil {
			printline(true, "cursor.Extent() error: ", err.Error())
		} else {
			printline(true, "cursor.Extent()", s)
		}
	}
	printline(printu, "USR", cursor.USR())

	printline(printt, "Type()", ct.Spelling())
	printline(printt, "Type().Kind()", ct.Kind().String(), "Invalid")
	extratype(printtc, "Type().CanonicalType()", ct, ct.CanonicalType())
	extratype(printtp, "Type().PointeeType()", ct, ct.PointeeType())
	printline(printtd, "Type().Declaration()", ct.Declaration().Spelling())
	// For a function type:
	extratype(printtr, "Type().ResultType()", ct, ct.ResultType())
	printlineint32(printtn, "Type().NumArgTypes()", ct.NumArgTypes())
	extratype(printte, "Type().ElementType()", ct, ct.ElementType())
	printlineint64(printtne, "Type().NumElements()", ct.NumElements())
	extratype(printtae, "Type().ArrayElementType()", ct, ct.ArrayElementType())
	printlineint64(printtas, "Type().ArraySize()", ct.ArraySize())
	extratype(printtnt, "Type().NamedType()", ct, ct.NamedType())
	if sizeOf, err := ct.SizeOf(); err == nil {
		printlineint64(printts, "Type().SizeOf()", int64(sizeOf))
	}

	switch kind {
	case cursorkind.MacroDefinition, cursorkind.MacroExpansion:
	// For other kinds, the results or not meaningful and worse, inconsistent.
	default:
		printbools(printbc, "bools for cursor", cursorBools(cursor))
	}
	printbools(printbk, "bools for cursor.Kind()", cursorKindBools(cursor.Kind()))
	/*
		switch cursor.Kind() {
		case cursorkind.MacroDefinition, cursorkind.MacroExpansion, cursorkind.MacroInstantiation:
			printtokens(true, "macro tokens", cursor)
		}
	*/
	printbools(printbt, "bools for cursor.Type()", typeBools(cursor.Type()))
	printbools(printbl, "bools for cursor.Location()", sourceLocationBools(cursor.Location()))
	printbools(printbe, "bools for cursor.Extent()", souceRangeBools(cursor.Extent()))
	/*
		// Compute a hash value for the given cursor.
		printlineuint32("hashCursor()", cursor.HashCursor())
	*/
	// Determine the linkage of the entity referred to by a given cursor.
	// func (c Cursor) Linkage() LinkageKind
	printline(printl, "Linkage()", cursor.Linkage().String(), "Linkage_Invalid", "Linkage_External")

	// Describe the visibility of the entity referred to by a cursor.
	printline(printv, "Visibility()", cursor.Visibility().String(), "Visibility_Default", "Visibility_Invalid")

	//func (c Cursor) Availability() AvailabilityKind
	printline(printav, "(Availability)", cursor.Availability().String(), "Avilability_Available")

	// Determine the "language" of the entity referred to by a given cursor.
	// func (c Cursor) Language() LanguageKind
	printline(printla, "Language()", cursor.Language().String(), "Language_C", "Language_Invalid")

	if prints || printa {
		sp := cursor.SemanticParent()
		lp := cursor.LexicalParent()
		if !sp.Equal(lp) {
			printline(prints, "SemanticParent()", "is not equal to LexicalParent()")

		}
	}
	printtop := func(b bool, title string, cursor clang.Cursor) {
		if b {
			ct := cursor.Type()
			kind := cursor.Kind()
			cursor_spelling := cursor.Spelling()
			cursor_kind_spelling := kind.String()

			printline(true, title, cursor_spelling)
			printline(true, title+" Kind()", cursor_kind_spelling)

			printline(true, title+" Type()", ct.Spelling())
			printline(true, title+" Type().Kind()", ct.Kind().String(), "Invalid")
			extratype(true, title+" Type().CanonicalType()", ct, ct.CanonicalType())
			extratype(true, title+" Type().PointeeType()", ct, ct.PointeeType())
			printline(true, title+" Type().Declaration()", ct.Declaration().Spelling())
		}
	}
	_ = printtop
	if n, err := macroexpansions.Find(cursor, sources); err == nil {
		printline(true, "cursor match found at macroexpansion position:", fmt.Sprintf("%d", n))
	}

	// Retrieve the file that is included by the given inclusion directive cursor.
	// func (c Cursor) IncludedFile() File
	printline(printi, "includedFile()", cursor.IncludedFile().Name())

	// Retrieve the underlying type of a typedef declaration.
	printline(printtf, "TypedefDeclUnderlyingType()", cursor.TypedefDeclUnderlyingType().Spelling())

	// Retrieve the integer type of an enum declaration.
	printline(printed, "EnumDeclIntegerType()", cursor.EnumDeclIntegerType().Spelling())

	// Retrieve the bit width of a bit field declaration as an integer.
	// If a cursor that is not a bit field declaration is passed in, -1 is returned.
	printlineint32(printfbw, "FieldDeclBitWidth()", cursor.FieldDeclBitWidth())

	// Retrieve the number of non-variadic arguments associated with a given cursor.
	// The number of arguments can be determined for calls as well as for
	// declarations of functions or methods. For other cursors -1 is returned.
	printlineint32(printn, "NumArguments()", cursor.NumArguments())

	// Retrieve the return type associated with a given cursor.
	// This only returns a valid type if the cursor refers to a function or method.
	printline(printr, "ResultType()", cursor.ResultType().Spelling())

	// Return the offset of the field represented by the Cursor.
	if offsetOfField, err := cursor.OffsetOfField(); err == nil {
		printlineint64(printo, "OffsetOfField()", int64(offsetOfField))
	}

	// Returns the storage class for a function or variable declaration.
	// If the passed in Cursor is not a function or variable declaration,
	// CX_SC_Invalid is returned else the storage class.
	// func (c Cursor) StorageClass() StorageClass
	switch kind {
	// For other kinds, the results or not meaningful and worse, inconsistent.
	case cursorkind.FunctionDecl, cursorkind.VarDecl:
		printline(printsc, "StorageClass()", cursor.StorageClass().String())
	}

	// Determine the number of overloaded declarations referenced by a
	// CXCursor_OverloadedDeclRef cursor.
	// If it is not a CXCursor_OverloadedDeclRef cursor, returns 0.
	// func (c Cursor) NumOverloadedDecls() uint32
	printlineuint32filter0(printno, "NumOverloadedDecls()", cursor.NumOverloadedDecls())

	printlinediff := func(b bool, hdr string, cursor, c2 clang.Cursor) {
		if !c2.IsNull() && !c2.Equal(cursor) {
			printline(b, hdr, c2.Spelling()+" "+extentdescription(c2)+" kind: "+c2.Kind().String())
		}
	}
	// For a cursor that is a reference, retrieve a cursor representing the
	// entity that it references.
	// func (c Cursor) Referenced() Cursor
	printlinediff(printre, "Referenced()", cursor, cursor.Referenced())

	// For a cursor that is either a reference to or a declaration
	// of some entity, retrieve a cursor that describes the definition of
	// that entity.
	// func (c Cursor) Definition() Cursor
	if !cursor.Referenced().Equal(cursor.Definition()) {
		printlinediff(printd, "Definition() is different", cursor, cursor.Definition())
	}

	// Retrieve the canonical cursor corresponding to the given cursor.
	// In the C family of languages, many kinds of entities can be declared several
	// times within a single translation unit. For example, a structure type can
	// be forward-declared (possibly multiple times) and later defined:
	// func (c Cursor) CanonicalCursor() Cursor
	printlinediff(printc, "CanonicalCursor()", cursor, cursor.CanonicalCursor())
	printer.printf(true, "\n")

	return printer.buf.String()
}

// TBD later
/*
	Retrieve the integer value of an enum constant declaration as a signed
	long long.

	If the cursor does not reference an enum constant declaration, LLONG_MIN is returned.
	Since this is also potentially a valid constant value, the kind of the cursor
	must be verified before calling this function.
*/
// func (c Cursor) EnumConstantDeclValue() int64

/*
	Retrieve the integer value of an enum constant declaration as an unsigned
	long long.

	If the cursor does not reference an enum constant declaration, ULLONG_MAX is returned.
	Since this is also potentially a valid constant value, the kind of the cursor
	must be verified before calling this function.
*/
// func (c Cursor) EnumConstantDeclUnsignedValue() uint64

/*
	Retrieve the argument cursor of a function or method.

	The argument cursor can be determined for calls as well as for declarations
	of functions or methods. For other cursors and for invalid indices, an
	invalid cursor is returned.
*/
// func (c Cursor) Argument(i uint32) Cursor

/*
	Returns the access control level for the referenced object.

	If the cursor refers to a C++ declaration, its access control level within its
	parent scope is returned. Otherwise, if the cursor refers to a base specifier or
	access specifier, the specifier itself is returned.
*/
// func (c Cursor) AccessSpecifier() AccessSpecifier

/*
	Retrieve a cursor for one of the overloaded declarations referenced
	by a CXCursor_OverloadedDeclRef cursor.

	Parameter cursor The cursor whose overloaded declarations are being queried.

	Parameter index The zero-based index into the set of overloaded declarations in
	the cursor.

	Returns A cursor representing the declaration referenced by the given
	cursor at the specified index. If the cursor does not have an
	associated set of overloaded declarations, or if the index is out of bounds,
	returns clang_getNullCursor();
*/
// func (c Cursor) OverloadedDecl(index uint32) Cursor

/*
	Retrieve the display name for the entity referenced by this cursor.

	The display name contains extra information that helps identify the cursor,
	such as the parameters of a function or template or the arguments of a
	class template specialization.
*/
// func (c Cursor) DisplayName() string

// Given a cursor that represents a declaration, return the associated comment's source range. The range may include multiple consecutive comments with whitespace in between.
// func (c Cursor) CommentRange() SourceRange

// Given a cursor that represents a declaration, return the associated comment text, including comment markers.
// func (c Cursor) RawCommentText() string

// Given a cursor that represents a documentable entity (e.g., declaration), return the associated \paragraph; otherwise return the first paragraph.
// func (c Cursor) BriefCommentText() string

// Given a CXCursor_ModuleImportDecl cursor, return the associated module.
// func (c Cursor) Module() Module
/*
func (c Cursor) DefinitionSpellingAndExtent() (string, string, uint32, uint32, uint32, uint32) {
	var startBuf *C.char
	defer C.free(unsafe.Pointer(startBuf))
	var endBuf *C.char
	defer C.free(unsafe.Pointer(endBuf))
	var startLine C.uint
	var startColumn C.uint
	var endLine C.uint
	var endColumn C.uint

	C.clang_getDefinitionSpellingAndExtent(c.c, &startBuf, &endBuf, &startLine, &startColumn, &endLine, &endColumn)

	return C.GoString(startBuf), C.GoString(endBuf), uint32(startLine), uint32(startColumn), uint32(endLine), uint32(endColumn)
}*/

/*
	Retrieve a completion string for an arbitrary declaration or macro
	definition cursor.

	Parameter cursor The cursor to query.

	Returns A non-context-sensitive completion string for declaration and macro
	definition cursors, or NULL for other kinds of cursors.
*/
// func (c Cursor) CompletionString() CompletionString

// If cursor is a statement declaration tries to evaluate the statement and if its variable, tries to evaluate its initializer, into its corresponding type.
// func (c Cursor) Evaluate() EvalResult

/*
	Find references of a declaration in a specific file.

	Parameter cursor pointing to a declaration or a reference of one.

	Parameter file to search for references.

	Parameter visitor callback that will receive pairs of CXCursor/CXSourceRange for
	each reference found.
	The CXSourceRange will point inside the file; if the reference is inside
	a macro (and not a macro argument) the CXSourceRange will be invalid.

	Returns one of the CXResult enumerators.
*/
// func (c Cursor) FindReferencesInFile(file File, visitor CursorAndRangeVisitor) Result

// func (c Cursor) Xdata() int32
/* {
	return int32(c.c.xdata)
}*/

type SourceRange struct {
	file    string
	soffset uint32
	eoffset uint32
}

func (sr *SourceRange) String() string {
	return fmt.Sprintf("%s[%d:%d]", sr.file, sr.soffset, sr.eoffset)
}

// ExtractString() method returns a string pulled from the file given by the start and end offsets.
// This does not try to be efficient.  It is understood that the same file will be opened
// many times and the entire read into a buffer each time, even for the smallest extent ranges.
// Future optimizations can work on reading the file a single time and reading into a buffer
// that is kept along with the translation unit somewhere.  Perhaps a hash of the file string.
func (sr *SourceRange) ExtractString(sources ast.Sources) (string, error) {
	if sr.file == "" {
		return "", fmt.Errorf("ExtractString: sr.file is null")
	}
	if sr.soffset > sr.eoffset {
		return "", fmt.Errorf("ExtractString: soffset > eoffset, %d, %d", sr.soffset ,sr.eoffset)
	}
	s, err :=   sources.Extract(sr.file, int(sr.soffset), int(sr.eoffset))
	return   s, err

	/*
		buf, err := ioutil.ReadFile(sr.file)
		if err != nil {
			fmt.Printf("read file %s error: %v\n", sr.file, err)
			return ""
		}
		if sr.eoffset > uint32(len(buf)) {
			return "error: eoffset greater than file size"
		}
		return string(buf[sr.soffset:sr.eoffset])
	*/
}

func extentsourcerange(cursor clang.Cursor) *SourceRange {
	sr := cursor.Extent() // SourceRange
	ssl := sr.Start()     // start SourceLocation
	esl := sr.End()       // end SourceLocation

	//file, line, column, offset := ssl.ExpansionLocation()
	file, _, _, soffset := ssl.ExpansionLocation()
	_, _, _, eoffset := esl.ExpansionLocation()

	return &SourceRange{
		file:    file.Name(),
		soffset: soffset,
		eoffset: eoffset,
	}
}

func extentdescription(cursor clang.Cursor) string {
	sr := extentsourcerange(cursor)

	return sr.String()
}
func extentsource(cursor clang.Cursor, sources ast.Sources) (string, error) {
	sr := extentsourcerange(cursor)

	return sr.ExtractString(sources)
}

type MacroExpansion struct {
	cursor clang.Cursor
	extent string
}
type MacroExpansions struct {
	me []MacroExpansion
}

var macroexpansions MacroExpansions

func (m *MacroExpansions) Add(cursor clang.Cursor, sources ast.Sources) error {
	// Only add if it would introduce a new extent since the only thing being
	// matched up so far is the extent.
	extent, err := extentsource(cursor, sources)
	if  err != nil {
		return  err
	}
	if _, err := m.find(extent); err != nil {
		m.me = append(m.me,
			MacroExpansion{
				cursor: cursor,
				extent: extent,
			})
		return nil
	} else {
		return ErrExists
	}
}
func (m *MacroExpansions) Len() int {
	return len(m.me)
}
func (m *MacroExpansions) find(extent string) (int, error) {
	if len(extent) > 0 {
		for i := range m.me {
			if extent == m.me[i].extent {
				return i, nil
			}
		}
	}
	return -1, ErrNotFound
}
func (m *MacroExpansions) Find(cursor clang.Cursor, sources ast.Sources) (int, error) {
	if !cursor.Location().IsFromMainFile() {
		s, err := extentsource(cursor, sources)
		if  err != nil {
			return -1, err
		}
		return m.find(s)
	}
	return -1, ErrNotFound
}

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNotFound = Error("not found")
	ErrExists   = Error("exists")
)

type Printer struct {
	buf strings.Builder
}

func (p *Printer) printf(b bool, format string, a ...interface{}) (n int, err error) {
	if b || printa {
		return fmt.Fprintf(&p.buf, format, a...)
	}
	return 0, nil
}

var printa = false   // flag.Bool("a", false, "print all")
var printav = true   // flag.Bool("av", false, "print cursor.Availability()")
var printbc = true   // flag.Bool("bc", false, "print bools for cursor")
var printbe = true   // flag.Bool("be", false, "print bools for cursor.Extent()")
var printbk = true   // flag.Bool("bk", false, "print bools for cursor.Kind()")
var printbl = true   // flag.Bool("bl", false, "print bools for cursor.Location()")
var printbt = true   // flag.Bool("bt", false, "print bools for cursor.Type()")
var printc = true    // flag.Bool("c", false, "print cursor.CanonicalCursor()")
var printd = true    // flag.Bool("d", false, "print cursor.Definition()")
var printe = true    // flag.Bool("e", false, "print cursor.Extent()")
var extentlimit = 30 // flag.Int("extentlimit", 0, "limit which cursor.Extent() sources to print")
var printed = true   // flag.Bool("ed", false, "print cursor.EnumDeclIntegerType()")
var printfbw = true  // flag.Bool("fbw", false, "print cursor.FieldDeclBitWidth()")
var printh = true    // flag.Bool("h", false, "print header")
var printi = true    // flag.Bool("i", false, "print cursor.IncludedFile()")
var printk = true    // flag.Bool("k", false, "print cursor.Kind()")
var printl = true    // flag.Bool("l", false, "print cursor.Linkage()")
var printla = true   // flag.Bool("la", false, "print cursor.Language()")
var printn = true    // flag.Bool("n", false, "print cursor.NumArguments()")
var printno = true   // flag.Bool("no", false, "print cursor.NumOverloadedDecls()")
var printo = true    // flag.Bool("o", false, "print cursor.OffsetOfField()")
var prints = false   // flag.Bool("s", false, "print cursor.SemanticParent() check against LexicalParent()")
var printsc = true   // flag.Bool("sc", false, "print cursor.StorageClass()")
var printr = true    // flag.Bool("r", false, "print cursor.ResultType()")
var printre = true   // flag.Bool("re", false, "print cursor.Referenced()")
var printt = true    // flag.Bool("t", false, "print cursor.Type() and cursor.Type().Kind()")
var printtae = true  // flag.Bool("tae", false, "print cursor.Type().ArrayElementType()")
var printtas = true  // flag.Bool("tas", false, "print cursor.Type().ArraySize()")
var printtc = true   // flag.Bool("tc", false, "print cursor.Type().CanonicalType()")
var printtd = true   // flag.Bool("td", false, "print cursor.Type().Declaration()")
var printte = true   // flag.Bool("te", false, "print cursor.Type().ElementType()")
var printtf = true   // flag.Bool("tf", false, "print cursor.TypedefDeclUnderlyingType()")
var printtn = true   // flag.Bool("tn", false, "print cursor.Type().NumArgTypes()")
var printtne = true  // flag.Bool("tne", false, "print cursor.Type().NumElements()")
var printtnt = true  // flag.Bool("tnt", false, "print cursor.Type().NamedType()")
var printtp = true   // flag.Bool("tp", false, "print cursor.Type().PointeeType()")
var printtr = true   // flag.Bool("tr", false, "print cursor.Type().ResultType()")
var printts = true   // flag.Bool("ts", false, "print cursor.Type().SizeOf()")
var printu = true    // flag.Bool("u", false, "print cursor.USR()")
var printv = true    // flag.Bool("v", false, "print cursor.Visibility()")

func hasChildOfKind(cursor clang.Cursor, ckind cursorkind.Kind) bool {
	has := false
	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {

		if cursor.Kind() == ckind {
			has = true
			return clang.ChildVisit_Break
		}

		return clang.ChildVisit_Continue
	})

	return has
}

func commentBools(c clang.Comment) []string {
	var r []string
	test := func(b bool, s string) {
		if b {
			r = append(r, s)
		}
	}
	test(c.IsWhitespace(), "IsWhitespace")
	test(c.InlineContentComment_HasTrailingNewline(), "InlineContentComment_HasTrailingNewline")
	test(c.HTMLStartTagComment_IsSelfClosing(), "HTMLStartTagComment_IsSelfClosing")
	test(c.ParamCommandComment_IsParamIndexValid(), "ParamCommandComment_IsParamIndexValid")
	test(c.ParamCommandComment_IsDirectionExplicit(), "ParamCommandComment_IsDirectionExplicit")
	test(c.TParamCommandComment_IsParamPositionValid(), "TParamCommandComment_IsParamPositionValid")
	return r
}
func cursorBools(c clang.Cursor) []string {
	var r []string
	test := func(b bool, s string) {
		if b {
			r = append(r, s)
		}
	}
	test(c.IsNull(), "IsNull")
	test(c.HasAttrs(), "HasAttrs")
	test(c.IsMacroFunctionLike(), "IsMacroFunctionLike")
	test(c.IsMacroBuiltin(), "IsMacroBuiltin")
	test(c.IsFunctionInlined(), "IsFunctionInlined")
	test(c.IsAnonymous(), "IsAnonymous")
	test(c.IsBitField(), "IsBitField")
	test(c.IsVirtualBase(), "IsVirtualBase")
	test(c.IsCursorDefinition(), "IsCursorDefinition")
	test(c.IsDynamicCall(), "IsDynamicCall")
	test(c.IsObjCOptional(), "IsObjCOptional")
	test(c.IsVariadic(), "IsVariadic")
	test(c.CXXConstructor_IsConvertingConstructor(), "CXXConstructor_IsConvertingConstructor")
	test(c.CXXConstructor_IsCopyConstructor(), "CXXConstructor_IsCopyConstructor")
	test(c.CXXConstructor_IsDefaultConstructor(), "CXXConstructor_IsDefaultConstructor")
	test(c.CXXConstructor_IsMoveConstructor(), "CXXConstructor_IsMoveConstructor")
	test(c.CXXField_IsMutable(), "CXXField_IsMutable")
	test(c.CXXMethod_IsDefaulted(), "CXXMethod_IsDefaulted")
	test(c.CXXMethod_IsPureVirtual(), "CXXMethod_IsPureVirtual")
	test(c.CXXMethod_IsStatic(), "CXXMethod_IsStatic")
	test(c.CXXMethod_IsVirtual(), "CXXMethod_IsVirtual")
	test(c.CXXMethod_IsConst(), "CXXMethod_IsConst")
	return r
}
func cursorKindBools(ck cursorkind.Kind) []string {
	var r []string
	test := func(b bool, s string) {
		if b {
			r = append(r, s)
		}
	}
	test(ck.IsDeclaration(), "IsDeclaration")
	test(ck.IsReference(), "IsReference")
	test(ck.IsExpression(), "IsExpression")
	test(ck.IsStatement(), "IsStatement")
	test(ck.IsAttribute(), "IsAttribute")
	test(ck.IsInvalid(), "IsInvalid")
	test(ck.IsTranslationUnit(), "IsTranslationUnit")
	test(ck.IsPreprocessing(), "IsPreprocessing")
	test(ck.IsUnexposed(), "IsUnexposed")
	return r
}
func idxDeclInfoBools(idi clang.IdxDeclInfo) []string {
	var r []string
	test := func(b bool, s string) {
		if b {
			r = append(r, s)
		}
	}
	test(idi.IsRedeclaration(), "IsRedeclaration")
	test(idi.IsDefinition(), "IsDefinition")
	test(idi.IsContainer(), "IsContainer")
	test(idi.IsImplicit(), "IsImplicit")
	return r
}
func idxImportedASTFileInfoBools(iiastfi clang.IdxImportedASTFileInfo) []string {
	var r []string
	test := func(b bool, s string) {
		if b {
			r = append(r, s)
		}
	}
	test(iiastfi.IsImplicit(), "IsImplicit")
	return r
}
func idxIncludedFileInfoBools(iifi clang.IdxIncludedFileInfo) []string {
	var r []string
	test := func(b bool, s string) {
		if b {
			r = append(r, s)
		}
	}
	test(iifi.IsImport(), "IsImport")
	test(iifi.IsAngled(), "IsAngled")
	test(iifi.IsModuleImport(), "IsModuleImport")
	return r
}
func moduleBools(m clang.Module) []string {
	var r []string
	test := func(b bool, s string) {
		if b {
			r = append(r, s)
		}
	}
	test(m.IsSystem(), "IsSystem")
	return r
}
func sourceLocationBools(sl clang.SourceLocation) []string {
	var r []string
	test := func(b bool, s string) {
		if b {
			r = append(r, s)
		}
	}
	test(sl.IsInSystemHeader(), "IsInSystemHeader")
	test(!sl.IsFromMainFile(), "! IsFromMainFile")
	return r
}
func souceRangeBools(sr clang.SourceRange) []string {
	var r []string
	test := func(b bool, s string) {
		if b {
			r = append(r, s)
		}
	}
	test(sr.IsNull(), "IsNull")
	return r
}
func typeBools(t clang.Type) []string {
	var r []string
	test := func(b bool, s string) {
		if b {
			r = append(r, s)
		}
	}
	test(t.IsConstQualifiedType(), "IsConstQualifiedType")
	test(t.IsVolatileQualifiedType(), "IsVolatileQualifiedType")
	test(t.IsRestrictQualifiedType(), "IsRestrictQualifiedType")
	test(t.IsFunctionTypeVariadic(), "IsFunctionTypeVariadic")
	test(!t.IsPODType(), "! IsPODType")
	return r
}
