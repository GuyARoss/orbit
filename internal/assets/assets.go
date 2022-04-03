// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/GuyARoss/orbit/pkg/fsutils"
)

//go:embed embed/*
var content embed.FS

type AssetKey string

const (
	WebPackConfig  AssetKey = "base.config.js"
	HotReload      AssetKey = "hotreload.js"
	Tests          AssetKey = "orbit_test.go"
	PrimaryPackage AssetKey = "orbit.go"
	SSRProtoFile   AssetKey = "com.proto"
)

func WriteFile(toDir string, f fs.DirEntry) error {
	newFile, err := os.Create(fsutils.NormalizePath(fmt.Sprintf("%s/%s", toDir, f.Name())))
	if err != nil {
		return err
	}

	data, err := content.ReadFile(fsutils.NormalizePath(fmt.Sprintf("embed/%s", f.Name())))
	if err != nil {
		return err
	}

	_, err = newFile.Write(data)
	if err != nil {
		return err
	}

	return nil
}

type AssetFileReader struct {
	dirEntry fs.DirEntry
}

func (s *AssetFileReader) Read() (fs.File, error) {
	return content.Open(fsutils.NormalizePath(fmt.Sprintf("embed/%s", s.dirEntry.Name())))
}

type AssetMap map[AssetKey]fs.DirEntry

func (m AssetMap) AssetEntry(key AssetKey) fs.DirEntry {
	return m[key]
}

func (m AssetMap) AssetKey(key AssetKey) *AssetFileReader {
	return &AssetFileReader{dirEntry: m[key]}
}

func AssetKeys() (AssetMap, error) {
	mp := make(map[AssetKey]fs.DirEntry)

	files, err := content.ReadDir("embed")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".template") {
			continue
		}

		mp[AssetKey(file.Name())] = file
	}

	return mp, nil
}
