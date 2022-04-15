// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package jsparse

import (
	"testing"
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
		{`import React from "../react"`, "import React from '../../../test/react.jsx'"},
		{"import { tool } from '../tools/test'", "import { tool } from '../../../test/tools/test.jsx'"},
		{"import { tool } from '../tools/test.js'", "import { tool } from '../../../test/tools/test.js'"},
		{"import 'thing.css'", "import 'thing.css'"},
	}

	p := DefaultJSDocument{webDir: "test", pageDir: "./thing/apple.js"}

	for i, c := range tt {
		got := p.formatImportLine(c.i)

		if c.o != got.FinalStatement {
			t.Errorf("(%d) expected %s got %s \n", i, c.o, got.FinalStatement)
		}
	}
}

func TestVerifyPath(t *testing.T) {
	var tt = []struct {
		i string
		o string
	}{
		{"./thing.jsx", "./thing.jsx"},
		{".thing.jsx", "./thing.jsx"},
		{"/cake.jsx", "./cake.jsx"},
		{"cake.jsx", "./cake.jsx"},
		{"../../cake.jsx", "../../cake.jsx"},
	}

	for i, d := range tt {
		got := verifyPath(d.i)
		if got != d.o {
			t.Errorf("(%d) expected '%s' got '%s'", i, d.o, got)
		}

	}
}

func TestTokenizeLine(t *testing.T) {
	var tt = []struct {
		i          string
		o          DefaultJSDocument
		exportName string
	}{
		{"import Thing from 'thing'", DefaultJSDocument{
			imports: []*ImportDependency{
				{"import Thing from 'thing'", "", ModuleImportType},
			},
		}, ""},
		{"// some random text", DefaultJSDocument{
			extension: "jsx",
			other:     []string{"some random text"},
		}, ""},
		{"export default Thing", DefaultJSDocument{
			extension: "jsx",
		}, "Thing"},
	}

	for i, d := range tt {
		cdoc := NewEmptyDocument()
		got := cdoc.tokenizeLine(d.i)

		if got != nil {
			t.Error("did not expect error during line tokenization")
			return
		}

		if cdoc.defaultExport.Name != d.exportName {
			t.Errorf("(%d) expected name %s got %s", i, d.exportName, cdoc.defaultExport.Name)
			return
		}

		if len(cdoc.imports) != len(d.o.imports) {
			t.Errorf("(%d) import missmatch", i)
			return
		}

		if len(cdoc.other) != len(d.o.other) {
			t.Errorf("(%d) other missmatch", i)
			return
		}

		cdoc = &DefaultJSDocument{}
	}
}

func TestTokenizeLine_DetectExport(t *testing.T) {
	cdoc := NewEmptyDocument()
	err := cdoc.tokenizeLine("function Thing() {}")
	if err != nil {
		t.Errorf("error occurred %s", err)
		return
	}

	err = cdoc.tokenizeLine("export default Thing")
	if err != nil {
		t.Errorf("error occurred %s", err)
		return
	}

	if cdoc.defaultExport == nil {
		t.Error("did not expect default export to be nil")
		return
	}

	if len(cdoc.defaultExport.Args) != 0 {
		t.Error("did not expect args to be present on resource default export")
		return
	}
}
