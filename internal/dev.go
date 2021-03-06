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
	allocatedstack "github.com/GuyARoss/orbit/pkg/allocated_stack"
	dependtree "github.com/GuyARoss/orbit/pkg/depend_tree"
	"github.com/GuyARoss/orbit/pkg/fsutils"
	"github.com/GuyARoss/orbit/pkg/hotreload"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/log"
	parseerror "github.com/GuyARoss/orbit/pkg/parse_error"
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

// devSession is the internal state for processing change requests during a development process
type devSession struct {
	*SessionOpts

	RootComponents srcpack.PackComponentFileMap
	SourceMap      dependtree.DependencySourceMap
	packer         srcpack.Packer
	libout         libout.BundleWriter
	ChangeRequest  *changeRequest
}

// ChangeRequestOpts options used for processing a change request
type ChangeRequestOpts struct {
	SafeFileTimeout time.Duration
	HotReload       hotreload.HotReloader
	Hook            *srcpack.SyncHook
	Parser          jsparse.JSParser
}

var ErrFileTooRecentlyProcessed = errors.New("change not accepted, file too recently processed")

// DoBundleKeyChangeRequest processes a change request for a bundle key
func (s *devSession) DoBundleKeyChangeRequest(bundleKey string, opts *ChangeRequestOpts) error {
	component := s.RootComponents.FindBundleKey(bundleKey)
	err := s.DirectFileChangeRequest("", component, opts)

	if err != nil {
		return parseerror.FromError(err, component.OriginalFilePath())
	}

	err = opts.HotReload.ReloadSignal()
	if err != nil {
		return parseerror.FromError(err, component.OriginalFilePath())
	}
	return nil
}

