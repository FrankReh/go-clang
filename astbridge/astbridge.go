// Package astbridge bridges the cgo functionality of the clang package
// with the ast package structures.
//
// The ast structures do not rely on clang and can therefore be cgo free.
// The primary ast structures are also pointer free so can be serialized.
package astbridge

import (
	"errors"
	"fmt"
	"strings"

	"github.com/frankreh/go-clang-v5.0/ast"
	"github.com/frankreh/go-clang-v5.0/clang"
	"github.com/frankreh/go-clang-v5.0/clang/cursorkind"
	"github.com/frankreh/go-clang-v5.0/clang/typekind"
)

// ClangTranslationUnit references the clang package components.
type ClangTranslationUnit struct {
	GoTu          ast.TranslationUnit
	ClangTu       *clang.TranslationUnit
	ClangTokens   []clang.Token
	ClangCursors  []clang.Cursor
	mapTokenIndex map[clang.SourceLocation]int
	typeIndexes   map[clang.Type]int // TypeIndex for those types already created.
}

// Tokenize
//		ensures GoTo slice and map structures used for tokens are initialized,
//		calls Tokenize on the cursor's Extent,
//		can return a range of already existing toeksn if the first token is found in the map,
//		otherwise appends the new clang tokens to those being tracked,
//		creates GoTo Tokens from the Token kind and spelling,
//		adds them to the TokenMap, which may recognize them as already in the map,
//		appends the GoTo Token indexes to the slice of tokens,
//		and puts all the Token source locations into the map to allow
//		subsequent calls to possibly find their subrange of tokens is already in
//		the slice.
func (ctu *ClangTranslationUnit) Tokenize(cursor clang.Cursor) ast.IndexPair {
	// Ensure setup of the slice and the map.
	if ctu.GoTu.TokenIds == nil {
		ctu.GoTu.TokenIds = make([]ast.TokenId, 0, 100) // cap size is arbitrary
	}
	s := &ctu.GoTu.TokenIds

	if ctu.mapTokenIndex == nil {
		ctu.mapTokenIndex = make(map[clang.SourceLocation]int)
	}
	m := ctu.mapTokenIndex

	// Process tokens from cursor. Most times they will already have been seen.

	var tokenRange ast.IndexPair
	clangTokens := ctu.ClangTu.Tokenize(cursor.Extent())

	if len(clangTokens) == 0 {
		return tokenRange
	}

	tokenRange.Len = len(clangTokens)

	if index, ok := m[ctu.ClangTu.TokenLocation(clangTokens[0])]; ok {
		tokenRange.Head = index
		return tokenRange
	}
	tokenRange.Head = len(m) // len tracks length of the slice and the map.

	// The location of the first token wasn't found in the map, so all will be appended
	// to the ClangTokens, and the slice and map.
	ctu.ClangTokens = append(ctu.ClangTokens, clangTokens...)

	// convert clang Tokens to ast tokens and appends them to ctu.GoTu.TokenIds.

	for _, clangToken := range clangTokens {
		tokenSpelling := ctu.ClangTu.TokenSpelling(clangToken)
		token := ast.Token{
			TokenKindId: clangToken.Kind(),
			TokenNameId: ctu.GoTu.TokenNameMap.Id(tokenSpelling),
		}
		tokenId := ctu.GoTu.TokenMap.Id(token)
		*s = append(*s, tokenId)

		clangSourceLocation := ctu.ClangTu.TokenLocation(clangToken)
		m[clangSourceLocation] = len(m) // len tracks length of the slice and the map.
	}

	return tokenRange
}

