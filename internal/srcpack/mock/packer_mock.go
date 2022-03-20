// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package mock

import (
	"log"

	"github.com/GuyARoss/orbit/internal/srcpack"
)

type MockPacker struct {
	Components []srcpack.Component
}

func (m *MockPacker) PackMany(pages []string) ([]srcpack.PackComponent, error) { return nil, nil }
func (m *MockPacker) PackSingle(logger log.Logger, file string) (srcpack.PackComponent, error) {
	return &m.Components[0], nil
}
func (m *MockPacker) ReattachLogger(logger log.Logger) srcpack.Packer { return nil }
