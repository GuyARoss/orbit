package dev

import (
	"fmt"

	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/pkg/fs"
)

type session struct {
	pageGenSettings *internal.GenPagesSettings
	sourceMap       map[string]*fs.PackedPage
}

func createSession(settings *internal.GenPagesSettings) (*session, error) {
	err := settings.CleanPathing()
	if err != nil {
		return nil, err
	}

	lib := settings.ApplyPages()

	sourceMap := make(map[string]*fs.PackedPage)
	for _, p := range lib.Pages {
		sourceMap[p.BaseDir] = p
	}

	return &session{
		settings, sourceMap,
	}, nil
}

func (s *session) executeChangeRequest(file string) {
	fmt.Printf("change request for %s", file)
}
