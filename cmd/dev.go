package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GuyARoss/orbit/internal"
	dependtree "github.com/GuyARoss/orbit/pkg/depend_tree"
	"github.com/GuyARoss/orbit/pkg/log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type devSession struct {
	pageGenSettings *internal.GenPagesSettings
	rootComponents  map[string]*internal.PackedComponent
	sourceMap       *dependtree.DependencySourceMap
}

var watcher *fsnotify.Watcher

func createSession(settings *internal.GenPagesSettings) (*devSession, error) {
	err := settings.CleanPathing()
	if err != nil {
		return nil, err
	}

	lib := settings.PackWebDir()

	rootComponents := make(map[string]*internal.PackedComponent)
	for _, p := range lib.Pages {
		rootComponents[p.OriginalFilePath] = p
	}

	sourceMap, err := internal.CreateSourceMap(settings.WebDir, lib.Pages)
	if err != nil {
		return nil, err
	}

	return &devSession{
		pageGenSettings: settings,
		rootComponents:  rootComponents,
		sourceMap:       sourceMap,
	}, nil
}

func (s *devSession) executeChangeRequest(file string) error {
	applied := false

	sources := s.sourceMap.FindRoot(file)
	for _, source := range sources {
		component := s.rootComponents[source]

		if component != nil {
			applied = true
			s.pageGenSettings.Repack(component)
		}
	}

	if !applied {
		pages := s.pageGenSettings.PackWebDir()
		writeErr := pages.WriteOut()
		return writeErr

	}

	return nil
}

func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

var devCMD = &cobra.Command{
	Use: "dev",
	Run: func(cmd *cobra.Command, args []string) {
		as := &internal.GenPagesSettings{
			PackageName:    viper.GetString("pacname"),
			OutDir:         viper.GetString("out"),
			WebDir:         viper.GetString("webdir"),
			BundlerMode:    viper.GetString("mode"),
			AssetDir:       viper.GetString("assetdir"),
			NodeModulePath: viper.GetString("nodemod"),
		}

		s, err := createSession(as)

		if err != nil {
			panic(err)
		}

		watcher, _ = fsnotify.NewWatcher()
		defer watcher.Close()

		if err := filepath.Walk("./", watchDir); err != nil {
			log.Error("invalid walk on watchDir")
			return
		}

		done := make(chan bool)

		go func() {
			for {
				time.Sleep(2 * time.Second)

				select {
				case e := <-watcher.Events:
					if !strings.Contains(e.Name, "node_modules") || !strings.Contains(e.Name, ".orbit") {
						err := s.executeChangeRequest(e.Name)
						if err != nil {
							log.Error(err.Error())
							os.Exit(1)
						}
					}
				case err := <-watcher.Errors:
					log.Error(fmt.Sprintf("watcher failed %s", err.Error()))
				}
			}
		}()

		<-done
	},
}

func init() {
	RootCMD.AddCommand(devCMD)
}
