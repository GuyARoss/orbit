// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/experiments"
	"github.com/GuyARoss/orbit/pkg/hotreload"
	"github.com/GuyARoss/orbit/pkg/jsparse"
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

		err := experiments.Load(logger, viper.GetStringSlice("experimental"))
		if err != nil {
			logger.Warn(err.Error())
		}

		s, err := internal.New(context.Background(), &internal.SessionOpts{
			BuildOpts:     internal.NewBuildOptsFromViper(),
			HotReloadPort: viper.GetInt("hotreloadport"),
		})

		if err != nil {
			logger.Warn(err.Error())
			return
		}

		watcher, _ := fsnotify.NewWatcher()
		defer watcher.Close()

		reloader := hotreload.New()

		if err := filepath.Walk(viper.GetString("webdir"), WatchDir(watcher)); err != nil {
			panic("invalid walk on watchDir")
		}

		timeout := time.Duration(viper.GetInt("timeout")) * time.Millisecond

		fileChangeOpts := &internal.ChangeRequestOpts{
			SafeFileTimeout: time.Duration(viper.GetInt("samefiletimeout")) * time.Millisecond,
			Hook:            srcpack.NewSyncHook(log.NewEmptyLogger()),
			HotReload:       reloader,
			Parser:          &jsparse.JSFileParser{},
		}

		if viper.GetBool("terminate_on_startup") {
			return
		}

		devServer := internal.NewDevServer(reloader, logger, s, fileChangeOpts)

		go devServer.FileWatcherBundler(timeout, watcher)
		go devServer.RedirectionBundler()

		http.HandleFunc("/ws", reloader.HandleWebSocket)

		logger.Info(fmt.Sprintf("Hot reload server started on port '%d'", viper.GetInt("hotreloadport")))
		logger.Info("You will still need to run your application")

		err = http.ListenAndServe(fmt.Sprintf("localhost:%d", viper.GetInt("hotreloadport")), nil)
		if err != nil {
			panic(err)
		}
	},
}

// watchDir is a utility function used by the file path walker that applies
// each sub directory found under a path to the file watcher
func WatchDir(watcher *fsnotify.Watcher) func(path string, fi os.FileInfo, err error) error {
	return func(path string, fi os.FileInfo, err error) error {
		if strings.Contains(path, "node_modules") {
			return nil
		}

		if fi.Mode().IsDir() {
			return watcher.Add(path)
		}

		return nil
	}
}

func init() {
	var timeoutDuration int
	var samefileTimeout int
	var port int
	var terminateStartup bool

	devCMD.PersistentFlags().IntVar(&timeoutDuration, "timeout", 500, "specifies the timeout duration in milliseconds until a change will be detected")
	viper.BindPFlag("timeout", devCMD.PersistentFlags().Lookup("timeout"))

	devCMD.PersistentFlags().IntVar(&samefileTimeout, "same_file_timeout", 500, "specifies the timeout duration in milliseconds until a change will be detected for repeating files")
	viper.BindPFlag("samefiletimeout", devCMD.PersistentFlags().Lookup("same_file_timeout"))

	devCMD.PersistentFlags().IntVar(&port, "hot_reload_port", 3005, "port used for hotreload")
	viper.BindPFlag("hotreloadport", devCMD.PersistentFlags().Lookup("hot_reload_port"))

	devCMD.PersistentFlags().BoolVar(&terminateStartup, "terminate_on_startup", false, "flag used for terminating the dev command after startup")
	viper.BindPFlag("terminate_on_startup", devCMD.PersistentFlags().Lookup("terminate_on_startup"))
}
