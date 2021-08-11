package cmd

import (
	"fmt"
	"log"

	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/pkg/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type devSession struct {
	pageGenSettings *internal.GenPagesSettings
	sourceMap       map[string]*fs.PackedPage
}

var devCMD = &cobra.Command{
	Use: "dev",
	Run: func(cmd *cobra.Command, args []string) {
		as := &internal.GenPagesSettings{
			PackageName: viper.GetString("pacname"),
			OutDir:      viper.GetString("out"),
			WebDir:      viper.GetString("webdir"),
			BundlerMode: viper.GetString("mode"),
		}

		s, err := createSession(as)

		if err != nil {
			panic(err)
		}

		watch, err := fs.DirectoryWatch(s.pageGenSettings.WebDir)
		if err != nil {
			panic(err)
		}
		done := make(chan bool)
		go func() {
			for {
				select {
				case event := <-watch.FileChange:
					s.executeChangeRequest(fmt.Sprintf("%s/%s", s.pageGenSettings.WebDir, event))
				case err := <-watch.Error:
					log.Fatal(err)
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
