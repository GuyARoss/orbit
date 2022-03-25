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
	CurrentBundleKey() string
	IsActive() bool
	IsActiveBundle(string) bool
}

type HotReload struct {
	m        *sync.Mutex
	socket   *websocket.Conn
	upgrader *websocket.Upgrader

	currentBundleKey string
}

type SocketRequest struct {
	Operation string `json:"operation"`
	Value     string `json:"value"`
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
		return s.currentBundleKey == key
	}

	return true
}

func (s *HotReload) IsActive() bool {
	return s.socket != nil
}

func (s *HotReload) CurrentBundleKey() string {
	return s.currentBundleKey
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
	case "page":
		s.currentBundleKey = sockRequest.Value
	}

	s.m.Unlock()
}

func New() *HotReload {
	u := &websocket.Upgrader{}
	u.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	return &HotReload{
		m:        &sync.Mutex{},
		upgrader: u,
	}
}
