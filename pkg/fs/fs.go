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

func condenseFilePath(filePath string) string {
	spt := strings.Split(filePath, "\\")

	return fmt.Sprintf("%s\\%s", strings.Join(spt[0:2], "\\"), strings.Join(spt[len(spt)-2:], "\\"))
}

func condenseDirPath(dirPath string) string {
	spt := strings.Split(dirPath, "\\")

	return fmt.Sprintf("%s\\%s", strings.Join(spt[0:2], "\\"), spt[len(spt)-1])
}

func copyDir(dir string, baseDir string, outDir string, condense bool) []*CopyResults {
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
			copied := copyDir(filepath.Join(dir, entry.Name()), baseDir, outDir, condense)

			for p := range copied {
				copiedDirs = append(copiedDirs, &CopyResults{
					BaseDir: filepath.Join(dir, entry.Name()),
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
