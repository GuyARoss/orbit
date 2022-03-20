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
	"strings"
	"time"

	"github.com/GuyARoss/orbit/internal/assets"
	"github.com/GuyARoss/orbit/internal/libout"
	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/bundler"
	dependtree "github.com/GuyARoss/orbit/pkg/depend_tree"
	"github.com/GuyARoss/orbit/pkg/fsutils"
	"github.com/GuyARoss/orbit/pkg/hotreload"
	"github.com/GuyARoss/orbit/pkg/log"
	"github.com/GuyARoss/orbit/pkg/webwrap"
)

// SessionOpts are options used for creating a new session
type SessionOpts struct {
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
	packer            srcpack.Packer
	lastProcessedFile *proccessedChangeRequest
	libout            libout.BundleWriter
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

		return ErrFileTooRecentlyProcessed
	}

	root := s.RootComponents[filePath]

	// if components' bundle is the current bundle that is open in the browser
	// recompute bundle and send refresh signal back to browser
	if root != nil && root.BundleKey() == opts.HotReload.CurrentBundleKey() {
		err := s.DirectFileChangeRequest(filePath, root, opts)
		if err != nil {
			return err
		}

		err = opts.HotReload.ReloadSignal()
		if err != nil {
			return err
		}

		// no need to continue, root file has already been processed.
		return nil
	}

	// if we assume that this is a new page, attempt to build it and add it to preexisting context
	// @@todo(guy) magic string : "pages" allow support for this keyword from a flag
	if strings.Contains(filePath, "pages/") {
		err := s.NewPageFileChangeRequest(context.Background(), filePath)

		return err
	}

	// component may exist as a page depencency, if so, recompute and send refresh signal
	sources := s.SourceMap.FindRoot(filePath)
	if len(sources) > 0 {
		// component is not root, we need to find in which tree(s) the component exists & execute
		// a repack for each of those components & their dependent branches.
		err := s.IndirectFileChangeRequest(sources, filePath, opts)
		if err != nil {
			return err
		}

		err = opts.HotReload.ReloadSignal()
		if err != nil {
			return err
		}
	}

	return nil
}

// DirectFileChangeRequest processes a change request for a root component directly
func (s *devSession) DirectFileChangeRequest(filePath string, component srcpack.PackComponent, opts *ChangeRequestOpts) error {
	// if component is one of the root components, we will just repack that component
	if component == nil {
		return nil
	}

	opts.Hook.WrapFunc(component.OriginalFilePath(), func() { component.Repack() })

	s.lastProcessedFile = &proccessedChangeRequest{
		FileName:    filePath,
		ProcessedAt: time.Now(),
	}

	return nil
}

// IndirectFileChangeRequest processes a change request for a file that may be a dependency of a root component
func (s *devSession) IndirectFileChangeRequest(sources []string, indirectFile string, opts *ChangeRequestOpts) error {
	// we iterate through each of the root sources for the source until the component bundle has been found.
	for _, source := range sources {
		source = verifyComponentPath(source)
		component := s.RootComponents[source]

		if component.BundleKey() == opts.HotReload.CurrentBundleKey() {
			opts.Hook.WrapFunc(component.OriginalFilePath(), func() { component.Repack() })

			s.lastProcessedFile = &proccessedChangeRequest{
				FileName:    indirectFile,
				ProcessedAt: time.Now(),
			}

			return nil
		}
	}

	return nil
}

// NewPageFileChangeRequest processes a change request for file that is detected as a new page
func (s *devSession) NewPageFileChangeRequest(ctx context.Context, file string) error {
	ats, err := assets.AssetKeys()
	if err != nil {
		panic(err)
	}

	component, err := s.packer.PackSingle(log.NewEmptyLogger(), file)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, bundler.BundlerID, s.Mode)
	s.libout.AcceptComponent(ctx, component, &webwrap.CacheDOMOpts{
		CacheDir:  ".orbit/dist",
		WebPrefix: "/p/",
	})

	err = s.libout.WriteLibout(libout.NewGOLibout(
		ats.AssetKey(assets.Tests),
		ats.AssetKey(assets.PrimaryPackage),
	), &libout.FilePathOpts{
		TestFile: fsutils.NormalizePath(fmt.Sprintf("%s/%s/orb_test.go", s.OutDir, s.Pacname)),
		EnvFile:  fsutils.NormalizePath(fmt.Sprintf("%s/%s/orb_env.go", s.OutDir, s.Pacname)),
		HTTPFile: fsutils.NormalizePath(fmt.Sprintf("%s/%s/orb_http.go", s.OutDir, s.Pacname)),
	})
	if err != nil {
		return err
	}

	sourceMap, err := srcpack.New(s.WebDir, []srcpack.PackComponent{component}, s.WebDir)
	if err != nil {
		return err
	}

	s.SourceMap = s.SourceMap.Merge(sourceMap)
	s.RootComponents[verifyComponentPath(component.OriginalFilePath())] = component

	return nil
}

// verifyComponentPath is a utility that verifies that the provided path is a file valid path
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

	// @@todo(guy) magic string : "pages" allow support for this keyword from a flag
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

	bg.AcceptComponents(ctx, components, &webwrap.CacheDOMOpts{
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
		libout:            bg,
	}, nil
}
