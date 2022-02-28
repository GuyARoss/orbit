package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/internal/assets"
	"github.com/GuyARoss/orbit/internal/libout"
	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/bundler"
	dependtree "github.com/GuyARoss/orbit/pkg/depend_tree"
	"github.com/GuyARoss/orbit/pkg/fsutils"
	"github.com/GuyARoss/orbit/pkg/log"
	webwrapper "github.com/GuyARoss/orbit/pkg/web_wrapper"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type proccessedChangeRequest struct {
	ProcessedAt time.Time
	FileName    string
}

type SessionOpts struct {
	UseDebug bool
	WebDir   string
}

type devSession struct {
	*SessionOpts

	rootComponents    map[string]*srcpack.Component
	sourceMap         *dependtree.DependencySourceMap
	lastProcessedFile *proccessedChangeRequest
	m                 *sync.Mutex
	packer            *srcpack.Packer
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

func createSession(ctx context.Context, opts *SessionOpts) (*devSession, error) {
	ats, err := assets.AssetKeys()
	if err != nil {
		panic(err)
	}

	err = internal.OrbitFileStructure(&internal.FileStructureOpts{
		PackageName: viper.GetString("pacname"),
		OutDir:      viper.GetString("out"),
		Assets:      []fs.DirEntry{ats.AssetKey(assets.WebPackConfig)},
	})

	if err != nil {
		return nil, err
	}

	packer := srcpack.NewDefaultPacker(log.NewEmptyLogger(), &srcpack.DefaultPackerOpts{
		WebDir:        viper.GetString("webdir"),
		BundlerMode:   viper.GetString("mode"),
		NodeModuleDir: viper.GetString("nodemod"),
	})

	pageFiles := fsutils.DirFiles(fsutils.NormalizePath(fmt.Sprintf("%s/pages", viper.GetString("webdir"))))
	components, err := packer.PackMany(pageFiles)
	if err != nil {
		panic(err)
	}

	bg := libout.New(&libout.BundleGroupOpts{
		PackageName:   viper.GetString("pacname"),
		BaseBundleOut: ".orbit/dist",
		BundleMode:    string(viper.GetString("mode")),
		PublicDir:     viper.GetString("publicdir"),
	})

	ctx = context.WithValue(ctx, bundler.BundlerID, viper.GetString("mode"))

	bg.AcceptComponents(ctx, components, &webwrapper.CacheDOMOpts{
		CacheDir:  ".orbit/dist",
		WebPrefix: "/p/",
	})

	err = bg.WriteLibout(libout.NewGOLibout(
		ats.AssetKey(assets.Tests),
		ats.AssetKey(assets.PrimaryPackage),
	), &libout.FilePathOpts{
		TestFile: fsutils.NormalizePath(fmt.Sprintf("%s/%s/orb_test.go", viper.GetString("webdir"), viper.GetString("pacname"))),
		EnvFile:  fsutils.NormalizePath(fmt.Sprintf("%s/%s/orb_env.go", viper.GetString("webdir"), viper.GetString("pacname"))),
		HTTPFile: fsutils.NormalizePath(fmt.Sprintf("%s/%s/orb_http.go", viper.GetString("webdir"), viper.GetString("pacname"))),
	})
	if err != nil {
		return nil, err
	}

	sourceMap, err := srcpack.New(opts.WebDir, components, opts.WebDir)
	if err != nil {
		return nil, err
	}

	rootComponents := make(map[string]*srcpack.Component)
	for _, p := range components {
		// verify that the path is clean before we apply it to the root component map
		path := verifyComponentPath(p.OriginalFilePath())

		rootComponents[path] = p
	}

	return &devSession{
		SessionOpts:       opts,
		rootComponents:    rootComponents,
		sourceMap:         sourceMap,
		lastProcessedFile: &proccessedChangeRequest{},
		m:                 &sync.Mutex{},
		packer:            packer.ReattachLogger(log.NewDefaultLogger()),
	}, nil
}

// executeChangeRequest attempts to find the file in the component tree, if found it
// will repack each of the branches that dependens on it.
func (s *devSession) executeChangeRequest(file string, timeoutDuration time.Duration, sh *srcpack.SyncHook) error {
	if s.UseDebug {
		s.packer.Logger.Info(fmt.Sprintf("change detected → %s", file))
	}

	// if this file has been recently processed (specificed by the timeout flag), do not process it.
	if file == s.lastProcessedFile.FileName &&
		time.Since(s.lastProcessedFile.ProcessedAt).Seconds() < timeoutDuration.Seconds() {

		if s.UseDebug {
			s.packer.Logger.Info(fmt.Sprintf("change not excepted → %s (too recently processed)", file))
		}
		return nil
	}

	component := s.rootComponents[file]

	// if component is one of the root components, we will just repack that component
	if component != nil {
		if s.UseDebug {
			s.packer.Logger.Info(fmt.Sprintf("change found → %s (root)", file))
		}

		sh.WrapFunc(component.OriginalFilePath(), func() { component.Repack() })

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

	if s.UseDebug {
		s.packer.Logger.Info(fmt.Sprintf("%d branch(s) found", len(sources)))
	}

	// we iterate through each of the root sources for the source
	activeNodes := make([]*srcpack.Component, 0)
	for _, source := range sources {
		if s.UseDebug {
			s.packer.Logger.Info(fmt.Sprintf("change found → %s (branch)", source))
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

		s, err := createSession(context.Background(), &SessionOpts{
			UseDebug: viper.GetBool("usedebug"),
			WebDir:   viper.GetString("webdir"),
		})
		if err != nil {
			panic(err)
		}

		logger.Success("dev server started successfully\n")

		watcher, _ = fsnotify.NewWatcher()
		defer watcher.Close()

		if err := filepath.Walk(fsutils.NormalizePath("./"), watchDir); err != nil {
			panic("invalid walk on watchDir")
		}

		done := make(chan bool)

		go func() {
			sh := srcpack.NewSyncHook(log.NewDefaultLogger())

			for {
				time.Sleep(time.Duration(viper.GetInt("timeout")) * time.Millisecond)

				select {
				case e := <-watcher.Events:
					if !strings.Contains(e.Name, "node_modules") && !strings.Contains(e.Name, ".orbit") {
						go s.executeChangeRequest(e.Name, time.Duration(viper.GetInt("samefiletimeout"))*time.Millisecond, sh)
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
