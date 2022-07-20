// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package parseerror

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	e := New("thing", "name")

	if e == nil {
		t.Errorf("error was not correcly created")
	}
}

func TestFromError(t *testing.T) {
	e := FromError(nil, "tea")
	if e != nil {
		t.Errorf("error should be nil, if error is nil")
	}

	e = FromError(errors.New("thing"), "123")
	if e == nil {
		t.Errorf("error was not correcly created")
	}
}