func (ctu *ClangTranslationUnit) Populate(tu *clang.TranslationUnit, topLevelNamesToSkip map[string]bool) error {
	if ctu.ClangTu != nil {
		return errors.New("Already populated")
	}
	ctu.typeIndexes = make(map[clang.Type]int)

	// For some tidyness, have the "" string map to the 0 id.
	_ = ctu.GoTu.CursorNameMap.Id("")

	ctu.ClangTu = tu
	clangRootCursor := tu.TranslationUnitCursor()

	if err := diagnosticError(tu); err != nil {
		return err
	}

	ctu.Tokenize(clangRootCursor)

	ctu.GoTu.TypeMap.Init()

	ctu.GoTu.Back = make(map[int]int)
	ctu.GoTu.Referenced = make(map[int]int)
	ctu.GoTu.Definition = make(map[int]int)

	// Layer children to end of list, one set of children at a time.

	// Seed list with the root.
	ctu.ClangCursors = append(ctu.ClangCursors, clangRootCursor)
	ctu.GoTu.Cursors = append(ctu.GoTu.Cursors, ast.Cursor{
		CursorKindId: clangRootCursor.Kind(),
		CursorNameId: ctu.GoTu.CursorNameMap.Id(clangRootCursor.Spelling()),
		ParentIndex:  -1,
		TypeIndex:    ctu.GoTu.TypeMap.MustAutoKeyIndex(typekind.Unexposed),
		Tokens: ast.IndexPair{ // By definition, all the clang tokens.
			Head: 0,
			Len:  len(ctu.ClangTokens),
		},

		// Index can be set manually later but we do it here
		// to better show what's going on.
		// Index: 0, // Index no longer exists in this structure.
	})

	debug := false
	if debug {
		fmt.Printf("%[1]s %[1]d\n", clang.ChildVisit_Break)
		fmt.Printf("%[1]s %[1]d\n", clang.ChildVisit_Continue)
		fmt.Printf("%[1]s %[1]d\n", clang.ChildVisit_Recurse)
	}

	// Map cursor to its index in the list being created.
	cursorsSeen := make(map[clang.Cursor]int)
	nullCursor := clang.NewNullCursor()

	// interestingCursor returns true if new cursor is not null and not same as
	// old cursor.
	interestingCursor := func(newc, oldc clang.Cursor) bool {
		return !newc.Equal(nullCursor) && !newc.Equal(oldc)
	}

	// Loop through every cursor found, calling Visit on each one.

	// Grow the list of clang cursors by visiting the list of clang cursors.
	for parentIndex := 0; parentIndex < len(ctu.ClangCursors); parentIndex++ {

		if ctu.ClangCursors[parentIndex] == nullCursor {
			if debug {
				fmt.Printf("for loop: parentIndex %d, nullCursor\n", parentIndex)
			}
			continue
		}

		childcount := 0

		if debug {
			fmt.Printf("for loop: parentIndex %d, hash %x len %d\n",
				parentIndex,
				ctu.ClangCursors[parentIndex].HashCursor(),
				len(ctu.ClangCursors))
		}

		// The all important Visit().

		ctu.ClangCursors[parentIndex].Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {

			childcount++

			ownIndex := len(ctu.GoTu.Cursors)

			seenIndex, seen := cursorsSeen[cursor]

			if seen {
				// If this cursor has been seen already, it means the AST clang is walking us through
				// has an additional parent for the same child. Rather than traverse this child as a new
				// cursor that needs to be recorded, along with its children, register this duplication in
				// the tree by creating a cursorkind.Back entry.

				ctu.GoTu.Back[ownIndex] = seenIndex

				backCursor := ast.Cursor{
					CursorKindId: cursorkind.Back,
					ParentIndex:  parentIndex,
				}
				// Keep the lists the same length. Use the nullCursor as a place holder.
				ctu.GoTu.Cursors = append(ctu.GoTu.Cursors, backCursor)
				ctu.ClangCursors = append(ctu.ClangCursors, nullCursor)

			} else {

				if debug {
					fmt.Printf("parentIndex %d, visiting cursor hash %x, parent hash %x\n",
						parentIndex,
						cursor.HashCursor(),
						parent.HashCursor())
					fmt.Printf("childcount %d\n", childcount)
				}

				// Skip over builtin macro definitions.
				if cursor.Kind() == cursorkind.MacroDefinition {
					// Keying off of cursor.IsMacroBuiltin() doesn't work.
					// These macros don't return true for this call.
					// So skip based on the high level names that were found when no code
					// was submitted.
					if topLevelNamesToSkip != nil {
						name := cursor.Spelling()

						if topLevelNamesToSkip[name] {
							// Skipping the builtin macros to save on output. If it becomes necessary to include,
							// may want to have a test mode that can skip these for most of the test cases.
							return clang.ChildVisit_Continue
						}
					}
				}

				tokenRange := ctu.Tokenize(cursor)

				// Create the appropriate type entry.
				typeIndex := ctu.determineTypeIndex(cursor.Type())

				newCursor := ast.Cursor{
					CursorKindId: cursor.Kind(),
					CursorNameId: ctu.GoTu.CursorNameMap.Id(cursor.Spelling()),
					ParentIndex:  parentIndex,
					TypeIndex:    typeIndex,

					// Index can be set manually later but we do it here
					// to better show what's going on.
					// Index:  ownIndex,
					Tokens: tokenRange,
				}

				cursorsSeen[cursor] = len(ctu.GoTu.Cursors) // Length of either list would suffice.

				// N.B. Append to the two lists back to back, else risk
				// a return getting added in between and having their
				// lengths not match.
				ctu.GoTu.Cursors = append(ctu.GoTu.Cursors, newCursor)
				ctu.ClangCursors = append(ctu.ClangCursors, cursor)

				// Record when cursor has a referenced cursor and when it has a difference definition cursor.
				if referenced := cursor.Referenced(); interestingCursor(referenced, cursor) {

					referencedIndex, ok := cursorsSeen[referenced]
					if ok {
						ctu.GoTu.Referenced[ownIndex] = referencedIndex
					} else {
						fmt.Println(`panic("referenced cursor not yet seen")`)
						//panic("referenced cursor not yet seen")
						ctu.GoTu.Referenced[ownIndex] = -1
					}

					if definition := cursor.Referenced(); interestingCursor(definition, referenced) {
						definitionIndex, ok := cursorsSeen[definition]
						if ok {
							ctu.GoTu.Definition[ownIndex] = definitionIndex
						} else {
							fmt.Println(`panic("definition cursor not yet seen")`)
							//panic("definition cursor not yet seen")
							ctu.GoTu.Definition[ownIndex] = -1
						}
					}
				}
			}

			// Determining the children doesn't have to be done here.
			// It is enough that the ParentIndex was set in this visit.
			// But it is done here to better show what's going on.

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
			// End of setup that could be done later.

			return clang.ChildVisit_Continue // Continue to next sibling
		})
	}

	// These could also be called from DecodeFinish if we knew we just wanted them done
	// after gob decoding and not when struct is first populated.
	ctu.GoTu.SetBinaryOperatorNames()
	// Maybe don't call this here either.
	//ctu.GoTu.SetBackChildren()

	return nil
}

