// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package internal

import (
	"context"
	"fmt"
	"testing"
	"time"

	liboutmock "github.com/GuyARoss/orbit/internal/libout/mock"
	srcpackmock "github.com/GuyARoss/orbit/internal/srcpack/mock"
	allocatedstack "github.com/GuyARoss/orbit/pkg/allocated_stack"
	hotreloadmock "github.com/GuyARoss/orbit/pkg/hotreload/mock"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/jsparse/mock"

	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/log"
)

func TestProcessChangeRequest_TooRecentlyProcessed(t *testing.T) {
	fn := "this_was_recently_processed.txt"

	s := devSession{
		ChangeRequest: &changeRequest{
			LastProcessedAt: time.Now(),
			LastFileName:    fn,
			changeRequests:  allocatedstack.New(1),
		},
		SessionOpts: &SessionOpts{},
	}

	err := s.DoFileChangeRequest(fn, &ChangeRequestOpts{
		SafeFileTimeout: time.Second * 50,
	})

	if err == nil {
		fmt.Println("err name", err)
		t.Errorf("expected err file too recently processed")
	}
}

func TestDoChangeRequest_DirectFile(t *testing.T) {
	fn := "direct_file_thing"
	comp := &srcpackmock.MockPackedComponent{
		FilePath: "./test/",
		Depends: []*jsparse.ImportDependency{
			{
				FinalStatement: `import React from '../react.js'`,
				InitialPath:    "./test/react.js",
				Type:           jsparse.LocalImportType,
			},
		},
	}
	s := devSession{
		ChangeRequest: &changeRequest{
			LastProcessedAt: time.Now(),
			LastFileName:    "",
			changeRequests:  allocatedstack.New(1),
		},
		SessionOpts: &SessionOpts{},
		RootComponents: map[string]srcpack.PackComponent{
			fn: comp,
		},
		SourceMap: map[string][]string{},
	}

	hotReloader := &hotreloadmock.MockHotReload{
		Active: true,
	}
	err := s.DoFileChangeRequest(fn, &ChangeRequestOpts{
		SafeFileTimeout: time.Hour * 2,
		HotReload:       hotReloader,
		Hook:            srcpack.NewSyncHook(log.NewEmptyLogger()),
		Parser: &mock.MockJSParser{
			ParseDocument: jsparse.NewEmptyDocument(),
			Err:           nil,
		},
	})

	if err != nil {
		t.Errorf("error should not have been thrown during direct file processing '%s'", err)
		return
	}

	if hotReloader.DidReload == false {
		t.Errorf("hot reloading did not occur after file processing")
	}

	if comp.WasRepacked == false {
		t.Errorf("packing did not occur during direct file processing")
		return
	}

	if len(s.SourceMap["./test/react.js"]) != 1 {
		t.Errorf("did not merge dependent trees")
		return
	}
}

func TestDoChangeRequest_IndirectFile(t *testing.T) {
	fn := "direct_file_thing"

	comp := &srcpackmock.MockPackedComponent{
		FilePath: "./test/",
		Depends: []*jsparse.ImportDependency{
			{
				FinalStatement: `import React from '../react.js'`,
				InitialPath:    "./test/react.js",
				Type:           jsparse.LocalImportType,
			},
		},
	}

	s := devSession{
		ChangeRequest: &changeRequest{
			LastProcessedAt: time.Now(),
			LastFileName:    "",
			changeRequests:  allocatedstack.New(1),
		},
		SessionOpts: &SessionOpts{},
		RootComponents: map[string]srcpack.PackComponent{
			"thing2": comp,
		},
		SourceMap: map[string][]string{
			fn: {"thing2"},
		},
	}

	hotReloader := &hotreloadmock.MockHotReload{
		Active: true,
	}

	err := s.DoFileChangeRequest(fn, &ChangeRequestOpts{
		SafeFileTimeout: time.Hour * 2,
		HotReload:       hotReloader,
		Hook:            srcpack.NewSyncHook(log.NewEmptyLogger()),
		Parser: &mock.MockJSParser{
			ParseDocument: jsparse.NewDocument("./test", "react.jsx"),
			Err:           nil,
		},
	})

	if err != nil {
		t.Errorf("error should not have been thrown during indirect file processing '%s'", err)
	}

	if hotReloader.DidReload == false {
		t.Errorf("hot reloading did not occur after file processing")
	}

	if comp.WasRepacked == false {
		t.Errorf("packing did not occur during indirect file processing")
	}

	if len(s.SourceMap["./test/react.js"]) != 1 {
		t.Errorf("did not merge dependent trees")
	}
}

type mockPacker struct {
	components []srcpack.Component
	failPack   bool
}

