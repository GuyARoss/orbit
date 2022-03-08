package jsparse

import (
	"errors"
	"strings"
	"testing"

	"github.com/GuyARoss/orbit/pkg/fsutils"
)

func TestFormatImportLine(t *testing.T) {
	tt := []struct {
		i string
		o string
	}{
		{"import React from 'react'", "import React from 'react'"},
		{"import Thing2 from './thing2'", "import Thing2 from '../../../thing/thing2.jsx'"},
		{"import { withMemo } from 'react-thing';", "import { withMemo } from 'react-thing';"},
		{"import React from '../react'", "import React from '../../../test/react.jsx'"},
		{"import React from \"../react\"", "import React from '../../../test/react.jsx'"},
		{"import { tool } from '../tools/test'", "import { tool } from '../../../test/tools/test.jsx'"},
	}

	p := DefaultJSDocument{webDir: "test", pageDir: "./thing/apple.js"}

	for i, c := range tt {
		got := p.formatImportLine(c.i)

		if fsutils.NormalizePath(c.o) != got.FinalStatement {
			t.Errorf("(%d) expected %s got %s \n", i, fsutils.NormalizePath(c.o), got.FinalStatement)
		}
	}
}

func Test_lineImportType(t *testing.T) {
	g := lineImportType(`import { thing } from "@test/util"`)
	if g != ModuleImportType {
		t.Error("expected module import type")
	}

	g = lineImportType(`import cat from "../../utils.jsx"`)
	if g != LocalImportType {
		t.Error("expected local import type")
	}
}

func TestCleanExportDefaultNameErrors(t *testing.T) {
	_, err := extractDefaultExportName("export default () => {}")
	if !errors.Is(ErrFunctionExport, err) {
		t.Error("expected function export to raise error")
	}

	_, err = extractDefaultExportName("export default test")
	if !errors.Is(ErrExportNotCapitalized, err) {
		t.Error("expected non capitalized component to raise exception")
	}
}

func TestCleanExportDefaultName(t *testing.T) {
	tt := []struct {
		i string
		o string
	}{
		{"export default Test", "Test"},
		{"export default SomethingCool  ", "SomethingCool"},
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

func TestDefaultPageName(t *testing.T) {
	pn := defaultPageName("thing_stuff")
	if pn != "ThingStuff" {
		t.Error("default page name mismatch")
	}

	pn = defaultPageName("sff_m.js")
	if pn != "SffM" {
		t.Error("default page name mismatch")
	}
}

func TestExtension(t *testing.T) {
	pn := &DefaultJSDocument{
		pageDir: "thing.png",
	}

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

func TestPathToken(t *testing.T) {
	var tt = []struct {
		i string
		o rune
	}{
		{`import thing from "thing.js"`, '"'},
		{"disinvalid", '"'},
	}

	for i, d := range tt {
		got := pathToken(d.i)

		if got != d.o {
			t.Errorf("(%d) expected %c got %c", i, d.o, got)
		}
	}
}

func TestPageExtension(t *testing.T) {
	var tt = []struct {
		i string
		o string
	}{
		{"test", ".jsx"},
		{"test.jsx", ""},
	}

	for i, d := range tt {
		got := pageExtension(d.i)

		if got != d.o {
			t.Errorf("(%d) expected %s got %s", i, d.o, got)
		}
	}
}

func TestTokenizeLine(t *testing.T) {
	var tt = []struct {
		i string
		o DefaultJSDocument
	}{
		{"import Thing from 'thing'", DefaultJSDocument{
			imports: []*ImportDependency{
				{"import Thing from 'thing'", "", ModuleImportType},
			},
		}},
		{"some random text", DefaultJSDocument{
			other: []string{"some random text"},
		}},
		{"export default Thing", DefaultJSDocument{
			name: "Thing",
		}},
	}

	cdoc := &DefaultJSDocument{}
	for i, d := range tt {
		got := cdoc.tokenizeLine(d.i)

		if got != nil {
			t.Error("did not expect error during line tokenization")
		}

		if cdoc.name != d.o.name {
			t.Errorf("(%d) expected name %s got %s", i, cdoc.name, d.o.name)
		}

		if len(cdoc.imports) != len(d.o.imports) {
			t.Errorf("(%d) import missmatch", i)
		}

		if len(cdoc.other) != len(d.o.other) {
			t.Errorf("(%d) other missmatch", i)
		}

		cdoc = &DefaultJSDocument{}
	}
}
