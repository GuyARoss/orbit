// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package hotreload

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type HotReloader interface {
	ReloadSignal() error
	HandleWebSocket(w http.ResponseWriter, r *http.Request)
	CurrentBundleKeys() []string
	IsActive() bool
	IsActiveBundle(string) bool
}

type BundleKeyList []string

func (l BundleKeyList) Diff(bundleList BundleKeyList) BundleKeyList {
	changes := make([]string, 0)
	for _, k := range l {
		hasMatch := false
		for _, l := range bundleList {
			if k == l {
				hasMatch = true
				break
			}
		}

		if !hasMatch {
			changes = append(changes, k)
		}
	}

	return changes
}

type RedirectionEvent struct {
	PreviousBundleKeys BundleKeyList
	BundleKeys         BundleKeyList
}

type HotReload struct {
	m        *sync.Mutex
	socket   *websocket.Conn
	upgrader *websocket.Upgrader

	currentBundleKeys BundleKeyList
	Redirected        chan RedirectionEvent
}

type SocketRequest struct {
	Operation string   `json:"operation"`
	Value     []string `json:"value"`
}

func (s *HotReload) ReloadSignal() error {
	if s.IsActive() {
		return s.socket.WriteJSON(&SocketRequest{
			Operation: "reload",
		})
	}

	return nil
}

func (s *HotReload) IsActiveBundle(key string) bool {
	if s.IsActive() {
		for _, k := range s.currentBundleKeys {
			if k == key {
				return true
			}
		}
	} else {
		return true
	}

	return false
}

func (s *HotReload) IsActive() bool {
	return s.socket != nil
}

func (s *HotReload) CurrentBundleKeys() []string {
	return s.currentBundleKeys
}

func (s *HotReload) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	s.m.Lock()

	// close previous socket conn
	if s.socket != nil {
		s.socket.Close()
	}

	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	sockRequest := &SocketRequest{}
	err = c.ReadJSON(sockRequest)

	if err != nil {
		panic(err)
	}

	s.socket = c

	switch sockRequest.Operation {
	case "pages":
		s.Redirected <- RedirectionEvent{
			PreviousBundleKeys: s.currentBundleKeys,
			BundleKeys:         sockRequest.Value,
		}
		s.currentBundleKeys = sockRequest.Value
	}

	s.m.Unlock()
}

func New() *HotReload {
	u := &websocket.Upgrader{}
	u.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	return &HotReload{
		m:                 &sync.Mutex{},
		upgrader:          u,
		Redirected:        make(chan RedirectionEvent),
		currentBundleKeys: make([]string, 0),
	}
}
