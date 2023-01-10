// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package srcpack

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/GuyARoss/orbit/pkg/log"
	"github.com/GuyARoss/orbit/pkg/webwrap"
)

// hooks for logging the pre & post operations of the packing process.
type PackHook interface {
	Pre(filePath string)                       // "pre" runs before each component packing iteration
	Post(filepath string, elapsedTime float64) // "post" runs after each component packing iteration
	Close()
}

type SyncHook struct {
	logger log.Logger

	m *sync.Mutex
}

func NewSyncHook(logger log.Logger) *SyncHook {
	return &SyncHook{
		logger: logger,
		m:      &sync.Mutex{},
	}
}

func (s *SyncHook) WrapFunc(filepath string, do func() *webwrap.WrapStats) {
	starttime := time.Now()
	stats := do()

	s.m.Lock()
	if stats == nil {
		s.m.Unlock()
		s.logger.Error(fmt.Sprintf("failed to bundle '%s'", filepath))
		return
	}
	elapsed := strings.Split(fmt.Sprintf("%f", time.Since(starttime).Seconds()), ".")
	s.logger.Info(fmt.Sprintf("%s - %s.%ss", filepath, elapsed[0], elapsed[1][0:1]))
	s.logger.Info(fmt.Sprintf("[web: %s, bundler: %s]\n", stats.WebVersion, stats.Bundler))
	s.m.Unlock()
}
