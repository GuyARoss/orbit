package jsparse

import (
	"errors"
	"testing"
)

func Test_formatImportLine_DefaultPkg(t *testing.T) {
	p := DefaultJSDocument{webDir: "test"}
	got := p.formatImportLine("import React from 'react'")
	expected := "import React from 'react'"

	if got.FinalStatement != expected {
		t.Errorf("got %s, expected %s", got.FinalStatement, expected)
	}
}

func Test_formatImportLine_AlternateEOLChar(t *testing.T) {
	p := DefaultJSDocument{webDir: "test"}
	got := p.formatImportLine("import { withMemo } from 'video-react';")
	expected := "import { withMemo } from 'video-react';"

	if got.FinalStatement != expected {
		t.Errorf("got %s, expected %s", got.FinalStatement, expected)
	}
}

func Test_formatImportLine_DefaultLocal(t *testing.T) {
	p := DefaultJSDocument{webDir: "test"}
	got := p.formatImportLine("import React from '../react'")
	expected := "import React from '../../../test/react.jsx'"

	if got.FinalStatement != expected {
		t.Errorf("got %s, expected %s", got.FinalStatement, expected)
	}
}

func Test_formatImportLine_AlternativeStrChar(t *testing.T) {
	p := DefaultJSDocument{webDir: "test"}
	got := p.formatImportLine("import React from \"../react\"")
	expected := "import React from '../../../test/react.jsx'"

	if got.FinalStatement != expected {
		t.Errorf("got %s, expected %s", got.FinalStatement, expected)
	}
}

func Test_formatImportLine_ConstLocal(t *testing.T) {
	p := DefaultJSDocument{webDir: "test"}
	got := p.formatImportLine("import { tool } from '../tools/test'")
	expected := "import { tool } from '../../../test/tools/test.jsx'"

	if got.FinalStatement != expected {
		t.Errorf("got %s, expected %s", got.FinalStatement, expected)
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

func Test_cleanExportDefaultName(t *testing.T) {
	name, _ := extractDefaultExportName("export default Test")
	if name != "Test" {
		t.Errorf("expected %s got %s \n", "Test", name)
	}

	_, err := extractDefaultExportName("export default () => {}")
	if !errors.Is(ErrFunctionExport, err) {
		t.Error("expected function export to raise error")
	}

	_, err = extractDefaultExportName("export default test")
	if !errors.Is(ErrExportNotCapitalized, err) {
		t.Error("expected non capitalized component to raise exception")
	}
}

func Test_cleanExportDefaultNameExtraSpace(t *testing.T) {
	name, err := extractDefaultExportName("export default SomethingCool  ")
	if err != nil {
		t.Error("error exception should not be thrown", err)
	}

	if name != "SomethingCool" {
		t.Errorf("expected %s got %s", "SomethingCool", name)
	}
}

func Test_defaultPageName(t *testing.T) {
	pn := defaultPageName("thing_stuff")
	if pn != "ThingStuff" {
		t.Error("default page name mismatch")
	}

	pn = defaultPageName("sff_m.js")
	if pn != "SffM" {
		t.Error("default page name mismatch")
	}
}

func Test_extension(t *testing.T) {
	pn := &DefaultJSDocument{
		pageDir: "thing.png",
	}

	if pn.Extension() != "png" {
		t.Errorf("got %s expected png", pn.Extension())
	}
}
