// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package fsutils

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestDirFiles(t *testing.T) {
	// in this test the max depth for dir files is 2
	paths := []string{"/cat", "/cat2", "/cat2/testdir"}

	tmp := fmt.Sprintf("%s%c%s", os.TempDir(), os.PathSeparator, uuid.NewString())
	os.Mkdir(tmp, 0755)
	os.Create(tmp + "/thing.delete")

	for i, path := range paths {
		np := fmt.Sprintf("%s/%s", tmp, path)
		os.Mkdir(np, 0755)

		f, _ := os.Create(np + "/thing.delete")
		f.Close()

		paths[i] = np
	}

	resp := DirFiles(tmp)

	if len(resp) != 3 {
		t.Errorf("expected %d got %d", 3, len(resp))
	}
}
