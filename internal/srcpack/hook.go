package srcpack

import (
	"fmt"

	"github.com/GuyARoss/orbit/pkg/log"
)

// PackHooks
// passing of "per" & "post" hooks for our iterative packing method "PackPages".
type Hooks interface {
	Pre(filePath string)      // "pre" runs before each component packing iteration
	Post(elapsedTime float64) // "post" runs after each component packing iteration
}

type DefaultHook struct{}

func (s *DefaultHook) Pre(filePath string) {
	log.Info(fmt.Sprintf("bundling %s → ...", filePath))
}

func (s *DefaultHook) Post(elapsedTime float64) {
	log.Success(fmt.Sprintf("completed in %fs\n", elapsedTime))
}
