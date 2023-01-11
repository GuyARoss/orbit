// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package hotreload

import (
	"net/http"
	"testing"
	"time"

	"github.com/GuyARoss/orbit/pkg/hotreload/mock"
	"github.com/gorilla/websocket"
)

func TestBundleKeyListDiff(t *testing.T) {
	f := BundleKeyList([]string{"thing", "cat", "dog"})
	ff := BundleKeyList([]string{"thing"})

	p := f.Diff(ff)

	if len(p) != 2 {
		t.Errorf("expected 'cat', 'dog' and got %s", p)
	}
}

func TestHotReloadReloadSingle(t *testing.T) {
	t.Run("active socket", func(t *testing.T) {
		s := &mock.MockSocket{}

		hr := &HotReload{
			currentBundleKeys: []string{},
			socket:            s,
		}

		err := hr.ReloadSignal()
		if err != nil {
			t.Errorf("should not throw err")
		}

		if !s.DidWrite {
			t.Errorf("did not write")
		}
	})
	t.Run("inactive socket", func(t *testing.T) {
		hr := &HotReload{
			currentBundleKeys: []string{},
			socket:            nil,
		}

		err := hr.ReloadSignal()
		if err != nil {
			t.Errorf("should not throw err")
		}
	})
}

func TestHotReloadActiveBundle(t *testing.T) {
	t.Run("active socket", func(t *testing.T) {
		hr := &HotReload{
			currentBundleKeys: []string{
				"thing", "cat",
			},
			socket: &websocket.Conn{},
		}

		tt := []struct {
			i string
			e bool
		}{
			{"cat", true},
			{"no_present", false},
		}

		for i, c := range tt {
			if g := hr.IsActiveBundle(c.i); c.e != g {
				t.Errorf("(%d) expected '%t' got '%t'", i, c.e, g)
			}
		}
	})

	t.Run("inactive socket", func(t *testing.T) {
		hr := &HotReload{
			currentBundleKeys: []string{
				"thing", "cat",
			},
			socket: nil,
		}

		tt := []struct {
			i string
			e bool
		}{
			{"cat", true},
			{"no_present", true},
		}

		for i, c := range tt {
			if g := hr.IsActiveBundle(c.i); c.e != g {
				t.Errorf("(%d) expected '%t' got '%t'", i, c.e, g)
			}
		}
	})
}

func TestEmitLog_SocketNotAvailable(t *testing.T) {
	hr := &HotReload{
		currentBundleKeys: []string{
			"thing", "cat",
		},
		socket: nil,
	}

	resp := hr.EmitLog(Warning, "should not work")
	if resp != nil {
		t.Errorf("socket not available")
	}
}

func TestEmitLog(t *testing.T) {
	ms := &mock.MockSocket{}

	hr := &HotReload{
		currentBundleKeys: []string{
			"thing", "cat",
		},
		socket: ms,
	}

	err := hr.EmitLog(Warning, "warning text")
	if err != nil {
		t.Errorf("error was not expected '%s'", err)
		return
	}

	if !ms.DidWrite {
		t.Error("did not write to socket")
		return
	}
}

func TestCurrentBundles(t *testing.T) {
	bundleKeys := []string{"test", "test2"}

	hr := &HotReload{
		currentBundleKeys: bundleKeys,
	}
	keys := hr.CurrentBundleKeys()
	if len(keys) != len(bundleKeys) {
		t.Errorf("bundle keys not equal got '%s'", keys)
	}
}

func TestNewSocket(t *testing.T) {
	redirectionKeys := []string{"apple", "orange"}

	n := New()
	type localSession struct {
		didRedirect bool
	}

	s := &localSession{
		didRedirect: false,
	}
	go func(r *localSession) {
		resp := <-n.Redirected
		if len(resp.BundleKeys) == len(redirectionKeys) {
			r.didRedirect = true
		}

		r.didRedirect = true
	}(s)

	n.skipUpgrade = true
	n.socket = &mock.MockSocket{
		ReadData: &SocketRequest{
			Operation: "pages",
			Value:     redirectionKeys,
		},
	}
	n.HandleWebSocket(&mock.MockResponseWriter{}, &http.Request{})

	time.Sleep(100 * time.Millisecond) // ensure that the channel has enough time to sync.
	if !s.didRedirect {
		t.Error("expected redirect on socket init\n")
	}

	if len(n.currentBundleKeys) != len(redirectionKeys) {
		t.Errorf("expected updated bundle keys got '%s'", n.currentBundleKeys)
	}
}