func (ctu *ClangTranslationUnit) determineTypeIndex(ctype clang.Type) int {

	typeIndex, found := ctu.typeIndexes[ctype]
	if found {
		return typeIndex
	}

	tkind := ctype.Kind()
	typeIndex = ctu.GoTu.TypeMap.AutoKeyIndex(tkind)
	if typeIndex == -1 {

		// The type spelling.

		typespelling := ctype.Spelling()

		// The type alignment.

		alignof := uint64(0)
		var err error
		switch tkind {
		case typekind.Void:
			break
		default:
			alignof, err = ctype.AlignOf()
			if err != nil {
				panic(fmt.Sprintf("In calling AlignOf, type %v %s %s: %s", ctype, tkind, typespelling, err))
			}
		}

		// The type size.

		sizeof := uint64(0)
		switch tkind {
		case typekind.Void,
			typekind.VariableArray:
			break
		default:
			sizeof, err = ctype.SizeOf()
			if err != nil {
				panic(fmt.Sprintf("In calling SizeOf, type %v %s %s: %s", ctype, tkind, typespelling, err))
			}
		}

		// Add the appropriate record to the TypeMap.

		switch tkind {
		case typekind.Record:

			typeIndex, err = ctu.GoTu.TypeMap.AddRecord(ast.TypeRecord{
				Align:        int(alignof), // TBD change api to return int rather than uint64.
				Size:         int(sizeof),  // TBD change api to return int rather than uint64.
				TypeSpelling: typespelling,
			})
			if err != nil {
				errmsg := fmt.Sprintf("%[1]s:%[1]d", ctype.Kind())
				panic(errmsg + ": " + err.Error())
			}

		case typekind.Enum:
			// Same as Record

			typeIndex, err = ctu.GoTu.TypeMap.AddEnum(ast.TypeEnum{
				Align:        int(alignof), // TBD change api to return int rather than uint64.
				Size:         int(sizeof),  // TBD change api to return int rather than uint64.
				TypeSpelling: typespelling,
			})
			if err != nil {
				errmsg := fmt.Sprintf("%[1]s:%[1]d", ctype.Kind())
				panic(errmsg + ": " + err.Error())
			}

		case typekind.Typedef:
			typeIndex = ctu.determineTypeIndex2(tkind,
				"Typedefee",
				ctype.CanonicalType(), // Convert Canonical type to TypeTypedef struct.
				func(pointeetypeindex int) (int, error) {
					return ctu.GoTu.TypeMap.AddTypedef(ast.TypeTypedef{
						TypeSpelling:        typespelling,
						UnderlyingTypeIndex: pointeetypeindex,
					})
				})
		case typekind.Pointer:
			typeIndex = ctu.determineTypeIndex2(tkind,
				"Pointee",
				ctype.PointeeType(), // Convert Pointee type to TypePointer struct.
				func(pointeetypeindex int) (int, error) {
					return ctu.GoTu.TypeMap.AddPointer(ast.TypePointer{
						UnderlyingTypeIndex: pointeetypeindex,
					})
				})
		case typekind.Elaborated:
			typeIndex = ctu.determineTypeIndex2(tkind,
				"Elaboratee",
				ctype.NamedType(), // Convert Named type to TypeElaborated struct.
				func(pointeetypeindex int) (int, error) {
					return ctu.GoTu.TypeMap.AddElaborated(ast.TypeElaborated{
						UnderlyingTypeIndex: pointeetypeindex,
					})
				})

		case typekind.FunctionNoProto,
			typekind.FunctionProto:

			// Build list of argument types.
			var argtypeindexes []int
			numargs := ctype.NumArgTypes()
			for i := int32(0); i < numargs; i++ {
				at := ctype.ArgType(uint32(i))
				ati := ctu.mustDetermineSubTypeIndex(tkind, at, "arg")
				argtypeindexes = append(argtypeindexes, ati)
			}

			typeIndex = ctu.addSuperWithOneSubType(tkind,
				ctype.ResultType(),

				func(resulttypeindex int) (int, error) {
					return ctu.GoTu.TypeMap.AddFunction(ast.TypeFunction{
						TypeKindKind: ast.TypeKindKind{tkind},
						ResultTypeId: resulttypeindex,
						ArgIds:       argtypeindexes,
						TypeSpelling: typespelling,
					})
				})

		case typekind.ConstantArray:
			numelem := ctype.NumElements()
			if numelem < 0 {
				panic(fmt.Sprintf("element count is negative?"))
			}

			typeIndex = ctu.addSuperWithOneSubType(tkind,
				ctype.ElementType(),

				func(elemtypeindex int) (int, error) {
					return ctu.GoTu.TypeMap.AddConstantArray(ast.TypeConstantArray{
						ElemCount: int(numelem),
						TypeVariableArray: ast.TypeVariableArray{
							ElemTypeId:   elemtypeindex,
							Align:        int(alignof), // TBD change api to return int rather than uint64.
							Size:         int(sizeof),  // TBD change api to return int rather than uint64.
							TypeSpelling: typespelling,
						},
					})
				})

		case typekind.VariableArray:

			typeIndex = ctu.addSuperWithOneSubType(tkind,
				ctype.ElementType(),

				func(elemtypeindex int) (int, error) {
					return ctu.GoTu.TypeMap.AddVariableArray(ast.TypeVariableArray{
						ElemTypeId:   elemtypeindex,
						Align:        int(alignof), // TBD change api to return int rather than uint64.
						Size:         int(sizeof),  // TBD change api to return int rather than uint64.
						TypeSpelling: typespelling,
					})
				})

		default:
			if tkind.IsBuiltin() {
				typeIndex, err = ctu.GoTu.TypeMap.AddIntrinsic(ast.TypeIntrinsic{
					TypeKindKind: ast.TypeKindKind{tkind},
					Align:        int(alignof), // TBD change api to return int rather than uint64.
					Size:         int(sizeof),  // TBD change api to return int rather than uint64.
					TypeSpelling: typespelling,
				})
				if err != nil {
					panic(err)
				}
			} else {
				errmsg := fmt.Sprintf("%[1]s:%[1]d", tkind)
				panic("ctu.determineTypeIndex type not yet handled: " + errmsg)
				/*
				 */
				// TBD Stop gap measure while other type kinds are implemented.
				typeIndex = ctu.GoTu.TypeMap.MustAutoKeyIndex(typekind.Unexposed) // TBD
			}
		}
	}

	// Cache typeIndex.
	ctu.typeIndexes[ctype] = typeIndex
	return typeIndex
}

