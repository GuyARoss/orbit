package assets

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"strings"
)

//go:embed embed/*
var content embed.FS

func WriteAssetsDir(toDir string) error {
	files, err := content.ReadDir("embed")
	if err != nil {
		return err
	}

	_, err = os.Stat(toDir)
	if errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(toDir, 0755)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".template") {
			continue
		}

		newFile, err := os.Create(fmt.Sprintf("%s/%s", toDir, file.Name()))
		if err != nil {
			return err
		}

		data, err := content.ReadFile(fmt.Sprintf("embed/%s", file.Name()))
		if err != nil {
			return err
		}

		_, err = newFile.Write(data)
		if err != nil {
			return err
		}
	}

	return nil
}
