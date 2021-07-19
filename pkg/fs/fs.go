package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func copyDir(dir string, baseDir string, outDir string) []*CopyResults {
	copiedDirs := make([]*CopyResults, 0)

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		if entry.Name() == ".orbit" {
			continue
		}

		// @@todo(guy): implement a .orbitignore?
		if entry.Name() == "node_modules" {
			continue
		}

		if entry.IsDir() {
			copied := copyDir(filepath.Join(dir, entry.Name()), baseDir, outDir)

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

		orbitDirPath := filepath.Join(outDir, ns)
		destPath := filepath.Join(orbitDirPath, entry.Name())

		if !DoesDirExist(orbitDirPath) {
			os.Mkdir(orbitDirPath, 0755)
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

type DirWatch struct {
	FileChange chan string
	Error      chan error
}

func DirectoryWatch(dir string) (*DirWatch, error) {
	fChangeChan := make(chan string)
	errorChan := make(chan error)

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		go func(fChan chan string, eChan chan error, entry os.FileInfo) {
			path := filepath.Join(dir, entry.Name())

			err := watchFile(path, fChan)
			if err != nil {
				eChan <- err
			}
		}(fChangeChan, errorChan, entry)
	}

	return &DirWatch{
		FileChange: fChangeChan,
		Error:      errorChan,
	}, nil
}

func watchFile(filePath string, fileChange chan string) error {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			fmt.Println("file changed")
			fileChange <- filePath
		}

		time.Sleep(1 * time.Second)
	}
}
