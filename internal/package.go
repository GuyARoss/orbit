package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/GuyARoss/orbit/internal/assets"
	"github.com/GuyARoss/orbit/pkg/fsutils"
)

type PackageJSONTemplate struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Author       string            `json:"author"`
	License      string            `json:"license"`
	Description  string            `json:"description"`
	Dependencies map[string]string `json:"dependencies"`
}

func (p *PackageJSONTemplate) Write(path string) error {
	newFile, err := os.Create(path)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(p)
	if err != nil {
		return err
	}

	newFile.Write(jsonData)

	return nil
}

type FileStructureOpts struct {
	PackageName string
	OutDir      string
	Assets      []fs.DirEntry
}

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
