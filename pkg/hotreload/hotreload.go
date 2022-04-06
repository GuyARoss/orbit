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

type RedirectionEvent struct {
	OldBundleKeys []string
	NewBundleKeys []string
}

type HotReload struct {
	m        *sync.Mutex
	socket   *websocket.Conn
	upgrader *websocket.Upgrader

	currentBundleKeys []string
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
			OldBundleKeys: s.currentBundleKeys,
			NewBundleKeys: sockRequest.Value,
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
