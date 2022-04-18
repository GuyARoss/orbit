// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package allocatedstack

import (
	"container/list"
)

type Stack struct {
	l    *list.List
	size uint
}

func (s *Stack) Add(v string) {
	if s.l.Len() > int(s.size)-1 {
		e := s.l.Back()
		s.l.Remove(e)
	}

	s.l.PushFront(v)
}

func (s *Stack) Arr() []string {
	final := make([]string, 0)

	e := s.l.Front()
	for i := 0; i < s.l.Len(); i++ {
		final = append(final, e.Value.(string))

		e = e.Next()
	}

	return final
}

func (s *Stack) Contains(v string) bool {
	e := s.l.Front()
	for i := 0; i < s.l.Len(); i++ {
		if e.Value == v {
			return true
		}

		e = e.Next()
	}

	return false
}

func (s *Stack) Len() int {
	return s.l.Len()
}

func New(size uint) *Stack {
	return &Stack{
		size: size,
		l:    list.New(),
	}
}
