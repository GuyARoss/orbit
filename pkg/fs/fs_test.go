package fs

import (
	"path/filepath"
	"testing"
)

func TestCondenseFilePath(t *testing.T) {
	path := filepath.Clean(".orbit\\base\\web\\pages\\home.jsx")

	got := condenseFilePath(path)
	if got != ".orbit\\base\\pages\\home.jsx" {
		t.Errorf("expected: %s, got %s", ".orbit\\pages\\home.jsx", got)
	}
}

func TestCondenseDirPath(t *testing.T) {
	path := filepath.Clean(".orbit\\base\\web\\pages")

	got := condenseDirPath(path)
	if got != ".orbit\\base\\pages" {
		t.Errorf("expected: %s, got %s", ".orbit\\base\\pages", got)
	}
}
