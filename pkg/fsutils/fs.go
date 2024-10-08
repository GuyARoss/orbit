// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package fsutils

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func depthFiles(dir string, maxDepth int, depth int) []string {
	if depth == maxDepth {
		return []string{}
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	simpleFiles := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			simpleFiles = append(simpleFiles, depthFiles(fmt.Sprintf("%s/%s", dir, file.Name()), maxDepth, depth+1)...)
		} else {
			simpleFiles = append(simpleFiles, fmt.Sprintf("%s/%s", dir, file.Name()))
		}
	}

	return simpleFiles
}

func DirFiles(dir string) []string {
	return depthFiles(dir, 2, 0)
}

func LastPathIndex(path string) string {
	paths := strings.Split(path, "/")
	s := strings.Split(paths[len(paths)-1], ".")

	return s[0]
}

func CanNotReadFile(path string) bool {
	_, err := os.Stat(path)
	return err != nil
}

// CopyFile copies a file from the src to the destination.
// this requires that the destination directory exists.
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	return err
}
