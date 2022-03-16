package internal

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/GuyARoss/orbit/internal/libout"
	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/log"
	webwrapper "github.com/GuyARoss/orbit/pkg/web_wrapper"
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

type mockHotReload struct {
	didReload        bool
	currentBundleKey string
	reloadErr        error
}

func (m *mockHotReload) ReloadSignal() error {
	m.didReload = true

	return m.reloadErr
}

func (m *mockHotReload) HandleWebSocket(w http.ResponseWriter, r *http.Request) {}
func (m *mockHotReload) CurrentBundleKey() string {
	return m.currentBundleKey
}

type mockPackedComponent struct {
	wasRepacked bool
}

func (m *mockPackedComponent) Repack() error {
	m.wasRepacked = true
	return nil
}

func (m *mockPackedComponent) RepackForWaitGroup(wg *sync.WaitGroup, c chan error) {}
func (m *mockPackedComponent) OriginalFilePath() string                            { return "" }
func (m *mockPackedComponent) Dependencies() []*jsparse.ImportDependency           { return nil }
func (m *mockPackedComponent) BundleKey() string                                   { return "" }
func (m *mockPackedComponent) Name() string                                        { return "" }
func (m *mockPackedComponent) WebWrapper() webwrapper.JSWebWrapper                 { return nil }

func TestDoChangeRequest_DirectFile(t *testing.T) {
	fn := "direct_file_thing"
	comp := &mockPackedComponent{}
	s := devSession{
		lastProcessedFile: &proccessedChangeRequest{},
		SessionOpts:       &SessionOpts{},
		RootComponents: map[string]srcpack.PackComponent{
			fn: comp,
		},
	}

	hotReloader := &mockHotReload{}
	err := s.DoChangeRequest(fn, &ChangeRequestOpts{
		SafeFileTimeout: time.Hour * 2,
		HotReload:       hotReloader,
		Hook:            srcpack.NewSyncHook(log.NewEmptyLogger()),
	})

	if err != nil {
		t.Errorf("error should not have been thrown during direct file processing")
	}

	if hotReloader.didReload == false {
		t.Errorf("hot reloading did not occur after file processing")
	}

	if comp.wasRepacked == false {
		t.Errorf("packing did not occur during direct file processing")
	}
}

func TestDoChangeRequest_IndirectFile(t *testing.T) {
	fn := "direct_file_thing"

	comp := &mockPackedComponent{}

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
	hotReloader := &mockHotReload{}
	err := s.DoChangeRequest(fn, &ChangeRequestOpts{
		SafeFileTimeout: time.Hour * 2,
		HotReload:       hotReloader,
		Hook:            srcpack.NewSyncHook(log.NewEmptyLogger()),
	})
	if err != nil {
		t.Errorf("error should not have been thrown during indirect file processing")
	}

	if hotReloader.didReload == false {
		t.Errorf("hot reloading did not occur after file processing")
	}
	if comp.wasRepacked == false {
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

type mockBundleWriter struct{}

func (m *mockBundleWriter) WriteLibout(files libout.Libout, fOpts *libout.FilePathOpts) error {
	return nil
}
func (m *mockBundleWriter) AcceptComponent(ctx context.Context, c srcpack.PackComponent, cacheOpts *webwrapper.CacheDOMOpts) {
}
func (m *mockBundleWriter) AcceptComponents(ctx context.Context, comps []srcpack.PackComponent, cacheOpts *webwrapper.CacheDOMOpts) {
}

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
		libout: &mockBundleWriter{},
	}
	hotReloader := &mockHotReload{}
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
