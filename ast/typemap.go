package ast

import (
	"fmt"
	"strings"

	"github.com/frankreh/go-clang-v5.0/clang/typekind"
)

// TypeKey keys a single index int to a kind of type and the
// index into the appropriate slice.
type TypeKey struct {
	TypeKind typekind.Kind
	TypeId   int // Index into relevant TypeMap slices.
}

type Type interface {
	Kind() typekind.Kind
}

type TypeKindKind struct {
	TypeKind typekind.Kind
}

func (tk *TypeKindKind) Kind() typekind.Kind {
	return tk.TypeKind
}

type TypeIntrinsic struct {
	TypeKindKind
	TypeSpelling string
	Align        int // TBD unsigned byte or unsigned short?
	Size         int
}

type TypePointer struct {
	UnderlyingTypeIndex int
}

func (p *TypePointer) Kind() typekind.Kind {
	return typekind.Pointer
}

type TypeElaborated struct {
	UnderlyingTypeIndex int
}

func (p *TypeElaborated) Kind() typekind.Kind {
	return typekind.Elaborated
}

type TypeRecord struct {
	Align        int // TBD unsigned byte or unsigned short?
	Size         int
	TypeSpelling string
}

func (p *TypeRecord) Kind() typekind.Kind {
	return typekind.Record
}

type TypeEnum struct {
	Align        int // TBD unsigned byte or unsigned short?
	Size         int
	TypeSpelling string
}

func (p *TypeEnum) Kind() typekind.Kind {
	return typekind.Enum
}

type TypeTypedef struct {
	UnderlyingTypeIndex int
	TypeSpelling        string
}

func (p *TypeTypedef) Kind() typekind.Kind {
	return typekind.Typedef
}

type TypeFunction struct {
	TypeKindKind
	TypeSpelling string
}

type TypeConstArray struct {
	ElemTypeId   int
	ElemCount    int
	Align        int
	Size         int
	TypeSpelling string
}

func (p *TypeConstArray) Kind() typekind.Kind {
	return typekind.ConstantArray
}

// AddIntrinsic adds TypeIntrinsic to its mapping and returns the new key index for it.
// But returns an error if the underlying type index is already in the mapping.
func (tm *TypeMap) AddIntrinsic(r TypeIntrinsic) (int, error) {
	if !r.TypeKind.IsBuiltin() {
		return -1, fmt.Errorf("TypeMap.AddIntrinsic: typekind %[1]s:%[1]d is not builtin", r.TypeKind)
	}
	if r.TypeSpelling == "" {
		return -1, fmt.Errorf("TypeMap.AddIntrinsic.TypeSpelling is null")
	}
	if r.Align == 0 {
		return -1, fmt.Errorf("TypeMap.AddIntrinsic.Align is 0")
	}
	if r.Size == 0 {
		// Allow a few type kinds to be coming in with a size of 0.
		switch r.TypeKind {
		case typekind.VariableArray:
			break
		default:
			return -1, fmt.Errorf("TypeMap.AddIntrinsic.Size is 0")
		}
	}

	l := &tm.Intrinsics
	for i, e := range *l {
		if r == e {
			panic(fmt.Sprintf("Entry already exists at %d", i))
		}
	}
	*l = append(*l, r)

	return tm.addToKeys(r.TypeKind, len(*l)-1), nil
}

// AddPointer adds TypePointer to its mapping and returns the new key index for it.
// But returns an error if the underlying type index is already in the mapping.
func (tm *TypeMap) AddPointer(r TypePointer) (int, error) {
	/*
		if u, err := tm.Type(r.UnderlyingTypeIndex); u == nil || err != nil {
			if u == nil {
				return -1, fmt.Errorf("TypeMap.AddPointer: underlying type index:%d key:%v resulted in nil Type interface",
					r.UnderlyingTypeIndex,
					tm.Keys[r.UnderlyingTypeIndex],
				)
			}
			return -1, fmt.Errorf("TypeMap.AddPointer: underlying type index error: %s", err)
		}

		l := &tm.Pointers
		for i, e := range *l {
			if r == e {
				panic(fmt.Sprintf("Entry already exists at %d", i))
			}
		}
		*l = append(*l, r)
		return tm.addToKeys(typekind.Pointer, len(*l)-1), nil
	*/
	// Pointer underlying indexes are not stored in their own list. Index is stored directly in key.
	if err := indexCheck(r.UnderlyingTypeIndex, len(tm.Keys), "Keys"); err != nil {
		return -1, fmt.Errorf("TypeMap.AddPointer: %s", err)
	}
	return tm.addToKeys(typekind.Pointer, r.UnderlyingTypeIndex), nil
}

