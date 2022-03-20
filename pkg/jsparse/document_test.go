// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package jsparse

import (
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
		{`import React from "../react"`, "import React from '../../../test/react.jsx'"},
		{"import { tool } from '../tools/test'", "import { tool } from '../../../test/tools/test.jsx'"},
		{"import { tool } from '../tools/test.js'", "import { tool } from '../../../test/tools/test.js'"},
	}

	p := DefaultJSDocument{webDir: "test", pageDir: "./thing/apple.js"}

	for i, c := range tt {
		got := p.formatImportLine(c.i)

		if fsutils.NormalizePath(c.o) != got.FinalStatement {
			t.Errorf("(%d) expected %s got %s \n", i, fsutils.NormalizePath(c.o), got.FinalStatement)
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
			extension: "jsx",
			other:     []string{"some random text"},
		}},
		{"export default Thing", DefaultJSDocument{
			extension: "jsx",
			name:      "Thing",
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
