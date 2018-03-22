package ast

import (
	"fmt"
)

// StringMap creates indexes for strings. Returning a new index when a new string is
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
type StringMap struct {
	m       map[string]int
	Strings []string
}

/*
func (s StringMap) String() string {
	return fmt.Sprintf("%v", s.Strings)
}
*/
func (s StringMap) GoString() string {
	return fmt.Sprintf("{Strings: %#v}", s.Strings)
}

// Len returns number of strings mapped. Valid indexes will be 0..Len-1.
func (s *StringMap) Len() int {
	return len(s.Strings)
}

func (s *StringMap) DecodeFinish() {
	if s.m == nil {
		s.m = make(map[string]int, len(s.Strings))
	}
	for i, str := range s.Strings {
		s.m[str] = i
	}
}

func (a *StringMap) AssertEqual(b *StringMap) error {
	if a == b {
		return nil
	}
	if err := a.assertEqualSlice(b); err != nil {
		return err
	}
	if err := a.assertEqualMap(b); err != nil {
		return err
	}
	return nil
}
func (a *StringMap) assertEqualSlice(b *StringMap) error {
	if len(a.Strings) != len(b.Strings) {
		return fmt.Errorf("StringMap unequal slice lengths, %d %d",
			len(a.Strings), len(b.Strings))
	}
	for i, v := range a.Strings {
		v2 := b.Strings[i]
		if v != v2 {
			return fmt.Errorf("StringMap unequal slice entry, %d %s %s",
				i, v, v2)
		}
	}
	return nil
}
func (a *StringMap) assertEqualMap(b *StringMap) error {
	if a.m == nil && b.m == nil {
		return nil
	}
	if len(a.m) != len(b.m) {
		return fmt.Errorf("StringMap unequal map lengths, %d %d",
			len(a.m), len(b.m))
	}
	for k, v := range a.m {
		v2 := b.m[k]
		if v != v2 {
			return fmt.Errorf("StringMap unequal map entry, %d %s %s",
				k, v, v2)
		}
	}
	return nil
}

func (s *StringMap) Init(strings []string) {
	s.m = make(map[string]int)
	s.Strings = make([]string, len(strings))
	for i, str := range strings {
		s.m[str] = i
		s.Strings[i] = str
	}
}

// Id returns the id for this string.
// First id will be zero.
func (s *StringMap) Id(str string) int {
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
func (s *StringMap) ToString(id int) string {
	return s.Strings[id]
}
