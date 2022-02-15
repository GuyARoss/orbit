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
