package fs

import (
	"fmt"
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
	if err != nil {
		return err
	}
	defer in.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

type CopyResults struct {
	BaseDir string
	CopyDir string
}

func pathDelimiter(path string) string {
	if strings.Contains(path, "//") {
		return "//"
	}

	return "/"
}

func condenseFilePath(filePath string) string {
	pathType := pathDelimiter(filePath)
	spt := strings.Split(filePath, pathType)

	return fmt.Sprintf("%s%s%s", strings.Join(spt[0:2], pathType), pathType, strings.Join(spt[len(spt)-2:], pathType))
}

func condenseDirPath(dirPath string) string {
	pathType := pathDelimiter(dirPath)
	spt := strings.Split(dirPath, pathType)

	return fmt.Sprintf("%s%s%s", strings.Join(spt[0:2], pathType), pathType, spt[len(spt)-1])
}

func CopyDir(dir string, baseDir string, outDir string, condense bool) []*CopyResults {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	copiedDirs := make([]*CopyResults, 0)
	for _, entry := range entries {
		if entry.Name() == ".orbit" {
			continue
		}

		// @@todo(guy): implement a .orbitignore?
		if entry.Name() == "node_modules" {
			continue
		}

		if entry.IsDir() {
			copied := CopyDir(filepath.Join(dir, entry.Name()), baseDir, outDir, condense)

			for p := range copied {
				copiedDirs = append(copiedDirs, &CopyResults{
					BaseDir: copied[p].BaseDir,
					CopyDir: copied[p].CopyDir,
				})
			}

			continue
		}

		sourcePath := filepath.Join(dir, entry.Name())
		ns := strings.Replace(dir, baseDir, "", 1)

		dirPath := filepath.Join(outDir, ns)
		destPath := filepath.Join(dirPath, entry.Name())
		if condense {
			destPath = condenseFilePath(destPath)
			dirPath = condenseDirPath(dirPath)
		}

		if !DoesDirExist(dirPath) {
			os.Mkdir(dirPath, 0755)
		}

		copiedDirs = append(copiedDirs, &CopyResults{
			BaseDir: filepath.Join(dir, entry.Name()),
			CopyDir: destPath,
		})
		copyFile(sourcePath, destPath)
	}

	return copiedDirs
}

func DoesDirExist(dir string) bool {
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}

func DirFiles(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	simpleFiles := make([]string, len(files))
	for idx, file := range files {
		// @@todo add support for non-shallow directories
		if !file.IsDir() {
			simpleFiles[idx] = file.Name()
		}
	}

	return simpleFiles
}