func (ctu *ClangTranslationUnit) determineTypeIndex2(tkind typekind.Kind,
	errname string,
	pointeetype clang.Type,
	addFn func(pointeetypeindex int) (int, error)) int {

	if pointeetype.Kind() == typekind.Invalid {
		panic(errname + " type is invalid?")
	}

	pointeetypeindex := ctu.mustDetermineSubTypeIndex(tkind, pointeetype, errname)

	typeIndex, err := addFn(pointeetypeindex)
	if err != nil {
		errmsg := fmt.Sprintf("%[1]s:%[1]d", tkind)
		errmsg += fmt.Sprintf(" %s %[2]s:%[2]d", errname, pointeetype.Kind())
		errmsg += fmt.Sprintf(" Key%v", ctu.GoTu.TypeMap.Keys[pointeetypeindex])
		errmsg += fmt.Sprintf(" %s typeindex %d", errname, pointeetypeindex)
		panic(errmsg + ": " + err.Error())
	}
	return typeIndex
}

func (ctu *ClangTranslationUnit) addSuperWithOneSubType(superTypeKind typekind.Kind,
	subType clang.Type,
	addFn func(pointeetypeindex int) (int, error)) int {

	if subType.Kind() == typekind.Invalid {
		panic(fmt.Sprintf("element type is invalid?"))
	}

	subtypeindex := ctu.mustDetermineSubTypeIndex(superTypeKind, subType, "subtype")

	typeIndex, err := addFn(subtypeindex)
	if err != nil {
		errmsg := fmt.Sprintf("%[1]s:%[1]d", superTypeKind)
		errmsg += fmt.Sprintf(" subType %[1]s:%[1]d", subType.Kind())
		errmsg += fmt.Sprintf(" Key[%d]%v", subtypeindex, ctu.GoTu.TypeMap.Keys[subtypeindex])
		panic(errmsg + ": " + err.Error())
	}
	return typeIndex
}

func (ctu *ClangTranslationUnit) mustDetermineSubTypeIndex(superTypeKind typekind.Kind, subType clang.Type, errname string) int {
	subtypeindex := ctu.determineTypeIndex(subType)
	if subtypeindex <= 1 {
		// Would indicate an array of an Invalid or Unexposed type kind.
		errmsg := fmt.Sprintf("%[1]s:%[1]d", superTypeKind)
		errmsg += fmt.Sprintf(" %s %[2]s:%[2]d", errname, subType.Kind())
		errmsg += fmt.Sprintf(" %s typeindex %d", errname, subtypeindex)
		errmsg += fmt.Sprintf(" Key%v", ctu.GoTu.TypeMap.Keys[subtypeindex])
		panic(errmsg + ": " + errname + " typeindex <= 1")
	}
	return subtypeindex
}

func diagnosticError(tu *clang.TranslationUnit) error {
	diags := tu.Diagnostics()
	if len(diags) == 0 {
		return nil
	}

	b := new(strings.Builder)
	for _, diag := range diags {
		fmt.Fprintln(b, "diagnostic", diag.FormatDiagnostic(0))
	}

	return fmt.Errorf("%s", b.String())
}
