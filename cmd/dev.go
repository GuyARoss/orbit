// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.
package cmd

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/fsutils"
	"github.com/GuyARoss/orbit/pkg/hotreload"
	"github.com/GuyARoss/orbit/pkg/log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var devCMD = &cobra.Command{
	Use:   "dev",
	Long:  "hot-reload bundle data given the specified pages in dev mode",
	Short: "hot-reload bundle data given the specified pages in dev mode",
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.NewDefaultLogger()

		s, err := internal.CreateSession(context.Background(), &internal.SessionOpts{
			UseDebug:      viper.GetBool("usedebug"),
			WebDir:        viper.GetString("webdir"),
			Mode:          viper.GetString("mode"),
			Pacname:       viper.GetString("pacname"),
			OutDir:        viper.GetString("out"),
			NodeModDir:    viper.GetString("nodemod"),
			PublicDir:     viper.GetString("publicdir"),
			HotReloadPort: viper.GetInt("hotreloadport"),
		})

		if err != nil {
			panic(err)
		}

		watcher, _ := fsnotify.NewWatcher()
		defer watcher.Close()

		hotReload := hotreload.New()

		if err := filepath.Walk(fsutils.NormalizePath("./"), internal.WatchDir(watcher)); err != nil {
			panic("invalid walk on watchDir")
		}

		go func(hr *hotreload.HotReload) {
			sh := srcpack.NewSyncHook(log.NewDefaultLogger())

			for {
				time.Sleep(time.Duration(viper.GetInt("timeout")) * time.Millisecond)

				select {
				case e := <-watcher.Events:
					root := s.RootComponents[e.Name]

					// page is the current bundle that is open in the browser
					// process change, recompute bundle and send refresh signal back to browser
					if root != nil && s.RootComponents[e.Name].BundleKey == hr.CurrentBundleKey {
						s.DirectFileChangeRequest(e.Name, time.Duration(viper.GetInt("samefiletimeout"))*time.Millisecond, sh)

						err := hr.ReloadSignal()
						if err != nil {
							fmt.Println(err)
						}
					}

					// component may exist as a page depencency, if so, recompute and send refresh signal
					if len(s.SourceMap.FindRoot(e.Name)) > 0 {
						s.IndirectFileChangeRequest(e.Name, hr.CurrentBundleKey, time.Duration(viper.GetInt("samefiletimeout"))*time.Millisecond, sh)

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

	devCMD.PersistentFlags().IntVar(&timeoutDuration, "timeout", 2000, "specifies the timeout duration in milliseconds until a change will be detected")
	viper.BindPFlag("timeout", devCMD.PersistentFlags().Lookup("timeout"))

	devCMD.PersistentFlags().IntVar(&samefileTimeout, "samefiletimeout", 2000, "specifies the timeout duration in milliseconds until a change will be detected for repeating files")
	viper.BindPFlag("samefiletimeout", devCMD.PersistentFlags().Lookup("samefiletimeout"))

	devCMD.PersistentFlags().IntVar(&port, "hotreloadport", 3005, "port used for hotreload")
	viper.BindPFlag("hotreloadport", devCMD.PersistentFlags().Lookup("hotreloadport"))

	RootCMD.AddCommand(devCMD)
}
