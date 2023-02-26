// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package internal

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/GuyARoss/orbit/internal/assets"
)

func TestPackageJSONTemplateWrite(t *testing.T) {
	j := &PackageJSONTemplate{
		Name:    "test",
		Version: "0.0.0",
	}

	path := fmt.Sprintf("%s/%s", t.TempDir(), "template")

	err := j.Write(path)
	if err != nil {
		t.Error("err writing json", err)
		return
	}

	_, err = os.Stat(path)
	if err != nil {
		t.Error("file does not exist", err)
	}
}

func TestMakeDeleteFileStructure(t *testing.T) {
	d := t.TempDir()
	mdir := fmt.Sprintf("%s/%s", d, "things")
	s := &FileStructure{
		PackageName: "test",
		OutDir:      d,
		Assets:      make([]fs.DirEntry, 0),
		Dist:        make([]fs.DirEntry, 0),
		Mkdirs:      []string{mdir},
	}
	err := s.Make()
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = os.Stat(mdir)
	if err != nil {
		t.Error("mkdir should exist")
	}

	_, err = os.Stat(fmt.Sprintf("%s/%s", d, "test"))
	if err != nil {
		t.Error("package dir should exist")
	}

	err = s.Cleanup()
	if err != nil {
		t.Error("delete should be successful")
	}

	_, err = os.Stat(mdir)
	if !os.IsNotExist(err) {
		t.Error("Cleanup failure")
	}
}

func BenchmarkMakeFileStructure(b *testing.B) {
	ats, err := assets.AssetKeys()
	if err != nil {
		b.Errorf("unexpected error '%s'", err)
		return
	}

	tDir := b.TempDir()
	s := &FileStructure{
		PackageName: "thing",
		OutDir:      tDir,
		Assets: []fs.DirEntry{
			ats.AssetEntry(assets.WebPackConfig),
			ats.AssetEntry(assets.SSRProtoFile),
			ats.AssetEntry(assets.JsWebPackConfig),
			ats.AssetEntry(assets.WebPackSWCConfig),
		},
		Mkdirs: []string{},
	}

	if err = s.Make(); err != nil {
		b.Errorf("unexpected error '%s'", err)
		return
	}
}
