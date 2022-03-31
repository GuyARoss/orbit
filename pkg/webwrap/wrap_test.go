// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package webwrap_test

import (
	"testing"

	"github.com/GuyARoss/orbit/pkg/webwrap"
	"github.com/GuyARoss/orbit/pkg/webwrap/mock"
)

func TestFirstMatch(t *testing.T) {
	l := webwrap.JSWebWrapperList([]webwrap.JSWebWrapper{&mock.MockWrapper{true, false}})

	if l.FirstMatch("tst") == nil {
		t.Errorf("expected match ")
	}
}
