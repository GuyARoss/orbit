// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.
package internal

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sync"
	"time"

	"github.com/GuyARoss/orbit/internal/assets"
	"github.com/GuyARoss/orbit/internal/libout"
	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/bundler"
	dependtree "github.com/GuyARoss/orbit/pkg/depend_tree"
	"github.com/GuyARoss/orbit/pkg/fsutils"
	"github.com/GuyARoss/orbit/pkg/log"
	webwrapper "github.com/GuyARoss/orbit/pkg/web_wrapper"
	"github.com/fsnotify/fsnotify"
)

type proccessedChangeRequest struct {
	ProcessedAt time.Time
	FileName    string
}

type SessionOpts struct {
	UseDebug      bool
	WebDir        string
	Mode          string
	Pacname       string
	OutDir        string
	NodeModDir    string
	PublicDir     string
	HotReloadPort int
}

type devSession struct {
	*SessionOpts

	RootComponents    map[string]*srcpack.Component
	SourceMap         *dependtree.DependencySourceMap
	lastProcessedFile *proccessedChangeRequest
	m                 *sync.Mutex
	packer            *srcpack.Packer
}

// verifyComponentPath is a utility that verifies
// that the provided path is a file valid path
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

// CreateSession creates a new active dev session with the following
// - a flat tree represented by a map of the root page in component form
// - initializes the development build process
func CreateSession(ctx context.Context, opts *SessionOpts) (*devSession, error) {
	ats, err := assets.AssetKeys()
	if err != nil {
		panic(err)
	}

	err = OrbitFileStructure(&FileStructureOpts{
		PackageName: opts.Pacname,
		OutDir:      opts.OutDir,
		Assets:      []fs.DirEntry{ats.AssetKey(assets.WebPackConfig)},
		Dist:        []fs.DirEntry{ats.AssetKey(assets.HotReload)},
	})

	if err != nil {
		return nil, err
	}

	c, err := CachedEnvFromFile(fsutils.NormalizePath(fmt.Sprintf("%s/%s/orb_env.go", opts.OutDir, opts.Pacname)))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	packer := srcpack.NewDefaultPacker(log.NewEmptyLogger(), &srcpack.DefaultPackerOpts{
		WebDir:           opts.WebDir,
		BundlerMode:      opts.Mode,
		NodeModuleDir:    opts.NodeModDir,
		CachedBundleKeys: c,
	})

	pageFiles := fsutils.DirFiles(fsutils.NormalizePath(fmt.Sprintf("%s/pages", opts.WebDir)))
	components, err := packer.PackMany(pageFiles)
	if err != nil {
		panic(err)
	}

	bg := libout.New(&libout.BundleGroupOpts{
		PackageName:   opts.Pacname,
		BaseBundleOut: ".orbit/dist",
		BundleMode:    opts.Mode,
		PublicDir:     opts.PublicDir,
		HotReloadPort: opts.HotReloadPort,
	})

	ctx = context.WithValue(ctx, bundler.BundlerID, opts.Mode)

	bg.AcceptComponents(ctx, components, &webwrapper.CacheDOMOpts{
		CacheDir:  ".orbit/dist",
		WebPrefix: "/p/",
	})

	err = bg.WriteLibout(libout.NewGOLibout(
		ats.AssetKey(assets.Tests),
		ats.AssetKey(assets.PrimaryPackage),
	), &libout.FilePathOpts{
		TestFile: fsutils.NormalizePath(fmt.Sprintf("%s/%s/orb_test.go", opts.OutDir, opts.Pacname)),
		EnvFile:  fsutils.NormalizePath(fmt.Sprintf("%s/%s/orb_env.go", opts.OutDir, opts.Pacname)),
		HTTPFile: fsutils.NormalizePath(fmt.Sprintf("%s/%s/orb_http.go", opts.OutDir, opts.Pacname)),
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
		RootComponents:    rootComponents,
		SourceMap:         sourceMap,
		lastProcessedFile: &proccessedChangeRequest{},
		m:                 &sync.Mutex{},
		packer:            packer.ReattachLogger(log.NewDefaultLogger()),
	}, nil
}

func (s *devSession) DirectFileChangeRequest(file string, timeoutDuration time.Duration, sh *srcpack.SyncHook) error {
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

	component := s.RootComponents[file]

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

	return nil
}

func (s *devSession) IndirectFileChangeRequest(indirectFile string, directComponentBundleKey string, timeoutDuration time.Duration, sh *srcpack.SyncHook) error {
	// if this file has been recently processed (specificed by the timeout flag), do not process it.
	if indirectFile == s.lastProcessedFile.FileName &&
		time.Since(s.lastProcessedFile.ProcessedAt).Seconds() < timeoutDuration.Seconds() {

		if s.UseDebug {
			s.packer.Logger.Info(fmt.Sprintf("change not excepted → %s (too recently processed)", indirectFile))
		}
		return nil
	}

	// component is not root, we need to find in which tree(s) the component exists & execute
	// a repack for each of those components & their dependent branches.
	sources := s.SourceMap.FindRoot(indirectFile)

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
		component := s.RootComponents[source]

		if component.BundleKey == directComponentBundleKey {
			if s.UseDebug {
				s.packer.Logger.Info(fmt.Sprintf("change found → %s (root)", indirectFile))
			}

			sh.WrapFunc(component.OriginalFilePath(), func() { component.Repack() })

			s.m.Lock()
			s.lastProcessedFile = &proccessedChangeRequest{
				FileName:    indirectFile,
				ProcessedAt: time.Now(),
			}
			s.m.Unlock()

			return nil
		}

		activeNodes = append(activeNodes, component)
	}

	cl := srcpack.PackedComponentList(activeNodes)

	return cl.RepackMany(log.NewDefaultLogger())
}

// watchDir is a utility function used by the file path walker that applies
// each sub directory found under a path to the file watcher
func WatchDir(watcher *fsnotify.Watcher) func(path string, fi os.FileInfo, err error) error {
	return func(path string, fi os.FileInfo, err error) error {
		if fi.Mode().IsDir() {
			return watcher.Add(path)
		}

		return nil
	}
}
