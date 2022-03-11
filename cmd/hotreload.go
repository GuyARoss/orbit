// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.
package cmd

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/fsutils"
	"github.com/GuyARoss/orbit/pkg/log"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SocketRequest struct {
	Operation string `json:"operation"`
	Value     string `json:"value"`
}

type HotReload struct {
	m        *sync.Mutex
	socket   *websocket.Conn
	upgrader *websocket.Upgrader

	CurrentBundleKey string
}

func NewHotReload() *HotReload {
	u := &websocket.Upgrader{}
	u.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	return &HotReload{
		m:        &sync.Mutex{},
		upgrader: u,
	}
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
		s.CurrentBundleKey = sockRequest.Value
	}

	s.m.Unlock()
}

func (s *HotReload) ReloadSignal() error {
	return s.socket.WriteJSON(&SocketRequest{
		Operation: "refresh",
	})
}

var hotreloadCMD = &cobra.Command{
	Use:   "hotreload",
	Long:  "hot-reload bundle data given the specified pages in dev mode",
	Short: "hot-reload bundle data given the specified pages in dev mode",
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.NewDefaultLogger()

		s, err := createSession(context.Background(), &SessionOpts{
			UseDebug: viper.GetBool("usedebug"),
			WebDir:   viper.GetString("webdir"),
		})

		if err != nil {
			panic(err)
		}

		watcher, _ := fsnotify.NewWatcher()
		defer watcher.Close()

		hotReload := NewHotReload()

		if err := filepath.Walk(fsutils.NormalizePath("./"), watchDir(watcher)); err != nil {
			panic("invalid walk on watchDir")
		}

		go func(hr *HotReload) {
			sh := srcpack.NewSyncHook(log.NewDefaultLogger())

			for {
				time.Sleep(time.Duration(viper.GetInt("timeout")) * time.Millisecond)

				select {
				case e := <-watcher.Events:
					root := s.rootComponents[e.Name]

					// page is the current bundle that is open in the browser
					// process change, recompute bundle and send refresh signal back to browser
					if root != nil && s.rootComponents[e.Name].BundleKey == hr.CurrentBundleKey {
						s.directFileChangeRequest(e.Name, time.Duration(viper.GetInt("samefiletimeout"))*time.Millisecond, sh)

						err := hr.ReloadSignal()
						if err != nil {
							fmt.Println(err)
						}
					}

					// component may exist as a page depencency, if so, recompute and send refresh signal
					if len(s.sourceMap.FindRoot(e.Name)) > 0 {
						s.indirectFileChangeRequest(e.Name, hr.CurrentBundleKey, time.Duration(viper.GetInt("samefiletimeout"))*time.Millisecond, sh)

						err := hr.ReloadSignal()
						if err != nil {
							fmt.Println(err)
						}
					}
				case err := <-watcher.Errors:
					panic(fmt.Sprintf("watcher failed %s", err.Error()))
				}
			}
		}(hotReload)

		http.HandleFunc("/ws", hotReload.HandleWebSocket)
		logger.Info(fmt.Sprintf("server started on port %d", viper.GetInt("hotreloadport")))

		err = http.ListenAndServe(fmt.Sprintf("localhost:%d", viper.GetInt("hotreloadport")), nil)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	var timeoutDuration int
	var samefileTimeout int
	var port int

	hotreloadCMD.PersistentFlags().IntVar(&timeoutDuration, "timeout", 2000, "specifies the timeout duration in milliseconds until a change will be detected")
	viper.BindPFlag("timeout", hotreloadCMD.PersistentFlags().Lookup("timeout"))

	hotreloadCMD.PersistentFlags().IntVar(&samefileTimeout, "samefiletimeout", 2000, "specifies the timeout duration in milliseconds until a change will be detected for repeating files")
	viper.BindPFlag("samefiletimeout", hotreloadCMD.PersistentFlags().Lookup("samefiletimeout"))

	hotreloadCMD.PersistentFlags().IntVar(&port, "hotreloadport", 3005, "port used for hotreload")
	viper.BindPFlag("hotreloadport", hotreloadCMD.PersistentFlags().Lookup("hotreloadport"))

	RootCMD.AddCommand(hotreloadCMD)
}
