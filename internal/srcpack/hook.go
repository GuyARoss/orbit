// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package srcpack

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/GuyARoss/orbit/pkg/log"
)

// hooks for logging the pre & post operations of the packing process.
type PackHook interface {
	Pre(filePath string)                       // "pre" runs before each component packing iteration
	Post(filepath string, elapsedTime float64) // "post" runs after each component packing iteration
	Close()
}

type SyncHook struct {
	postmap map[string]float64
	premap  map[string]bool
	logger  log.Logger

	order       *list.List
	initialized bool
	m           *sync.Mutex
}

func NewSyncHook(logger log.Logger) *SyncHook {
	return &SyncHook{
		logger:  logger,
		postmap: make(map[string]float64),
		premap:  make(map[string]bool),
		order:   list.New(),
		m:       &sync.Mutex{},
	}
}

func (s *SyncHook) Pre(filePath string) {
	f := s.order.Front()

	s.order.PushBack(filePath)

	// if nothing is in queue, then we can write.
	if f == nil && !s.initialized {
		s.premap[filePath] = true
		s.logger.Info(fmt.Sprintf("(1) bundling %s → ...", filePath))
	}

	s.initialized = true
}

func (s *SyncHook) Post(filePath string, elapsedTime float64) {
	current := s.order.Front()

	// if the filepath is the current queued filepath
	// we write the output & remove the item from the queue.
	if current == nil && !s.premap[filePath] && current.Value == filePath {
		s.premap[filePath] = true
		s.logger.Info(fmt.Sprintf("(2) bundling %s → ...", filePath))
		s.logger.Success(fmt.Sprintf("completed in %fs\n", elapsedTime))
		s.order.Remove(current)
	}

	// this can either be referenced later for metrics or be used in the case that
	// the queue does not resolve all of the items.
	s.postmap[filePath] = elapsedTime
}

func (s *SyncHook) WrapFunc(filepath string, do func()) {
	starttime := time.Now()
	s.m.Lock()
	s.Pre(filepath)

	do()

	s.Post(filepath, time.Since(starttime).Seconds())
	s.postmap[filepath] = time.Since(starttime).Seconds()
	s.m.Unlock()
}

func (s *SyncHook) Close() {
	for s.order.Len() > 0 {
		c := s.order.Front()
		s.order.Remove(c)

		filename := c.Value.(string)

		// if we have not yet proccessed the filename then do it
		if !s.premap[filename] {
			s.logger.Info(fmt.Sprintf("(1) bundling %s → ...", filename))
		}
		s.logger.Success(fmt.Sprintf("completed in %fs\n", s.postmap[filename]))
	}
}
