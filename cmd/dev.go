package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/bundler"
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
	rootComponents  map[string]*srcpack.Component
	sourceMap       *dependtree.DependencySourceMap

	lastProcessedFile *proccessedChangeRequest
	m                 *sync.Mutex
	packSettings      *srcpack.Packer
}

var watcher *fsnotify.Watcher

func verifyComponentPath(in string) string {
	skip := 0
	for _, c := range in {
		if c == '.' || c == '/' {
			skip += 1
			continue
		}

		break
	}

	return in[skip:]
}

func createSession(ctx context.Context, settings *internal.GenPagesSettings) (*devSession, error) {
	err := settings.CleanPathing()
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, bundler.BundlerID, settings.BundlerMode)

	// we use the empty logger here to to prevent the initial build from being shown
	// during the creation of the dev session
	lib, err := settings.PackWebDir(ctx, log.NewEmptyLogger())
	if err != nil {
		return nil, err
	}

	rootComponents := make(map[string]*srcpack.Component)
	for _, p := range lib.Pages {
		// verify that the path is clean before we apply it to the root component map
		path := verifyComponentPath(p.OriginalFilePath())

		rootComponents[path] = p
	}

	sourceMap, err := srcpack.New(settings.WebDir, lib.Pages, settings.WebDir)
	if err != nil {
		return nil, err
	}

	// we use the default logger here to log the dev build events
	_, packSettings := settings.DefaultPacker(ctx, log.NewDefaultLogger())

	return &devSession{
		pageGenSettings:   settings,
		rootComponents:    rootComponents,
		sourceMap:         sourceMap,
		lastProcessedFile: &proccessedChangeRequest{},
		m:                 &sync.Mutex{},
		packSettings:      packSettings,
	}, nil
}

// executeChangeRequest attempts to find the file in the component tree, if found it
// will repack each of the branches that dependens on it.
func (s *devSession) executeChangeRequest(file string, timeoutDuration time.Duration) error {
	if s.pageGenSettings.UseDebug {
		s.packSettings.Logger.Info(fmt.Sprintf("change detected → %s", file))
	}

	// if this file has been recently processed (specificed by the timeout flag), do not process it.
	if file == s.lastProcessedFile.FileName &&
		time.Since(s.lastProcessedFile.ProcessedAt).Seconds() < timeoutDuration.Seconds() {

		if s.pageGenSettings.UseDebug {
			s.packSettings.Logger.Info(fmt.Sprintf("change not excepted → %s (too recently processed)", file))
		}
		return nil
	}

	component := s.rootComponents[file]

	// if component is one of the root components, we will just repack that component
	if component != nil {
		if s.pageGenSettings.UseDebug {
			s.packSettings.Logger.Info(fmt.Sprintf("change found → %s (root)", file))
		}

		err := s.pageGenSettings.Repack(component, srcpack.NewSyncHook(log.NewDefaultLogger()))
		if err != nil {
			return err
		}

		s.m.Lock()
		s.lastProcessedFile = &proccessedChangeRequest{
			FileName:    file,
			ProcessedAt: time.Now(),
		}
		s.m.Unlock()

		return nil
	}

	// component is not root, we need to find in which tree(s) the component exists & execute
	// a repack for each of those components & their dependent branches.
	sources := s.sourceMap.FindRoot(file)

	if s.pageGenSettings.UseDebug {
		s.packSettings.Logger.Info(fmt.Sprintf("%d branch(s) found", len(sources)))
	}

	// we iterate through each of the root sources for the source
	activeNodes := make([]*srcpack.Component, 0)
	for _, source := range sources {
		if s.pageGenSettings.UseDebug {
			s.packSettings.Logger.Info(fmt.Sprintf("change found → %s (branch)", source))
		}

		source = verifyComponentPath(source)
		component = s.rootComponents[source]

		activeNodes = append(activeNodes, component)
	}

	cl := srcpack.PackedComponentList(activeNodes)

	return cl.RepackMany(log.NewDefaultLogger())
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
		logger := log.NewDefaultLogger()

		if viper.GetBool("usedebug") {
			logger.Warn("debug mode enabled")
		}

		logger.Info("starting dev server...")

		as := &internal.GenPagesSettings{
			PackageName:    viper.GetString("pacname"),
			OutDir:         viper.GetString("out"),
			WebDir:         viper.GetString("webdir"),
			BundlerMode:    viper.GetString("mode"),
			NodeModulePath: viper.GetString("nodemod"),
			PublicDir:      viper.GetString("publicdir"),
			UseDebug:       viper.GetBool("usedebug"),
		}

		s, err := createSession(context.Background(), as)
		if err != nil {
			panic(err)
		}

		logger.Success("dev server started successfully\n")

		watcher, _ = fsnotify.NewWatcher()
		defer watcher.Close()

		if err := filepath.Walk("./", watchDir); err != nil {
			panic("invalid walk on watchDir")
		}

		done := make(chan bool)

		go func() {
			for {
				time.Sleep(time.Duration(viper.GetInt("timeout")) * time.Millisecond)

				select {
				case e := <-watcher.Events:
					if !strings.Contains(e.Name, "node_modules") && !strings.Contains(e.Name, ".orbit") {
						go s.executeChangeRequest(e.Name, time.Duration(viper.GetInt("samefiletimeout"))*time.Millisecond)
					}
				case err := <-watcher.Errors:
					panic(fmt.Sprintf("watcher failed %s", err.Error()))
				}
			}
		}()

		<-done
	},
}

func init() {
	var timeoutDuration int
	var samefileTimeout int

	devCMD.PersistentFlags().IntVar(&timeoutDuration, "timeout", 2000, "specifies the timeout duration in milliseconds until a change will be detected")
	viper.BindPFlag("timeout", devCMD.PersistentFlags().Lookup("timeout"))

	devCMD.PersistentFlags().IntVar(&samefileTimeout, "samefiletimeout", 2000, "specifies the timeout duration in milliseconds until a change will be detected for repeating files")
	viper.BindPFlag("samefiletimeout", devCMD.PersistentFlags().Lookup("samefiletimeout"))

	RootCMD.AddCommand(devCMD)
}
