// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package jsparse

import (
	"testing"
)

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
