package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/pkg/fs"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type devSession struct {
	pageGenSettings *internal.GenPagesSettings
	sourceMap       map[string]*fs.PackedPage
}

var watcher *fsnotify.Watcher

func createSession(settings *internal.GenPagesSettings) (*devSession, error) {
	err := settings.CleanPathing()
	if err != nil {
		return nil, err
	}

	lib := settings.PackWebDir()

	sourceMap := make(map[string]*fs.PackedPage)
	for _, p := range lib.Pages {
		sourceMap[p.BaseDir] = p
	}

	return &devSession{
		settings, sourceMap,
	}, nil
}

func (s *devSession) executeChangeRequest(file string) {
	source := s.sourceMap[file]
	if source != nil {
		s.pageGenSettings.Repack(source)
	}

	// @@todo: re-enable me
	// s.pageGenSettings.PackWebDir()

	if _, err := os.Stat(".orbit/hotreload"); err != nil {
		syscall.Mkfifo(".orbit/hotreload", 0666)
	}

	f, err := os.OpenFile(".orbit/hotreload", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Errorf("error with hotreload")
	}

	f.WriteString(fmt.Sprintf("cr|%s", source.BaseDir))
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
			fmt.Println("ERROR", err)
		}

		done := make(chan bool)

		go func() {
			for {
				time.Sleep(2 * time.Second)

				select {
				case e := <-watcher.Events:
					if !strings.Contains(e.Name, "node_modules") || !strings.Contains(e.Name, ".orbit") {
						s.executeChangeRequest(e.Name)
					}
				case err := <-watcher.Errors:
					fmt.Println("err", err)
				}
			}
		}()

		<-done
	},
}

func init() {
	RootCMD.AddCommand(devCMD)
}