// ProcessChangeRequest will determine which type of change request is required for computation of the request file
func (s *devSession) DoFileChangeRequest(filePath string, opts *ChangeRequestOpts) error {
	// if this file has been recently processed (specificed by the timeout flag), do not process it.
	if !s.ChangeRequest.IsWithinRage(filePath, opts.SafeFileTimeout) {
		return parseerror.FromError(ErrFileTooRecentlyProcessed, filePath)
	}

	// file detected in the orbit output, we don't want to process any of these
	if strings.Contains(filePath, ".orbit") {
		return nil
	}

	root := s.RootComponents[filePath]
	// if components' bundle is the current bundle that is open in the browser
	// recompute bundle and send refresh signal back to browser
	if root != nil && opts.HotReload.IsActiveBundle(root.BundleKey()) {
		err := s.DirectFileChangeRequest(filePath, root, opts)
		if err != nil {
			return parseerror.FromError(err, filePath)
		}

		err = opts.HotReload.ReloadSignal()
		if err != nil {
			return parseerror.FromError(err, filePath)
		}

		// no need to continue, root file has already been processed.
		return nil
	}

	// if we assume that this is a new page, attempt to build it and add it to preexisting context
	// @@todo(guy) magic string : "pages" allow support for this keyword from a flag
	if strings.Contains(filePath, "pages/") {
		err := s.NewPageFileChangeRequest(context.Background(), filePath)

		if err != nil {
			return parseerror.FromError(err, filePath)
		}
	}

	// component may exist as a page dependency, if so, recompute and send refresh signal
	sources := s.SourceMap.FindRoot(filePath)
	if len(sources) > 0 {
		// component is not root, we need to find in which tree(s) the component exists & execute
		// a repack for each of those components & their dependent branches.
		err := s.IndirectFileChangeRequest(sources, filePath, opts)
		if err != nil {
			return parseerror.FromError(err, filePath)
		}

		err = opts.HotReload.ReloadSignal()
		if err != nil {
			return parseerror.FromError(err, filePath)
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

	if filePath == "" {
		filePath = component.OriginalFilePath()
	}

	opts.Hook.WrapFunc(component.OriginalFilePath(), func() *webwrap.WrapStats {
		component.Repack()

		return component.WebWrapper().Stats()
	})

	s.ChangeRequest.Push(filePath, component.BundleKey())

	sourceMap, err := srcpack.New(s.WebDir, []srcpack.PackComponent{component}, &srcpack.NewSourceMapOpts{
		Parser:     opts.Parser,
		WebDirPath: s.WebDir,
	})
	if err != nil {
		return err
	}

	s.SourceMap = s.SourceMap.MergeOverKey(sourceMap)

	return nil
}

// IndirectFileChangeRequest processes a change request for a file that may be a dependency of a root component
func (s *devSession) IndirectFileChangeRequest(sources []string, indirectFile string, opts *ChangeRequestOpts) error {
	// we iterate through each of the root sources for the source until the component bundle has been found.
	for _, source := range sources {
		component := s.RootComponents.Find(source)

		if !opts.HotReload.IsActiveBundle(component.BundleKey()) {
			continue
		}

		opts.Hook.WrapFunc(component.OriginalFilePath(), func() *webwrap.WrapStats {
			component.Repack()

			return component.WebWrapper().Stats()
		})

		s.ChangeRequest.Push(indirectFile, component.BundleKey())

		sourceMap, err := srcpack.New(s.WebDir, []srcpack.PackComponent{component}, &srcpack.NewSourceMapOpts{
			Parser:     opts.Parser,
			WebDirPath: s.WebDir,
		})
		if err != nil {
			return err
		}

		s.SourceMap = s.SourceMap.MergeOverKey(sourceMap)
		return nil
	}

	return nil
}

var ErrCannotBuildAssetKeys = errors.New("cannot build asset keys")

// NewPageFileChangeRequest processes a change request for file that is detected as a new page
func (s *devSession) NewPageFileChangeRequest(ctx context.Context, file string) error {
	ats, err := assets.AssetKeys()
	if err != nil {
		return ErrCannotBuildAssetKeys
	}

	component, err := s.packer.PackSingle(log.NewEmptyLogger(), file)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, webwrap.BundlerID, s.Mode)
	s.libout.AcceptComponent(ctx, component, &webwrap.CacheDOMOpts{
		CacheDir:  ".orbit/dist",
		WebPrefix: "/p/",
	})

	if err = s.libout.WriteLibout(libout.NewGOLibout(
		ats.AssetKey(assets.Tests),
		ats.AssetKey(assets.PrimaryPackage),
	), &libout.FilePathOpts{
		TestFile: fmt.Sprintf("%s/%s/orb_test.go", s.OutDir, s.Pacname),
		EnvFile:  fmt.Sprintf("%s/%s/orb_env.go", s.OutDir, s.Pacname),
		HTTPFile: fmt.Sprintf("%s/%s/orb_http.go", s.OutDir, s.Pacname),
	}); err != nil {
		return err
	}

	sourceMap, err := srcpack.New(s.WebDir, []srcpack.PackComponent{component}, &srcpack.NewSourceMapOpts{
		Parser:     &jsparse.JSFileParser{},
		WebDirPath: s.WebDir,
	})
	if err != nil {
		return err
	}

	s.SourceMap = s.SourceMap.Merge(sourceMap)
	s.RootComponents.Set(component)

	s.ChangeRequest.Push(file, component.BundleKey())

	return nil
}

// New creates a new active dev session with the following:
//  1. a flat tree represented by a map of the root page in component form
//  2. initializes the development build process
func New(ctx context.Context, opts *SessionOpts) (*devSession, error) {
	ats, err := assets.AssetKeys()
	if err != nil {
		panic(err)
	}

	err = OrbitFileStructure(&FileStructureOpts{
		PackageName: opts.Pacname,
		OutDir:      opts.OutDir,
		Assets: []fs.DirEntry{
			ats.AssetEntry(assets.WebPackConfig),
			ats.AssetEntry(assets.SSRProtoFile),
			ats.AssetEntry(assets.JsWebPackConfig),
		},
		Dist: []fs.DirEntry{ats.AssetEntry(assets.HotReload)},
	})

	if err != nil {
		return nil, err
	}

	c, err := CachedEnvFromFile(fmt.Sprintf("%s/%s/orb_env.go", opts.OutDir, opts.Pacname))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	packer := srcpack.NewDefaultPacker(log.NewEmptyLogger(), &srcpack.DefaultPackerOpts{
		WebDir:              opts.WebDir,
		BundlerMode:         opts.Mode,
		NodeModuleDir:       opts.NodeModDir,
		CachedBundleKeys:    c,
		SkipFirstPassBundle: true,
	})

	// @@todo(guy) magic string : "pages" allow support for this keyword from a flag
	pageFiles := fsutils.DirFiles(fmt.Sprintf("%s/pages", opts.WebDir))
	components, err := packer.PackMany(pageFiles)
	if err != nil {
		return nil, err
	}

	bg := libout.New(&libout.BundleGroupOpts{
		PackageName:   opts.Pacname,
		BaseBundleOut: ".orbit/dist",
		BundleMode:    opts.Mode,
		PublicDir:     opts.PublicDir,
		HotReloadPort: opts.HotReloadPort,
	})

	ctx = context.WithValue(ctx, webwrap.BundlerID, opts.Mode)

	if err = bg.AcceptComponents(ctx, components, &webwrap.CacheDOMOpts{
		CacheDir:  ".orbit/dist",
		WebPrefix: "/p/",
	}); err != nil {
		return nil, err
	}

	err = bg.WriteLibout(libout.NewGOLibout(
		ats.AssetKey(assets.Tests),
		ats.AssetKey(assets.PrimaryPackage),
	), &libout.FilePathOpts{
		TestFile: fmt.Sprintf("%s/%s/orb_test.go", opts.OutDir, opts.Pacname),
		EnvFile:  fmt.Sprintf("%s/%s/orb_env.go", opts.OutDir, opts.Pacname),
		HTTPFile: fmt.Sprintf("%s/%s/orb_http.go", opts.OutDir, opts.Pacname),
	})
	if err != nil {
		return nil, err
	}

	sourceMap, err := srcpack.New(opts.WebDir, components, &srcpack.NewSourceMapOpts{
		Parser:     &jsparse.JSFileParser{},
		WebDirPath: opts.WebDir,
	})
	if err != nil {
		return nil, err
	}

	rootComponents := make(srcpack.PackComponentFileMap)
	for _, p := range components {
		rootComponents.Set(p)
	}

	return &devSession{
		SessionOpts:    opts,
		RootComponents: rootComponents,
		SourceMap:      sourceMap,
		packer:         packer,
		libout:         bg,
		ChangeRequest: &changeRequest{
			changeRequests: allocatedstack.New(10),
		},
	}, nil
}

// changeRequest holds the most recent file changes that have happened in the development cycle
type changeRequest struct {
	LastProcessedAt time.Time
	LastFileName    string

	changeRequests *allocatedstack.Stack
}

func (c *changeRequest) ExistsInCache(file string) bool {
	return c.changeRequests.Contains(file)
}

func (c *changeRequest) Push(fileName string, bundleKey string) {
	c.LastFileName = fileName
	c.LastProcessedAt = time.Now()

	c.changeRequests.Add(bundleKey)
}

func (c *changeRequest) IsWithinRage(file string, t time.Duration) bool {
	if c != nil && file == c.LastFileName {
		return time.Since(c.LastProcessedAt).Seconds() > t.Seconds()
	}

	return true
}
