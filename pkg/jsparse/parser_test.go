package jsparse

import "testing"

func Test_formatImportLine_DefaultPkg(t *testing.T) {
	p := Page{webDir: "test"}
	got := p.formatImportLine("import React from 'react'")
	expected := "import React from 'react'"

	if got.FinalStatement != expected {
		t.Errorf("got %s, expected %s", got.FinalStatement, expected)
	}
}

func Test_formatImportLine_DefaultLocal(t *testing.T) {
	p := Page{webDir: "test"}
	got := p.formatImportLine("import React from '../react'")
	expected := "import React from '../../../test/react.jsx'"

	if got.FinalStatement != expected {
		t.Errorf("got %s, expected %s", got.FinalStatement, expected)
	}
}

func Test_formatImportLine_AlternativeStrChar(t *testing.T) {
	p := Page{webDir: "test"}
	got := p.formatImportLine("import React from \"../react\"")
	expected := "import React from '../../../test/react.jsx'"

	if got.FinalStatement != expected {
		t.Errorf("got %s, expected %s", got.FinalStatement, expected)
	}
}

func Test_formatImportLine_ConstLocal(t *testing.T) {
	p := Page{webDir: "test"}
	got := p.formatImportLine("import { tool } from '../tools/test'")
	expected := "import { tool } from '../../../test/tools/test.jsx'"

	if got.FinalStatement != expected {
		t.Errorf("got %s, expected %s", got.FinalStatement, expected)
	}
}
