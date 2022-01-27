package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/GuyARoss/orbit/internal"
	dependtree "github.com/GuyARoss/orbit/pkg/depend_tree"
	"github.com/GuyARoss/orbit/pkg/log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type proccessedChangeRequest struct {
	ProcessedAt time.Time
	FileName    string
}

type devSession struct {
	pageGenSettings *internal.GenPagesSettings
	rootComponents  map[string]*internal.PackedComponent
	sourceMap       *dependtree.DependencySourceMap

	lastProcessedFile *proccessedChangeRequest
	m                 *sync.Mutex
}

var watcher *fsnotify.Watcher

func createSession(settings *internal.GenPagesSettings) (*devSession, error) {
	err := settings.CleanPathing()
	if err != nil {
		return nil, err
	}

	lib, err := settings.PackWebDir(nil)
	if err != nil {
		return nil, err
	}

	rootComponents := make(map[string]*internal.PackedComponent)
	for _, p := range lib.Pages {
		rootComponents[p.OriginalFilePath] = p
	}

	sourceMap, err := internal.CreateSourceMap(settings.WebDir, lib.Pages, settings.WebDir)
	if err != nil {
		return nil, err
	}

	return &devSession{
		pageGenSettings:   settings,
		rootComponents:    rootComponents,
		sourceMap:         sourceMap,
		lastProcessedFile: &proccessedChangeRequest{},
		m:                 &sync.Mutex{},
	}, nil
}

func (s *devSession) executeChangeRequest(file string, timeoutDuration time.Duration) error {
	if file == s.lastProcessedFile.FileName &&
		time.Since(s.lastProcessedFile.ProcessedAt).Seconds() < timeoutDuration.Seconds() {
		return nil
	}

	component := s.rootComponents[fmt.Sprintf("./%s", file)]
	if component != nil {
		s.m.Lock()
		s.lastProcessedFile = &proccessedChangeRequest{
			FileName:    file,
			ProcessedAt: time.Now(),
		}
		s.m.Unlock()

		s.pageGenSettings.Repack(component)
	}

	sources := s.sourceMap.FindRoot(file)
	for _, source := range sources {
		component = s.rootComponents[source]

		if component != nil {
			s.m.Lock()
			s.lastProcessedFile = &proccessedChangeRequest{
				FileName:    file,
				ProcessedAt: time.Now(),
			}
			s.m.Unlock()

			s.pageGenSettings.Repack(component)
		}
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
	Use:   "dev",
	Long:  "hot-reload bundle data given the specified pages in dev mode",
	Short: "hot-reload bundle data given the specified pages in dev mode",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("starting dev server...")
		as := &internal.GenPagesSettings{
			PackageName:    viper.GetString("pacname"),
			OutDir:         viper.GetString("out"),
			WebDir:         viper.GetString("webdir"),
			BundlerMode:    viper.GetString("mode"),
			NodeModulePath: viper.GetString("nodemod"),
			PublicDir:      viper.GetString("publicdir"),
		}

		s, err := createSession(as)
		if err != nil {
			panic(err)
		}

		log.Success("dev server started successfully")

		watcher, _ = fsnotify.NewWatcher()
		defer watcher.Close()

		if err := filepath.Walk("./", watchDir); err != nil {
			log.Error("invalid walk on watchDir")
			return
		}

		done := make(chan bool)

		go func() {
			for {
				time.Sleep(time.Duration(viper.GetInt("timeout")) * time.Second)

				select {
				case e := <-watcher.Events:
					if !strings.Contains(e.Name, "node_modules") || !strings.Contains(e.Name, ".orbit") {
						go s.executeChangeRequest(e.Name, time.Duration(viper.GetInt("samefiletimeout"))*time.Second)
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
	var timeoutDuration int
	var samefileTimeout int

	devCMD.PersistentFlags().IntVar(&timeoutDuration, "timeout", 2, "specifies the timeout duration in seconds until a change will be detected")
	viper.BindPFlag("timeout", devCMD.PersistentFlags().Lookup("timeout"))

	devCMD.PersistentFlags().IntVar(&samefileTimeout, "samefiletimeout", 5, "specifies the timeout duration in seconds until a change will be detected for repeating files")
	viper.BindPFlag("samefiletimeout", devCMD.PersistentFlags().Lookup("samefiletimeout"))

	RootCMD.AddCommand(devCMD)
}