func (m *mockPacker) PackMany(pages []string) (srcpack.PackedComponentList, error) {
	if m.failPack {
		return nil, fmt.Errorf("error")
	}
	return nil, nil
}
func (m *mockPacker) PackSingle(logger log.Logger, file string) (srcpack.PackComponent, error) {
	if m.failPack {
		return nil, fmt.Errorf("error")
	}
	return &m.components[0], nil
}
func (m *mockPacker) ReattachLogger(logger log.Logger) srcpack.Packer { return nil }

func TestDoChangeRequest_UnknownPage(t *testing.T) {
	fn := "/pages/filename.jsx"

	s := devSession{
		ChangeRequest: &changeRequest{
			LastProcessedAt: time.Now(),
			LastFileName:    "",
			changeRequests:  allocatedstack.New(1),
		},
		SessionOpts:    &SessionOpts{},
		RootComponents: map[string]srcpack.PackComponent{},
		SourceMap:      map[string][]string{},
		packer: &mockPacker{
			components: []srcpack.Component{
				{},
			},
		},
		libout: &liboutmock.MockBundleWriter{},
	}
	hotReloader := &hotreloadmock.MockHotReload{}
	err := s.DoFileChangeRequest(fn, &ChangeRequestOpts{
		SafeFileTimeout: time.Hour * 2,
		HotReload:       hotReloader,
		Hook:            srcpack.NewSyncHook(log.NewEmptyLogger()),
		Parser: &mock.MockJSParser{
			ParseDocument: jsparse.NewEmptyDocument(),
			Err:           nil,
		},
	})

	if err != nil {
		t.Errorf("error should not have been thrown during processing of an unknown page")
	}

	if len(s.RootComponents) != 1 {
		t.Errorf("page was not correctly identified")
	}
}

func TestDoBundleChangeRequest(t *testing.T) {
	bundle := "test_bundle"
	comp := &srcpackmock.MockPackedComponent{
		FilePath: "./test/",
		Key:      bundle,
		Depends: []*jsparse.ImportDependency{
			{
				FinalStatement: `import React from '../react.js'`,
				InitialPath:    "./test/react.js",
				Type:           jsparse.LocalImportType,
			},
		},
	}

	s := devSession{
		ChangeRequest: &changeRequest{
			LastProcessedAt: time.Now(),
			LastFileName:    "",
			changeRequests:  allocatedstack.New(1),
		},
		SessionOpts: &SessionOpts{},
		RootComponents: map[string]srcpack.PackComponent{
			"test": comp,
		},
		SourceMap: map[string][]string{},
		packer: &mockPacker{
			components: []srcpack.Component{
				{},
			},
		},
		libout: &liboutmock.MockBundleWriter{},
	}
	hotReloader := &hotreloadmock.MockHotReload{}
	err := s.DoBundleKeyChangeRequest(bundle, &ChangeRequestOpts{
		SafeFileTimeout: time.Hour * 2,
		HotReload:       hotReloader,
		Hook:            srcpack.NewSyncHook(log.NewEmptyLogger()),
		Parser: &mock.MockJSParser{
			ParseDocument: jsparse.NewEmptyDocument(),
			Err:           nil,
		},
	})

	if err != nil {
		t.Errorf("error should not have been thrown during processing of an unknown page")
	}

	if len(s.RootComponents) != 1 {
		t.Errorf("page was not correctly identified")
	}
}

func TestUnknownPageError(t *testing.T) {
	fn := "/pages/filename.jsx"

	s := devSession{
		ChangeRequest: &changeRequest{
			LastProcessedAt: time.Now(),
			LastFileName:    "",
			changeRequests:  allocatedstack.New(1),
		},
		SessionOpts:    &SessionOpts{},
		RootComponents: map[string]srcpack.PackComponent{},
		SourceMap:      map[string][]string{},
		packer: &mockPacker{
			failPack: true,
			components: []srcpack.Component{
				{},
			},
		},
		libout: &liboutmock.MockBundleWriter{},
	}
	hotReloader := &hotreloadmock.MockHotReload{}
	err := s.DoFileChangeRequest(fn, &ChangeRequestOpts{
		SafeFileTimeout: time.Hour * 2,
		HotReload:       hotReloader,
		Hook:            srcpack.NewSyncHook(log.NewEmptyLogger()),
		Parser: &mock.MockJSParser{
			ParseDocument: jsparse.NewEmptyDocument(),
			Err:           nil,
		},
	})

	if err == nil {
		t.Errorf("error should have been thrown during invalid pack")
	}
}

func TestNew(t *testing.T) {
	ops := &SessionOpts{
		WebDir:     "",
		Mode:       "development",
		Pacname:    "test",
		OutDir:     t.TempDir() + "/out",
		NodeModDir: t.TempDir() + "/node_module",
		PublicDir:  t.TempDir() + "/publicdir",
	}

	_, err := New(context.TODO(), ops)
	if err == nil {
		t.Errorf("")
	}
}
