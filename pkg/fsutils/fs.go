// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package fsutils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func NormalizePath(path string) string {
	return strings.ReplaceAll(path, "/", fmt.Sprintf("%c", os.PathSeparator))
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
			simpleFiles[idx] = fmt.Sprintf("%s/%s", dir, file.Name())
		}
	}

	return simpleFiles
}
