// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package internal

import (
	"errors"
	"testing"
	"time"

	liboutmock "github.com/GuyARoss/orbit/internal/libout/mock"
	srcpackmock "github.com/GuyARoss/orbit/internal/srcpack/mock"
	hotreloadmock "github.com/GuyARoss/orbit/pkg/hotreload/mock"

	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/log"
)

func TestProcessChangeRequest_TooRecentlyProcessed(t *testing.T) {
	fn := "this_was_recently_processed.txt"

	s := devSession{
		lastProcessedFile: &proccessedChangeRequest{
			FileName:    fn,
			ProcessedAt: time.Now(),
		},
		SessionOpts: &SessionOpts{},
	}

	err := s.DoChangeRequest(fn, &ChangeRequestOpts{
		SafeFileTimeout: time.Second * 50,
	})

	if err == nil || !errors.Is(err, ErrFileTooRecentlyProcessed) {
		t.Errorf("expected err file too recently processed")
	}
}

func TestDoChangeRequest_DirectFile(t *testing.T) {
	fn := "direct_file_thing"
	comp := &srcpackmock.MockPackedComponent{}
	s := devSession{
		lastProcessedFile: &proccessedChangeRequest{},
		SessionOpts:       &SessionOpts{},
		RootComponents: map[string]srcpack.PackComponent{
			fn: comp,
		},
	}

	hotReloader := &hotreloadmock.MockHotReload{}
	err := s.DoChangeRequest(fn, &ChangeRequestOpts{
		SafeFileTimeout: time.Hour * 2,
		HotReload:       hotReloader,
		Hook:            srcpack.NewSyncHook(log.NewEmptyLogger()),
	})

	if err != nil {
		t.Errorf("error should not have been thrown during direct file processing")
	}

	if hotReloader.DidReload == false {
		t.Errorf("hot reloading did not occur after file processing")
	}

	if comp.WasRepacked == false {
		t.Errorf("packing did not occur during direct file processing")
	}
}

func TestDoChangeRequest_IndirectFile(t *testing.T) {
	fn := "direct_file_thing"

	comp := &srcpackmock.MockPackedComponent{}

	s := devSession{
		lastProcessedFile: &proccessedChangeRequest{},
		SessionOpts:       &SessionOpts{},
		RootComponents: map[string]srcpack.PackComponent{
			"thing2": comp,
		},
		SourceMap: map[string][]string{
			fn: {"thing2"},
		},
	}
	hotReloader := &hotreloadmock.MockHotReload{}
	err := s.DoChangeRequest(fn, &ChangeRequestOpts{
		SafeFileTimeout: time.Hour * 2,
		HotReload:       hotReloader,
		Hook:            srcpack.NewSyncHook(log.NewEmptyLogger()),
	})
	if err != nil {
		t.Errorf("error should not have been thrown during indirect file processing")
	}

	if hotReloader.DidReload == false {
		t.Errorf("hot reloading did not occur after file processing")
	}
	if comp.WasRepacked == false {
		t.Errorf("packing did not occur during direct file processing")
	}
}

type mockPacker struct {
	components []srcpack.Component
}

func (m *mockPacker) PackMany(pages []string) ([]srcpack.PackComponent, error) { return nil, nil }
func (m *mockPacker) PackSingle(logger log.Logger, file string) (srcpack.PackComponent, error) {
	return &m.components[0], nil
}
func (m *mockPacker) ReattachLogger(logger log.Logger) srcpack.Packer { return nil }

func TestDoChangeRequest_UnknownPage(t *testing.T) {
	fn := "/pages/filename.jsx"

	s := devSession{
		lastProcessedFile: &proccessedChangeRequest{},
		SessionOpts:       &SessionOpts{},
		RootComponents:    map[string]srcpack.PackComponent{},
		SourceMap:         map[string][]string{},
		packer: &mockPacker{
			components: []srcpack.Component{
				{},
			},
		},
		libout: &liboutmock.MockBundleWriter{},
	}
	hotReloader := &hotreloadmock.MockHotReload{}
	err := s.DoChangeRequest(fn, &ChangeRequestOpts{
		SafeFileTimeout: time.Hour * 2,
		HotReload:       hotReloader,
		Hook:            srcpack.NewSyncHook(log.NewEmptyLogger()),
	})

	if err != nil {
		t.Errorf("error should not have been thrown during processing of an unknown page")
	}

	if len(s.RootComponents) != 1 {
		t.Errorf("page was not correctly identified")
	}
}
