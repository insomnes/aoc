package solution

import (
	"fmt"
	"strings"
)

type Set[T comparable] map[T]struct{}

func (s Set[T]) String() string {
	if s.Empty() {
		return "Set{}"
	}
	var sb strings.Builder
	sb.WriteString("Set{")
	i := 0
	for val := range s {
		if i > 0 {
			sb.WriteString(", ")
		}
		i++
		sb.WriteString(fmt.Sprintf("%v", val))
	}
	sb.WriteString(fmt.Sprintf("} (%d)", s.Size()))
	return sb.String()
}

func NewSet[T comparable](n int) Set[T] {
	return make(map[T]struct{}, n)
}

func NewSetFromValues[T comparable](values ...T) Set[T] {
	s := NewSet[T](len(values))
	for _, value := range values {
		s[value] = struct{}{}
	}
	return s
}

func (s Set[T]) Add(val T) {
	s[val] = struct{}{}
}

func (s Set[T]) Remove(val T) {
	delete(s, val)
}

func (s Set[T]) Contains(val T) bool {
	_, ok := s[val]
	return ok
}

func (s Set[T]) Empty() bool {
	return s.Size() == 0
}

func (s Set[T]) Size() int {
	return len(s)
}

func (s Set[T]) Intersection(other Set[T]) Set[T] {
	minimal, toCheck := s, other
	if other.Size() < s.Size() {
		minimal, toCheck = other, s
	}
	intersect := NewSet[T](minimal.Size())
	for val := range minimal {
		if toCheck.Contains(val) {
			intersect.Add(val)
		}
	}
	return intersect
}

func (s Set[T]) Union(other Set[T]) Set[T] {
	union := NewSet[T](s.Size() + other.Size())
	for val := range s {
		union.Add(val)
	}
	for val := range other {
		union.Add(val)
	}
	return union
}

func (s Set[T]) Copy() Set[T] {
	c := NewSet[T](s.Size())
	for val := range s {
		c.Add(val)
	}
	return c
}
