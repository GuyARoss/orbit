// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package libout

import (
	"testing"

	"github.com/GuyARoss/orbit/pkg/embedutils"
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

func TestEnvFile(t *testing.T) {
	f := &GOLibout{}
	loboutFile, err := f.EnvFile(&BundleGroup{
		pages: []*page{
			{
				name: "SomePage",
			},
			{
				name: "SomeSecondPage",
			},
		},
		wrapDocRender: make(map[string][]embedutils.FileReader),
		BundleGroupOpts: &BundleGroupOpts{
			BaseBundleOut: "SomeDirThing",
			PackageName:   "TestPackage",
			PublicDir:     "/directory/here",
			HotReloadPort: 2012,
		},
	})
	if err != nil {
		t.Error("did not expect error", err)
		return
	}

	got := len(loboutFile.(*GOLibFile).Body)
	if got != 992 {
		t.Errorf("got '%d', expected '%d'", got, 992)
		return
	}
}
