package internal

import (
	"os"
	"testing"
)

func TestBuild_NoPaths(t *testing.T) {
	tdir := t.TempDir()

	os.Mkdir(tdir+"/pages", 0777)

	opts := &BuildOpts{
		Packname:       "",
		OutDir:         tdir,
		WebDir:         tdir,
		Mode:           "development",
		NodeModulePath: "",
		PublicDir:      "./thing",
		Dirs:           []string{},
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
