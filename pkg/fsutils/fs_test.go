// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package fsutils

import (
	"testing"
)

func TestCondenseFilePath_LongPath(t *testing.T) {
	path := ".orbit/base/web/pages/home.jsx"

	got := condenseFilePath(path)
	if got != ".orbit/base/pages/home.jsx" {
		t.Errorf("expected: %s, got %s", ".orbit/pages/home.jsx", got)
	}
}

func TestCondenseDirPath_LongDir(t *testing.T) {
	path := ".orbit/base/web/pages"

	got := condenseDirPath(path)
	if got != ".orbit/base/pages" {
		t.Errorf("expected: %s, got %s", ".orbit/base/pages", got)
	}
}

func TestCondenseDirPath_ShortDir(t *testing.T) {
	path := ".orbit/base/pages"

	got := condenseDirPath(path)
	if got != ".orbit/base/pages" {
		t.Errorf("expected: %s, got %s", ".orbit/base/pages", got)
	}
}
