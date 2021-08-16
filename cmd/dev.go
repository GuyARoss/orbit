package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	lib := settings.ApplyPages()

	sourceMap := make(map[string]*fs.PackedPage)
	for _, p := range lib.Pages {
		sourceMap[p.BaseDir] = p
	}

	return &devSession{
		settings, sourceMap,
	}, nil
}

func (s *devSession) executeChangeRequest(file string) {
	fmt.Printf("change request for %s", file)
	fmt.Println(s.sourceMap[file])
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

		fmt.Println(s.sourceMap)

		go func() {
			for {
				select {
				// watch for events
				case e := <-watcher.Events:
					if !strings.Contains(e.Name, "node_modules") || !strings.Contains(e.Name, ".orbit") {
						s.executeChangeRequest(e.Name)
					}
					// watch for errors
				case err := <-watcher.Errors:
					fmt.Println("ERROR", err)
				}
			}
		}()

		<-done
	},
}

func init() {
	RootCMD.AddCommand(devCMD)
}
