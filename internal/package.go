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
	"github.com/GuyARoss/orbit/pkg/fsutils"
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
type FileStructureOpts struct {
	PackageName string
	OutDir      string
	Assets      []fs.DirEntry
}

// OrbitFileStructure creates the foundation for orbits file structure, this includes:
// 1. creating the out package directory
// 2. the ./.orbit file structure
// 3. asset directory
func OrbitFileStructure(s *FileStructureOpts) error {
	err := os.RemoveAll(".orbit/")
	if err != nil {
		return err
	}

	if _, err := os.Stat(fsutils.NormalizePath(fmt.Sprintf("%s/%s", s.OutDir, s.PackageName))); os.IsNotExist(err) {
		err := os.Mkdir(fsutils.NormalizePath(fmt.Sprintf("%s/%s", s.OutDir, s.PackageName)), os.ModePerm)
		if err != nil {
			return err
		}
	}

	dirs := []string{
		".orbit", fsutils.NormalizePath(".orbit/base"), fsutils.NormalizePath(".orbit/base/pages"),
		fsutils.NormalizePath(".orbit/dist"), fsutils.NormalizePath(".orbit/assets"),
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
		err = assets.WriteFile(fsutils.NormalizePath(".orbit/assets"), a)
		if err != nil {
			return err
		}
	}

	return nil
}

// CachedEnvFromFile creates a new cached environment provided a file path
// for this method, we prefer using a single pass file reader over something like
// reflection due to the speed constraints of reflection
func CachedEnvFromFile(path string) (srcpack.CachedEnvKeys, error) {
	// @@todo: if we plan to add support for another output langauge, this
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

		// @@todo: to provide extensibility this should be specific to the language parser
		// rather than a hardcoded value.
		if current != "" && strings.Contains(line, "PageRender") {
			k[current] = strings.ReplaceAll(strings.Split(line, " ")[3], `"`, ``)

			current = ""
		}
	}

	return k, nil
}
