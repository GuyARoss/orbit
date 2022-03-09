// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// LICENSE file in the root directory of this source tree.
package srcpack

import (
	"fmt"
	"testing"

	"github.com/GuyARoss/orbit/pkg/log"
)

func TestSyncHook(t *testing.T) {
	logger := log.NewDefaultLogger()

	sh := NewSyncHook(logger)

	defer sh.Close()

	for i := 0; i < 4; i++ {
		sh.WrapFunc(fmt.Sprintf("thing_%d", i), func() {

		})
	}
}
