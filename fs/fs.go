package fs

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func copyFile(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	defer in.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func copyDir(dir string, baseDir string, outDir string) []string {
	copiedDirs := make([]string, 0)

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		if entry.Name() == ".orbit" {
			continue
		}

		if entry.IsDir() {
			copied := copyDir(filepath.Join(dir, entry.Name()), baseDir, outDir)

			for p := range copied {
				copiedDirs = append(copiedDirs, copied[p])
			}

			continue
		}

		sourcePath := filepath.Join(dir, entry.Name())
		ns := strings.Replace(dir, baseDir, "", 1)

		orbitDirPath := filepath.Join(outDir, ns)
		destPath := filepath.Join(orbitDirPath, entry.Name())

		if !doesDirExist(orbitDirPath) {
			os.Mkdir(orbitDirPath, 0755)
		}

		copiedDirs = append(copiedDirs, destPath)
		copyFile(sourcePath, destPath)
	}

	return copiedDirs
}

func doesDirExist(dir string) bool {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func appendFile(fileName string, content string) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(content); err != nil {
		return err
	}

	return nil
}
