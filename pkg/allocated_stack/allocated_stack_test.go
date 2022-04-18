// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package allocatedstack

import (
	"testing"

	"github.com/google/uuid"
)

func TestStack(t *testing.T) {
	t.Run("not exceed maxSize", func(t *testing.T) {
		s := New(5)

		for i := 0; i < 15; i++ {
			s.Add("stuff")

			if i >= 5 && s.Len() > 5 {
				t.Errorf("incorrect size expected %d got %d", 5, s.Len())
			}

			if i < 5 && s.Len() != i+1 {
				t.Errorf("incorrect size expected %d got %d", i+1, s.Len())
			}
		}
	})

	t.Run("arr", func(t *testing.T) {
		s := New(2)
		s.Add("1")
		s.Add("2")
		s.Add("3")
		s.Add("4")
		s.Add("5")
		s.Add("6")

		a := s.Arr()
		if a[0] != "6" && a[1] != "5" {
			t.Errorf("got %s", a)
		}
	})

	t.Run("remove correctly", func(t *testing.T) {
		length := 3
		s := New(uint(length))

		lst := make([]string, 20)
		for i := 0; i < 10; i++ {
			str := uuid.NewString()
			lst[i] = str
			s.Add(str)

			if !s.Contains(str) {
				t.Errorf("does not contain str")
			}

			if i > length {
				offset := i - length
				if offset > 1 {
					shouldNotExist := lst[:offset-1]

					for _, e := range shouldNotExist {
						if s.Contains(e) {
							t.Errorf("element should not exist on stack")
						}
					}
				} else {
					if s.Contains(lst[0]) {
						t.Errorf("element should not exist on stack")
					}
				}
			}
		}

	})
}
