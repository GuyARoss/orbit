package dev

import (
	"fmt"
	"log"

	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/pkg/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CMD = &cobra.Command{
	Use: "dev",
	Run: func(cmd *cobra.Command, args []string) {
		as := &internal.GenPagesSettings{
			PackageName: viper.GetString("pacname"),
			OutDir:      viper.GetString("out"),
			WebDir:      viper.GetString("webdir"),
			BundlerMode: viper.GetString("mode"),
		}

		session, err := createSession(as)

		if err != nil {
			panic(err)
		}

		execute(session)
	},
}

func execute(s *session) {
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
