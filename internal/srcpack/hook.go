package srcpack

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/GuyARoss/orbit/pkg/log"
)

// hooks for logging the pre & post operations of the packing process.
type Hooks interface {
	Pre(filePath string)                       // "pre" runs before each component packing iteration
	Post(filepath string, elapsedTime float64) // "post" runs after each component packing iteration
	Close()
}

type SyncHook struct {
	postmap map[string]float64
	premap  map[string]bool
	logger  log.Logger

	order *list.List
}

func NewSyncHook(logger log.Logger) *SyncHook {
	return &SyncHook{
		logger:  logger,
		postmap: make(map[string]float64),
		premap:  make(map[string]bool),
		order:   list.New(),
	}
}

func (s *SyncHook) Pre(filePath string) {
	// if nothing is in queue, then we can write.
	if s.order.Front() == nil {
		s.premap[filePath] = true
		s.logger.Info(fmt.Sprintf("bundling %s → ...", filePath))
	}

	s.order.PushBack(filePath)
}

func (s *SyncHook) Post(filePath string, elapsedTime float64) {
	current := s.order.Front()

	if !s.premap[filePath] {
		s.logger.Info(fmt.Sprintf("bundling %s → ...", filePath))
	}

	if s.postmap[filePath] != 0 {
		s.logger.Success(fmt.Sprintf("completed in %fs\n", elapsedTime))
	} else {
		if current != nil && current.Value == filePath {
			s.logger.Success(fmt.Sprintf("completed in %fs\n", elapsedTime))
		}
	}

	if current != nil {
		s.order.Remove(current)
	}

	s.postmap[filePath] = elapsedTime
}

func (s *SyncHook) WrapFunc(filepath string, do func()) {
	starttime := time.Now()
	m := sync.Mutex{}

	m.Lock()
	s.Pre(filepath)
	m.Unlock()

	do()

	m.Lock()
	s.Post(filepath, time.Since(starttime).Seconds())
	s.postmap[filepath] = time.Since(starttime).Seconds()

	m.Unlock()
}

func (s *SyncHook) Close() {
	for s.order.Len() > 0 {
		c := s.order.Front()
		s.order.Remove(c)

		valStr := c.Value.(string)
		s.Post(valStr, s.postmap[valStr])
	}
}
