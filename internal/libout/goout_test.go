package libout

import (
	"testing"
)

func TestMergeImports(t *testing.T) {
	p := &parsedGoFile{
		Imports: []string{
			"1",
			"2",
			"3",
		},
	}

	p.MergeImports([]string{
		"4",
		"5",
		"6",
	})

	cm := []string{
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
	}

	for i, v := range p.Imports {
		if cm[i] != v {
			t.Errorf("import %d expected %s got %s", i, cm[i], v)
		}
	}
}

func TestGoParserParseLine(t *testing.T) {
	p := &goParser{
		imports:         make([]string, 0),
		contextOfImport: false,
		softImports:     make(map[string]bool),
	}

	f := []struct {
		in              string
		out             string
		contextOfImport bool
		importIndex     int
	}{
		{
			out:             "",
			in:              `import "abs"`,
			contextOfImport: false,
			importIndex:     1,
		},
		{
			out:             "",
			in:              `import "cat"`,
			contextOfImport: false,
			importIndex:     2,
		},
		{
			out:             "",
			in:              `import (`,
			contextOfImport: true,
			importIndex:     2,
		},
		{
			out:             "",
			in:              `"face"`,
			contextOfImport: true,
			importIndex:     3,
		},
	}

	for i, ff := range f {
		o := p.parseLine(ff.in)
		if o != ff.out {
			t.Errorf("(%d) out expected '%s' got '%s'", i, ff.out, o)
		}
		if p.contextOfImport != ff.contextOfImport {
			t.Errorf("(%d) contextOfImport expected '%t' got '%t'", i, ff.contextOfImport, p.contextOfImport)
		}
		if len(p.imports) != ff.importIndex {
			t.Errorf("(%d) import count expected '%d' got '%d'", i, ff.importIndex, len(p.imports))
		}
	}
}

func TestGoParsedFileSerialize(t *testing.T) {
	p := &parsedGoFile{
		Body: "",
		Imports: []string{
			"1",
			"2",
			"3",
			"4",
		},
	}

	expected := `import (
	"1"
	"2"
	"3"
	"4"
)
`
	o := p.Serialize()
	if o != expected {
		t.Errorf("expected '%s' got '%s'", expected, o)
	}
}
