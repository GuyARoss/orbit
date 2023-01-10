// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package hotreload

import (
	"testing"

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

func TestHotReloadReloadSingal(t *testing.T) {
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
