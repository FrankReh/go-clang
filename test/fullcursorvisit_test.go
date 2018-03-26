package clang_test

import (
	"fmt"
	"strings"

	"github.com/frankreh/go-clang-v5.0/clang"
	"github.com/frankreh/go-clang-v5.0/clang/tokenkind"
	"github.com/frankreh/go-clang-v5.0/clang/typekind"
	run "github.com/frankreh/go-clang-v5.0/clangrun"
)

//-- 3.
// fullCursorStrings implements run.FullCursorVisiter and collects cursorString results.
type fullCursorStrings struct {
	topLevelNamesToSkip map[string]bool
	depthMap            map[clang.Cursor]int // 0 When parent is root, -1 when parent is one we are skipping.
	seenCursors         map[clang.Cursor]int // Cardinal order the cursor was already seen, not descended further.
	pad                 string
	list                []string

	seenTypes map[clang.Type]int // Ordinal number of the type already seen, starting at 1.
}

// FullCursorVisit implements run.FullCursorVisiter, collecting results of calls to cursorString.
func (x *fullCursorStrings) FullCursorVisit(tu clang.TranslationUnit, cursor, parent clang.Cursor) {
	// parent hash is not expected to be in the map so fact zero is returned is useful.
	if x.depthMap == nil {
		x.depthMap = make(map[clang.Cursor]int)
		x.seenCursors = make(map[clang.Cursor]int) // Init all at once.
		x.seenTypes = make(map[clang.Type]int)
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

	// Cursor Description

	offset, found := x.seenCursors[cursor]
	if found {
		s += fmt.Sprintf("[%d backreference]", offset)
		x.depthMap[cursor] = -1 // Keep from descending further.
		x.list = append(x.list, s)
		return
	}
	x.seenCursors[cursor] = len(x.list)

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

	// Type Description

	seenMsg := func(ctype clang.Type) (ordinal string, firsttime bool) {
		i := x.seenTypes[ctype] // seen ordinal
		if i == 0 {
			i = len(x.seenTypes) + 1
			x.seenTypes[ctype] = i
			firsttime = true
		}
		ordinal = fmt.Sprintf("%d", i)
		return
	}

	// Add Type info, done recursively so declare func first.
	var typeDesc func(ctype clang.Type) string
	typeDesc = func(ctype clang.Type) string {

		tkind := ctype.Kind()

		// If Invalid type, just return "". Make sure cursor was Statement first.
		// If it was anything else, we want to figure out the type.
		if typekind.Invalid == tkind {
			if !kind.IsStatement() {
				return "{Invalid - but not Statement}"
			}
			return ""
		}

		b := new(strings.Builder)
		fmt.Fprintf(b, "{")

		seenmsg, firsttime := seenMsg(ctype)

		if !firsttime {
			fmt.Fprintf(b, "seen-before:%s", seenmsg)
			fmt.Fprintf(b, "}")
			return b.String()
		}
		fmt.Fprintf(b, "first-seen:%s", seenmsg)

		fmt.Fprintf(b, " %s '%s'", tkind, ctype.Spelling())

		if !ctype.IsPODType() {
			fmt.Fprintf(b, " !POD")
		}

		if canonical := ctype.CanonicalType(); !canonical.Equal(ctype) {
			fmt.Fprintf(b, " Canon:%s", typeDesc(canonical))
		}

		/* need a new callingconv package first
		if callingconv := ctype.FunctionTypeCallingConv(); !canonical {
			fmt.Fprintf(b, " CallingConv:%s", callingconv)
		}
		*/

		// Seems same as Spelling()
		//if typedefname := ctype.TypedefName(); typedefname != "" {
		//	fmt.Fprintf(b, "typedefname:%s", typedefname)
		//}

		if numargs := ctype.NumArgTypes(); numargs != -1 {
			fmt.Fprintf(b, " numargs:%d", numargs)
		}

		if resulttype := ctype.ResultType(); resulttype.Kind() != typekind.Invalid {
			fmt.Fprintf(b, " result:%s", typeDesc(resulttype))
		}

		// Seems redundant to canonical
		//if namedType := ctype.NamedType(); namedType.Kind() != typekind.Invalid {
		//	fmt.Fprintf(b, " namedType:%s", typeDesc(namedType))
		//}

		if pointeetype := ctype.PointeeType(); pointeetype.Kind() != typekind.Invalid {
			fmt.Fprintf(b, " *%s", typeDesc(pointeetype))
		}

		if alignof, err := ctype.AlignOf(); err == nil {
			fmt.Fprintf(b, " align:%d", alignof)
		}

		if sizeof, err := ctype.SizeOf(); err == nil {
			fmt.Fprintf(b, " size:%d", sizeof)
		}

		// Some bools
		if ctype.IsConstQualifiedType() {
			fmt.Fprintf(b, " Const")
		}
		if ctype.IsVolatileQualifiedType() {
			fmt.Fprintf(b, " Volatile")
		}
		if ctype.IsRestrictQualifiedType() {
			fmt.Fprintf(b, " Restrict")
		}
		if ctype.IsFunctionTypeVariadic() {
			fmt.Fprintf(b, " Variadic")
		}
		if ctype.IsTransparentTagTypedef() {
			fmt.Fprintf(b, " TransparentTagTypedef")
		}

		if numelem := ctype.NumElements(); numelem != -1 {
			fmt.Fprintf(b, " len:%d", numelem)

			// ArraySize seems redundant.
			if arraysz := ctype.ArraySize(); arraysz != numelem {
				fmt.Fprintf(b, " arraysize:%d", arraysz)
			}

		}

		if elemtype := ctype.ElementType(); elemtype.Kind() != typekind.Invalid {
			// ContantArray, IncompleteArray, VariableArray, DependentSizedArray
			fmt.Fprintf(b, " elem:%s", typeDesc(elemtype))
		}

		fmt.Fprintf(b, "}")
		return b.String()
	}
	s += " " + typeDesc(cursor.Type())
	/* TBD For investigation purposes.
	// Do the same cursors appears in multiple parts of the tree?
	// Yes they do.
	s += fmt.Sprintf(" %v %v", cursor.HashCursor(), cursor)
	*/

	x.list = append(x.list, s)
}
func init() {
	// assert it implements the desired interface.
	var a interface{} = &fullCursorStrings{}
	_ = a.(run.FullCursorVisiter)
}

// tokenDescriptions return strings describing tokens, format "Spelling : Kind".
func tokenDescriptions(tu clang.TranslationUnit, tokens []clang.Token) []string {
	var r []string
	for _, token := range tokens {
		r = append(r, fmt.Sprintf("%s : %s", tu.TokenSpelling(token), token.Kind().String()))
	}
	return r
}