// AddElaborated adds TypeElaborated to its mapping and returns the new key index for it.
// But returns an error if the underlying type index is already in the mapping.
func (tm *TypeMap) AddElaborated(r TypeElaborated) (int, error) {
	/*
		if u, err := tm.Type(r.UnderlyingTypeIndex); u == nil || err != nil {
			if u == nil {
				return -1, fmt.Errorf("TypeMap.AddElaborated: underlying type index:%d key:%v resulted in nil Type interface",
					r.UnderlyingTypeIndex,
					tm.Keys[r.UnderlyingTypeIndex],
				)
			}
			return -1, fmt.Errorf("TypeMap.AddElaborated: underlying type index error: %s", err)
		}

		l := &tm.Elaborateds
		for i, e := range *l {
			if r == e {
				panic(fmt.Sprintf("Entry already exists at %d", i))
			}
		}
		*l = append(*l, r)
		return tm.addToKeys(typekind.Elaborated, len(*l)-1), nil
	*/
	// Elaborated underlying indexes are not stored in their own list. Index is stored directly in key.
	if err := indexCheck(r.UnderlyingTypeIndex, len(tm.Keys), "Keys"); err != nil {
		return -1, fmt.Errorf("TypeMap.AddElaborated: %s", err)
	}
	return tm.addToKeys(typekind.Elaborated, r.UnderlyingTypeIndex), nil
}

// AddRecord adds TypeRecord to its mapping and returns the new key index for it.
func (tm *TypeMap) AddRecord(r TypeRecord) (int, error) {
	l := &tm.Records
	for i, e := range *l {
		if r == e {
			panic(fmt.Sprintf("Entry already exists at %d", i))
		}
	}
	*l = append(*l, r)
	return tm.addToKeys(typekind.Record, len(*l)-1), nil
}

// AddEnum adds TypeEnum to its mapping and returns the new key index for it.
func (tm *TypeMap) AddEnum(r TypeEnum) (int, error) {
	l := &tm.Enums
	for i, e := range *l {
		if r == e {
			panic(fmt.Sprintf("Entry already exists at %d", i))
		}
	}
	*l = append(*l, r)
	return tm.addToKeys(typekind.Enum, len(*l)-1), nil
}

// AddTypedef adds TypeTypedef to its mapping and returns the new key index for it.
func (tm *TypeMap) AddTypedef(r TypeTypedef) (int, error) {
	l := &tm.Typedefs
	for i, e := range *l {
		if r == e {
			panic(fmt.Sprintf("Entry already exists at %d", i))
		}
	}
	*l = append(*l, r)
	return tm.addToKeys(typekind.Typedef, len(*l)-1), nil
}

// AddConstArray adds TypeFunction to its mapping and returns the new key index for it.
func (tm *TypeMap) AddFunction(r TypeFunction) (int, error) {
	l := &tm.Functions
	for i, e := range *l {
		if r == e {
			panic(fmt.Sprintf("Entry already exists at %d", i))
		}
	}
	*l = append(*l, r)
	return tm.addToKeys(r.TypeKind, len(*l)-1), nil
}

// AddConstArray adds TypeConstArray to its mapping and returns the new key index for it.
func (tm *TypeMap) AddConstArray(r TypeConstArray) (int, error) {
	l := &tm.ConstArrays
	for i, e := range *l {
		if r == e {
			panic(fmt.Sprintf("Entry already exists at %d", i))
		}
	}
	*l = append(*l, r)
	return tm.addToKeys(typekind.ConstantArray, len(*l)-1), nil
}

// TypeMap maps a Cursor.TypeIndex (an int) to a Type.
type TypeMap struct {
	Keys       []TypeKey
	Intrinsics []TypeIntrinsic
	//Pointers   []TypePointer
	//Elaborateds   []TypeElaborated
	Records     []TypeRecord
	Enums       []TypeEnum
	Typedefs    []TypeTypedef
	Functions   []TypeFunction
	ConstArrays []TypeConstArray
}

