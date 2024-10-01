// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package internal

import (
	"os"
	"testing"
)

func TestBuild_NoPaths(t *testing.T) {
	tdir := t.TempDir()

	os.Mkdir(tdir+"/pages", 0777)

	opts := &BuildOpts{
		PackageName:    "",
		OutDir:         tdir,
		ApplicationDir: tdir,
		Mode:           "development",
		NodeModulePath: "",
		PublicPath:     "./thing",
		RequiredDirs:   []string{},
		NoWrite:        true,
	}

	final, err := Build(opts)
	if err != nil {
		t.Errorf("should not fail during build '%s'", err)
		return
	}

	if len(final) != 0 {
		t.Errorf("should not not exceed 0 members got '%d'", len(final))
		return
	}
}
