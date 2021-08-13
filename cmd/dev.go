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

		// starting at the root of the project, walk each file/directory searching for
		// directories
		if err := filepath.Walk("./", watchDir); err != nil {
			fmt.Println("ERROR", err)
		}

		//
		done := make(chan bool)

		//
		go func() {
			for {
				select {
				// watch for events
				case e := <-watcher.Events:
					if !strings.Contains(e.Name, "node_modules") {
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
}

func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

// watcher, err := fsnotify.NewWatcher()
// if err != nil {
// 	panic(err)
// }
// defer watcher.Close()

// done := make(chan bool)
// go func() {
// 	for {
// 		select {
// 		case event := <-watcher.Events:
// 			s.executeChangeRequest(fmt.Sprintf("%s/%s", s.pageGenSettings.WebDir, event.Name))
// 		case err := <-watcher.Errors:
// 			log.Fatal(err)
// 			done <- true
// 		}
// 	}
// }()

// if err := watcher.Add(s.pageGenSettings.WebDir); err != nil {
// 	panic(err)
// }

// <-done