// Init ensures the struct is setup properly before first use.
func (tm *TypeMap) Init() {
	if len(tm.Keys) > 0 {
		return
	}
	// Keys gets entry
	// [0] representing typekind.Invalid and
	// [1] representing typekind.Unexposed.
	// Use index of -1 to catch incorrect use of these keys with any mapping slice.
	tm.Keys = append(tm.Keys, TypeKey{typekind.Invalid, -1})
	tm.Keys = append(tm.Keys, TypeKey{typekind.Unexposed, -1})
}

// addToKeys adds a TypeKey record to the slice for the given type kind and corresponding index.
// Returns the index of the record that was added.
func (tm *TypeMap) addToKeys(kind typekind.Kind, index int) int {
	tm.Keys = append(tm.Keys, TypeKey{kind, index})
	return len(tm.Keys) - 1
}

// AutoKeyIndex returns a key index that is 0 or greater if the typekind maps
// directly to a key index. But returns -1 if a specific Add method needs to be
// used for the typekind in question.
func (tm *TypeMap) AutoKeyIndex(t typekind.Kind) int {
	switch t {
	case typekind.Invalid:
		return 0 // Keys[0]
	case typekind.Unexposed:
		return 1 // Keys[1]
	}
	return -1
}

func (tm *TypeMap) MustAutoKeyIndex(t typekind.Kind) int {
	r := tm.AutoKeyIndex(t)
	if r >= 0 {
		return r
	}
	panic(fmt.Sprintf("typekind expected to have automatic key index: %[1]s:%[1]d", t))
}

func (tm *TypeMap) Type(i int) (Type, error) {
	if i < 0 || i >= len(tm.Keys) {
		return nil, fmt.Errorf("TypeMap.Type(%d) out of range", i, len(tm.Keys)) // TBD change with ast.Error
	}
	li := tm.Keys[i].TypeId

	if tm.Keys[i].TypeKind.IsBuiltin() {
		l := tm.Intrinsics
		if err := indexCheck(li, len(l), "Intrinics"); err != nil {
			return nil, err
		}
		return &l[li], nil
	}

	switch tm.Keys[i].TypeKind {
	case typekind.Invalid:
		return nil, nil
	case typekind.Unexposed:
		// libclang may not expose the type. But this package or the astbridge package may not expose it either.
		return nil, nil
	case typekind.Pointer:
		/*
			l := tm.Pointers
			if err := indexCheck(li, len(l), "Pointers"); err != nil {
				return nil, err
			}
			return &l[li], nil
		*/
		l := tm.Keys
		if err := indexCheck(li, len(l), "Keys"); err != nil {
			return nil, err
		}
		// Create the instance on the fly.
		return &TypePointer{li}, nil
	case typekind.Elaborated:
		/*
			l := tm.Elaborateds
			if err := indexCheck(li, len(l), "Elaborateds"); err != nil {
				return nil, err
			}
			return &l[li], nil
		*/
		l := tm.Keys
		if err := indexCheck(li, len(l), "Keys"); err != nil {
			return nil, err
		}
		// Create the instance on the fly.
		return &TypeElaborated{li}, nil
	case typekind.Record:
		l := tm.Records
		if err := indexCheck(li, len(l), "Records"); err != nil {
			return nil, err
		}
		return &l[li], nil
	case typekind.Enum:
		l := tm.Enums
		if err := indexCheck(li, len(l), "Enums"); err != nil {
			return nil, err
		}
		return &l[li], nil
	case typekind.Typedef:
		l := tm.Typedefs
		if err := indexCheck(li, len(l), "Typedefs"); err != nil {
			return nil, err
		}
		return &l[li], nil
	case typekind.FunctionProto,
		typekind.FunctionNoProto:
		l := tm.Functions
		if err := indexCheck(li, len(l), "Functions"); err != nil {
			return nil, err
		}
		return &l[li], nil
	case typekind.ConstantArray:
		l := tm.ConstArrays
		if err := indexCheck(li, len(l), "ConstArrays"); err != nil {
			return nil, err
		}
		return &l[li], nil
	}
	return nil, nil
}

