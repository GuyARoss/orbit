package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/GuyARoss/orbit/pkg/hotreload"
	"github.com/GuyARoss/orbit/pkg/log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type DevServer struct {
	hr             *hotreload.HotReload
	logger         log.Logger
	session        *devSession
	fileChangeOpts *ChangeRequestOpts
}

// RedirectionBundler waits for a redirection event from the client and performs a re-bundle if needed.
func (s *DevServer) RedirectionBundler() {
	for {
		// during dev mode when the browser redirects, we want to process
		// the file only if the bundle has not already been processed
		event := <-s.hr.Redirected

		for _, bundleKey := range event.BundleKeys.Diff(event.PreviousBundleKeys) {
			// the change request maintains a cache of recently bundled pages
			// if it exists on the cache, then we don't care to process it
			if !s.session.ChangeRequest.ExistsInCache(bundleKey) {
				go func(change string) {
					err := s.session.DoBundleKeyChangeRequest(change, s.fileChangeOpts)

					if err != nil {
						s.hr.EmitLog(hotreload.Warning, err.Error())
					}
				}(bundleKey)
			}
		}
	}
}

var blacklistedDirectories = []string{
	".orbit/",
}

func isBlacklistedDirectory(dir string) bool {
	for _, b := range blacklistedDirectories {
		if strings.Contains(dir, b) {
			return true
		}
	}
	return false
}

// FileWatcherBundler watches for events given the file watcher and processes change requests as found
func (s *DevServer) FileWatcherBundler(timeout time.Duration, watcher *fsnotify.Watcher) {
	for {
		time.Sleep(timeout)

		select {
		case e := <-watcher.Events:
			if isBlacklistedDirectory(e.Name) {
				continue
			}

			err := s.session.DoFileChangeRequest(e.Name, s.fileChangeOpts)

			switch err {
			case nil, ErrFileTooRecentlyProcessed:
				//
			default:
				s.hr.EmitLog(hotreload.Error, err.Error())
				s.logger.Error(err.Error())
			}

			if err == nil && len(viper.GetString("depout")) > 0 {
				s.session.SourceMap.Write(viper.GetString("depout"))
			}
		case err := <-watcher.Errors:
			panic(fmt.Sprintf("watcher failed %s", err.Error()))
		}
	}
}

func NewDevServer(hotReload *hotreload.HotReload, logger log.Logger, session *devSession, changeOpts *ChangeRequestOpts) *DevServer {
	return &DevServer{
		hr:             hotReload,
		logger:         logger,
		session:        session,
		fileChangeOpts: changeOpts,
	}
}
