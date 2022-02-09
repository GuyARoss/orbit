package srcpack

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/GuyARoss/orbit/pkg/log"
)

// PackHooks
// passing of "per" & "post" hooks for our iterative packing method "PackPages".
type Hooks interface {
	Pre(filePath string)                       // "pre" runs before each component packing iteration
	Post(filepath string, elapsedTime float64) // "post" runs after each component packing iteration
	Finalize()
}

type SyncHook struct {
	keys    []string
	postmap map[string]float64
	logger  log.Logger

	l *list.List
}

func NewSyncHook(logger log.Logger) *SyncHook {
	return &SyncHook{
		logger:  logger,
		postmap: make(map[string]float64),
		keys:    make([]string, 0),
		l:       list.New(),
	}
}

func (s *SyncHook) Pre(filePath string) {
	s.logger.Info(fmt.Sprintf("bundling %s → ...", filePath))

	s.l.PushBack(filePath)
}

func (s *SyncHook) Post(filePath string, elapsedTime float64) {
	current := s.l.Front()

	if s.postmap[filePath] != 0 {
		s.logger.Success(fmt.Sprintf("completed in %fs\n", elapsedTime))
	} else {
		if current != nil && current.Value == filePath {
			s.logger.Success(fmt.Sprintf("completed in %fs\n", elapsedTime))
		}
	}

	if current != nil {
		s.l.Remove(current)
	}

	s.postmap[filePath] = elapsedTime
}

func (s *SyncHook) WrapFunc(filepath string, do func()) {
	starttime := time.Now()
	m := sync.Mutex{}

	m.Lock()
	s.keys = append(s.keys, filepath)
	s.Pre(filepath)
	m.Unlock()

	do()

	m.Lock()
	s.Post(filepath, time.Since(starttime).Seconds())
	s.postmap[filepath] = time.Since(starttime).Seconds()

	m.Unlock()
}

func (s *SyncHook) Finalize() {
	for s.l.Len() > 0 {
		c := s.l.Front()
		s.l.Remove(c)

		valStr := c.Value.(string)
		s.Post(valStr, s.postmap[valStr])
	}
}