func (tm TypeMap) GoString() string {
	b := new(strings.Builder)
	fmt.Fprintf(b, "TypeMap{\n")
	if len(tm.Keys) > 0 {
		fmt.Fprintf(b, "    Keys: %v\n", tm.Keys)
	}
	if len(tm.Intrinsics) > 0 {
		fmt.Fprintf(b, "    Intrinsics: %v\n", tm.Intrinsics)
	}
	/*
		if len(tm.Pointers) > 0 {
			fmt.Fprintf(b, "    Pointers: %v\n", tm.Pointers)
		}
		if len(tm.Elaborateds) > 0 {
			fmt.Fprintf(b, "    Elaborateds: %v\n", tm.Elaborateds)
		}
	*/
	if len(tm.Records) > 0 {
		fmt.Fprintf(b, "    Records: %v\n", tm.Records)
	}
	if len(tm.Enums) > 0 {
		fmt.Fprintf(b, "    Enums: %v\n", tm.Enums)
	}
	if len(tm.Typedefs) > 0 {
		fmt.Fprintf(b, "    Typedefs: %v\n", tm.Typedefs)
	}
	if len(tm.Functions) > 0 {
		fmt.Fprintf(b, "    Functions: %v\n", tm.Functions)
	}
	if len(tm.ConstArrays) > 0 {
		fmt.Fprintf(b, "    ConstArrays: %v\n", tm.ConstArrays)
	}
	fmt.Fprintf(b, "}")

	return b.String()
}

// Len returns number of Types mapped. Valid indexes will be 0..Len-1.
func (tm *TypeMap) Len() int {
	return len(tm.Keys)
}

func (a *TypeMap) AssertEqual(b *TypeMap) error {
	if a == b {
		return nil
	}
	if err := a.assertEqualKeys(b); err != nil {
		return err
	}
	if err := a.assertEqualIntrinsics(b); err != nil {
		return err
	}
	/*
		if err := a.assertEqualPointers(b); err != nil {
			return err
		}
		if err := a.assertEqualElaborateds(b); err != nil {
			return err
		}
	*/
	if err := a.assertEqualRecords(b); err != nil {
		return err
	}
	if err := a.assertEqualEnums(b); err != nil {
		return err
	}
	if err := a.assertEqualTypedefs(b); err != nil {
		return err
	}
	if err := a.assertEqualFunctions(b); err != nil {
		return err
	}
	if err := a.assertEqualConstArrays(b); err != nil {
		return err
	}
	return nil
}

func (a *TypeMap) assertEqualKeys(b *TypeMap) error {
	if len(a.Keys) != len(b.Keys) {
		return fmt.Errorf("TypeMap unequal keys lengths, %d %d",
			len(a.Keys), len(b.Keys))
	}
	for i, v := range a.Keys {
		v2 := b.Keys[i]
		if v != v2 {
			return fmt.Errorf("TypeMap unequal keys entry, %d %s %s",
				i, v, v2)
		}
	}
	return nil
}

func (a *TypeMap) assertEqualIntrinsics(b *TypeMap) error {
	if len(a.Intrinsics) != len(b.Intrinsics) {
		return fmt.Errorf("TypeMap unequal Intrinsics lengths, %d %d",
			len(a.Intrinsics), len(b.Intrinsics))
	}
	for i, v := range a.Intrinsics {
		v2 := b.Intrinsics[i]
		if v != v2 {
			return fmt.Errorf("TypeMap unequal Intrinsics entry, %d %s %s",
				i, v, v2)
		}
	}
	return nil
}

/*
func (a *TypeMap) assertEqualPointers(b *TypeMap) error {
	if len(a.Pointers) != len(b.Pointers) {
		return fmt.Errorf("TypeMap unequal Pointers lengths, %d %d",
			len(a.Pointers), len(b.Pointers))
	}
	for i, v := range a.Pointers {
		v2 := b.Pointers[i]
		if v != v2 {
			return fmt.Errorf("TypeMap unequal Pointers entry, %d %s %s",
				i, v, v2)
		}
	}
	return nil
}
func (a *TypeMap) assertEqualElaborateds(b *TypeMap) error {
	if len(a.Elaborateds) != len(b.Elaborateds) {
		return fmt.Errorf("TypeMap unequal Elaborateds lengths, %d %d",
			len(a.Elaborateds), len(b.Elaborateds))
	}
	for i, v := range a.Elaborateds {
		v2 := b.Elaborateds[i]
		if v != v2 {
			return fmt.Errorf("TypeMap unequal Elaborateds entry, %d %s %s",
				i, v, v2)
		}
	}
	return nil
}
*/

