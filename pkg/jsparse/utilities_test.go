// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package jsparse

import (
	"errors"
	"strings"
	"testing"
)

func TestCleanExportDefaultName_Errors(t *testing.T) {
	tt := []struct {
		i string
		o error
	}{
		{"export default", ErrFunctionExport},
	}

	for i, d := range tt {
		_, err := extractDefaultExportName(d.i)
		if !errors.Is(d.o, err) {
			t.Errorf("(%d) expected error", i)
		}
	}
}

func TestExportDefaultName(t *testing.T) {
	tt := []struct {
		i string
		o string
	}{
		{"export default Test", "Test"},
		{"export default SomethingCool  ", "SomethingCool"},
		{"export default () => {}", ""},
	}

	for i, c := range tt {
		name, err := extractDefaultExportName(c.i)

		if err != nil {
			t.Errorf("(%d) error exception should not be thrown %d", i, err)
		}

		if name != c.o {
			t.Errorf("expected %s got %s \n", "Test", name)
		}
	}
}

func TestExtension(t *testing.T) {
	pn := NewDocument("", "./thing.png")

	if pn.Extension() != "png" {
		t.Errorf("got %s expected png", pn.Extension())
	}
}

func TestSubsetRune(t *testing.T) {
	var tt = []struct {
		f string
		s rune
		t rune
		e string
	}{
		{`"DATA"`, '"', '"', "DATA"},
	}

	for _, d := range tt {
		c := subsetRune(d.f, d.s, d.t)

		if !strings.Contains(c, d.e) {
			t.Errorf("expected %s got %s", d.e, c)
		}
	}
}

func TestPageExtension(t *testing.T) {
	var tt = []struct {
		i string
		o string
	}{
		{"test", "jsx"},
		{"test.jsx", "jsx"},
	}

	for i, d := range tt {
		got := pageExtension(d.i)

		if got != d.o {
			t.Errorf("(%d) expected %s got %s", i, d.o, got)
		}
	}
}

func TestExtractJsTokenName(t *testing.T) {
	var tt = []struct {
		i string
		t JSToken
		o string
	}{
		{"export function thing() {}", FuncToken, "thing"},
		{"let cats = '2'", LetToken, "cats"},
	}

	for i, d := range tt {
		got, err := extractJSTokenName(d.i, d.t)
		if err != nil {
			t.Errorf("error should not have been thrown")
		}

		if got != d.o {
			t.Errorf("(%d) expected %s got %s", i, d.o, got)
		}
	}
}

func TestParseArgs(t *testing.T) {
	var tt = []struct {
		i string
		l int
	}{
		{"let thing = '5'", 0},
		{"let thing = (ice) => {}", 1},
		{"let obthing = ({}, {}, ince) => {}", 3},
		{"function ({}, ince = 5) => {}", 2},
	}

	for i, d := range tt {
		got, err := parseArgs(d.i)
		if err != nil {
			t.Errorf("error should not have been thrown")
		}

		if len(got) != d.l {
			t.Errorf("(%d) expected %d got %d", i, d.l, len(got))
		}
	}
}
