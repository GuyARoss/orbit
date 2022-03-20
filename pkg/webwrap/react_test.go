// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package webwrap

import (
	"errors"
	"testing"

	"github.com/GuyARoss/orbit/pkg/jsparse/mock"
)

func TestApplyReact_Error(t *testing.T) {
	tt := []struct {
		doc *mock.MockJsDocument
		err error
	}{
		{mock.NewMockJSDocument("", "blah"), ErrInvalidComponent},
		{mock.NewMockJSDocument("lowercaseThing", "jsx"), ErrComponentExport},
	}

	r := &ReactWebWrapper{}

	for i, d := range tt {
		_, err := r.Apply(d.doc)

		if !errors.Is(err, d.err) {
			t.Errorf("(%d) got incorrect error", i)
		}
	}
}

func TestApplyReact(t *testing.T) {
	r := &ReactWebWrapper{}

	p, err := r.Apply(mock.NewMockJSDocument("Thing", "jsx"))
	if err != nil {
		t.Error("should not expect error during valid jsx parsing")
	}

	if p.Name() != "Thing" {
		t.Errorf("expected name 'Thing' got '%s'", p.Name())
	}
}
