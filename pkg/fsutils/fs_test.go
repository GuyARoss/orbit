package fsutils

import (
	"path/filepath"
	"testing"

	"github.com/GuyARoss/orbit/pkg/fsutils"
)

func TestCondenseFilePath_LongPath(t *testing.T) {
	path := filepath.Clean(fsutils.NormalizePath(".orbit/base/web/pages/home.jsx"))

	got := condenseFilePath(path)
	if got != fsutils.NormalizePath(".orbit/base/pages/home.jsx") {
		t.Errorf("expected: %s, got %s", fsutils.NormalizePath(".orbit/pages/home.jsx"), got)
	}
}

func TestCondenseDirPath_LongDir(t *testing.T) {
	path := filepath.Clean(fsutils.NormalizePath(".orbit/base/web/pages"))

	got := condenseDirPath(path)
	if got != fsutils.NormalizePath(".orbit/base/pages") {
		t.Errorf("expected: %s, got %s", fsutils.NormalizePath(".orbit/base/pages"), got)
	}
}

func TestCondenseDirPath_ShortDir(t *testing.T) {
	path := filepath.Clean(fsutils.NormalizePath(".orbit/base/pages"))

	got := condenseDirPath(path)
	if got != fsutils.NormalizePath(".orbit/base/pages") {
		t.Errorf("expected: %s, got %s", fsutils.NormalizePath(".orbit/base/pages"), got)
	}
}
