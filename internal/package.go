// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package internal

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/GuyARoss/orbit/internal/assets"
	"github.com/GuyARoss/orbit/internal/srcpack"
)

// PackageJSONTemplate struct for nodejs package.json file.
type PackageJSONTemplate struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Author       string            `json:"author"`
	License      string            `json:"license"`
	Description  string            `json:"description"`
	Dependencies map[string]string `json:"dependencies"`
}

// Write creates a new package.json to the provided path
func (p *PackageJSONTemplate) Write(path string) error {
	newFile, err := os.Create(path)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(p)
	if err != nil {
		return err
	}

	_, err = newFile.Write(jsonData)

	return err
}

// FileStructureOpts reqired structure options for the file structure creation
type FileStructure struct {
	PackageName string
	OutDir      string
	Assets      []fs.DirEntry
	Dist        []fs.DirEntry
	Mkdirs      []string
}

var dirs = []string{
	".orbit", ".orbit/base", ".orbit/base/pages",
	".orbit/dist", ".orbit/assets",
}

// Cleanup removes all required file paths
func (s *FileStructure) Cleanup() error {
	for _, dir := range dirs {
		err := os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}

	for _, dir := range s.Mkdirs {
		os.RemoveAll(dir)
	}

	return nil
}

// Make creates the foundation for orbits file structure
func (s *FileStructure) Make() error {
	for _, dir := range s.Mkdirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.Mkdir(dir, os.ModePerm); err != nil {
				return err
			}
		}
	}

	if _, err := os.Stat(fmt.Sprintf("%s/%s", s.OutDir, s.PackageName)); os.IsNotExist(err) {
		err := os.Mkdir(fmt.Sprintf("%s/%s", s.OutDir, s.PackageName), os.ModePerm)
		if err != nil {
			return err
		}
	}

	for _, dir := range dirs {
		_, err := os.Stat(dir)
		if errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(dir, 0755)
			if err != nil {
				return err
			}
		}
	}

	for _, a := range s.Assets {
		if err := assets.WriteFile(".orbit/assets", a); err != nil {
			return err
		}
	}

	for _, a := range s.Dist {
		if err := assets.WriteFile(".orbit/dist", a); err != nil {
			return err
		}
	}

	return nil
}

// CachedEnvFromFile creates a new cached environment provided a file path
// for this method, we prefer using a single pass file reader over something like
// reflection due to the speed constraints of reflection
func CachedEnvFromFile(path string) (srcpack.CachedEnvKeys, error) {
	// TODO(language-support): if we plan to add support for another output langauge, this
	// function needs to validate the extension to determine parsing method.
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	k := make(srcpack.CachedEnvKeys)

	current := ""
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "orbit:page") {
			current = strings.Split(line, " ")[2]
			continue
		}

		// TODO(language-support): to provide extensibility this should be specific to the language parser
		// rather than a hardcoded value.
		if current != "" && strings.Contains(line, "PageRender") {
			k[current] = strings.ReplaceAll(strings.Split(line, " ")[3], `"`, ``)

			current = ""
		}
	}

	return k, nil
}
