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
	"time"

	"github.com/GuyARoss/orbit/internal/assets"
	"github.com/GuyARoss/orbit/internal/libout"
	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/bundler"
	dependtree "github.com/GuyARoss/orbit/pkg/depend_tree"
	"github.com/GuyARoss/orbit/pkg/fsutils"
	"github.com/GuyARoss/orbit/pkg/hotreload"
	"github.com/GuyARoss/orbit/pkg/log"
	webwrapper "github.com/GuyARoss/orbit/pkg/web_wrapper"
)

// SessionOpts are options used for creating a new session
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

// proccessedChangeRequest is the most recent file change that has happended within the development process
type proccessedChangeRequest struct {
	ProcessedAt time.Time
	FileName    string
}

// devSession is the internal state for processing change requests during a development process
type devSession struct {
	*SessionOpts

	RootComponents    map[string]srcpack.PackComponent
	SourceMap         dependtree.DependencySourceMap
	lastProcessedFile *proccessedChangeRequest
	packer            *srcpack.Packer
}

// ChangeRequestOpts options used for processing a change request
type ChangeRequestOpts struct {
	SafeFileTimeout time.Duration
	HotReload       hotreload.HotReloader
	Hook            *srcpack.SyncHook
}

var ErrFileTooRecentlyProcessed = errors.New("change not accepted, file too recently processed")

// ProcessChangeRequest will determine which type of change request is required for computation of the request file
func (s *devSession) DoChangeRequest(filePath string, opts *ChangeRequestOpts) error {
	// if this file has been recently processed (specificed by the timeout flag), do not process it.
	if filePath == s.lastProcessedFile.FileName &&
		time.Since(s.lastProcessedFile.ProcessedAt).Seconds() < opts.SafeFileTimeout.Seconds() {

		if s.UseDebug {
			s.packer.Logger.Info(fmt.Sprintf("change not excepted → %s (too recently processed)", filePath))
		}

		return ErrFileTooRecentlyProcessed
	}

	root := s.RootComponents[filePath]

	// if components' bundle is the current bundle that is open in the browser
	// recompute bundle and send refresh signal back to browser
	if root != nil && root.BundleKey() == opts.HotReload.CurrentBundleKey() {
		s.DirectFileChangeRequest(filePath, root, opts.SafeFileTimeout, opts.Hook)

		err := opts.HotReload.ReloadSignal()
		if err != nil {
			return err
		}

		// no need to continue, root file has already been processed.
		return nil
	}

	// component may exist as a page depencency, if so, recompute and send refresh signal
	if len(s.SourceMap.FindRoot(filePath)) > 0 {
		s.IndirectFileChangeRequest(filePath, opts.HotReload.CurrentBundleKey(), opts.SafeFileTimeout, opts.Hook)

		err := opts.HotReload.ReloadSignal()
		if err != nil {
			return err
		}
	}

	return nil
}

// DirectFileChangeRequest processes a change request for a root component directly
func (s *devSession) DirectFileChangeRequest(filePath string, component srcpack.PackComponent, timeoutDuration time.Duration, sh *srcpack.SyncHook) error {
	// if component is one of the root components, we will just repack that component
	if component != nil {
		if s.UseDebug {
			s.packer.Logger.Info(fmt.Sprintf("change found → %s (root)", filePath))
		}

		sh.WrapFunc(component.OriginalFilePath(), func() { component.Repack() })

		s.lastProcessedFile = &proccessedChangeRequest{
			FileName:    filePath,
			ProcessedAt: time.Now(),
		}

		return nil
	}

	return nil
}

// IndirectFileChangeRequest processes a change request for a file that may be a dependency of a root component
func (s *devSession) IndirectFileChangeRequest(indirectFile string, parentBundleKey string, timeoutDuration time.Duration, sh *srcpack.SyncHook) error {
	// component is not root, we need to find in which tree(s) the component exists & execute
	// a repack for each of those components & their dependent branches.
	sources := s.SourceMap.FindRoot(indirectFile)

	if s.UseDebug {
		s.packer.Logger.Info(fmt.Sprintf("%d branch(s) found", len(sources)))
	}

	// we iterate through each of the root sources for the source
	activeNodes := make([]srcpack.PackComponent, 0)
	for _, source := range sources {
		if s.UseDebug {
			s.packer.Logger.Info(fmt.Sprintf("change found → %s (branch)", source))
		}

		source = verifyComponentPath(source)
		component := s.RootComponents[source]

		if component.BundleKey() == parentBundleKey {
			if s.UseDebug {
				s.packer.Logger.Info(fmt.Sprintf("change found → %s (root)", indirectFile))
			}

			sh.WrapFunc(component.OriginalFilePath(), func() { component.Repack() })

			s.lastProcessedFile = &proccessedChangeRequest{
				FileName:    indirectFile,
				ProcessedAt: time.Now(),
			}

			return nil
		}

		activeNodes = append(activeNodes, component)
	}

	cl := srcpack.PackedComponentList(activeNodes)

	return cl.RepackMany(log.NewDefaultLogger())
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

// NewDevSession creates a new active dev session with the following:
//  1. a flat tree represented by a map of the root page in component form
//  2. initializes the development build process
func NewDevSession(ctx context.Context, opts *SessionOpts) (*devSession, error) {
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

	rootComponents := make(map[string]srcpack.PackComponent)
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
		packer:            packer.ReattachLogger(log.NewDefaultLogger()),
	}, nil
}