func (a *TypeMap) assertEqualRecords(b *TypeMap) error {
	if len(a.Records) != len(b.Records) {
		return fmt.Errorf("TypeMap unequal Records lengths, %d %d",
			len(a.Records), len(b.Records))
	}
	for i, v := range a.Records {
		v2 := b.Records[i]
		if v != v2 {
			return fmt.Errorf("TypeMap unequal Records entry, %d %s %s",
				i, v, v2)
		}
	}
	return nil
}

func (a *TypeMap) assertEqualEnums(b *TypeMap) error {
	if len(a.Enums) != len(b.Enums) {
		return fmt.Errorf("TypeMap unequal Enums lengths, %d %d",
			len(a.Enums), len(b.Enums))
	}
	for i, v := range a.Enums {
		v2 := b.Enums[i]
		if v != v2 {
			return fmt.Errorf("TypeMap unequal Enums entry, %d %s %s",
				i, v, v2)
		}
	}
	return nil
}

func (a *TypeMap) assertEqualTypedefs(b *TypeMap) error {
	if len(a.Typedefs) != len(b.Typedefs) {
		return fmt.Errorf("TypeMap unequal Typedefs lengths, %d %d",
			len(a.Typedefs), len(b.Typedefs))
	}
	for i, v := range a.Typedefs {
		v2 := b.Typedefs[i]
		if v != v2 {
			return fmt.Errorf("TypeMap unequal Typedefs entry, %d %s %s",
				i, v, v2)
		}
	}
	return nil
}

func (a *TypeMap) assertEqualFunctions(b *TypeMap) error {
	if len(a.Functions) != len(b.Functions) {
		return fmt.Errorf("TypeMap unequal Functions lengths, %d %d",
			len(a.Functions), len(b.Functions))
	}
	for i, v := range a.Functions {
		v2 := b.Functions[i]
		if v != v2 {
			return fmt.Errorf("TypeMap unequal Functions entry, %d %s %s",
				i, v, v2)
		}
	}
	return nil
}

func (a *TypeMap) assertEqualConstArrays(b *TypeMap) error {
	if len(a.ConstArrays) != len(b.ConstArrays) {
		return fmt.Errorf("TypeMap unequal ConstArrays lengths, %d %d",
			len(a.ConstArrays), len(b.ConstArrays))
	}
	for i, v := range a.ConstArrays {
		v2 := b.ConstArrays[i]
		if v != v2 {
			return fmt.Errorf("TypeMap unequal ConstArrays entry, %d %s %s",
				i, v, v2)
		}
	}
	return nil
}

// TBD kind of got to here.
/*
// Does this struct even need a map? The astbridge code will keep its own
// map of seenTypes so will already know if there is a TypeId
func (tm *TypeMap) DecodeFinish() {
	if tm.m == nil {
		tm.m = make(map[Type]TokenId, len(tm.Tokens))
	}
	for i, t := range tm.Tokens {
		tm.m[t] = TokenId(i) // cast index to TokenId
	}
}

func (a *TokenMap) assertEqualMap(b *TokenMap) error {
	if a.m == nil && b.m == nil {
		return nil
	}
	if len(a.m) != len(b.m) {
		return fmt.Errorf("TokenMap unequal map lengths, %d %d",
			len(a.m), len(b.m))
	}
	for k, v := range a.m {
		v2 := b.m[k]
		if v != v2 {
			return fmt.Errorf("TokenMap unequal map entry, %d %s %s",
				k, v, v2)
		}
	}
	return nil
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
*/
func indexCheck(i, length int, name string) error {
	if i >= 0 && i < length {
		return nil
	}
	return fmt.Errorf("out of range: %d %s:len:%d", i, name, length)
}
