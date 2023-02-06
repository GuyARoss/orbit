// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package webwrap

import (
	"context"
	"crypto/md5"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/GuyARoss/orbit/pkg/embedutils"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/log"
)

type JSWebWrapper interface {
	RequiredBodyDOMElements(context.Context, *CacheDOMOpts) []string
	Setup(context.Context, *BundleOpts) (*BundledResource, error)
	Apply(jsparse.JSDocument) (map[string]jsparse.JSDocument, error)
	Version() string
	Stats() *WrapStats
	Bundle(configuratorFile string, originalFilePath string) error
	HydrationFile() []embedutils.FileReader
	VerifyRequirements() error
	DoesSatisfyConstraints(jsparse.JSDocument) bool
}

type WrapStats struct {
	WebVersion string
	Bundler    string
}

type JSWrapGroup struct {
	Wrappers []JSWebWrapper
	Stats    *WrapStats
}

func NewJSWrapGroup(wrappers []JSWebWrapper, stats *WrapStats) *JSWrapGroup {
	return &JSWrapGroup{
		Wrappers: wrappers,
		Stats:    stats,
	}
}

type JSWebWrapperList map[string]JSWebWrapper

// @@ find better name "FindAll"
// func (j JSWebWrapperList) FindAll(page jsparse.JSDocument) *JSWrapGroup {
// 	if page.Extension() == "jsx" && page.DefaultExport() != nil {
// 		if experiments.GlobalExperimentalFeatures.PreferSSR {
// 			return NewJSWrapGroup([]JSWebWrapper{j["react_ssr"], j["react_hydrate"]}, &WrapStats{
// 				WebVersion: "react",
// 				Bundler:    fmt.Sprintf("ssr + %s", j["react_hydrate"].Stats().Bundler),
// 			})
// 			// return NewJSWrapGroup([]JSWebWrapper{j["react_hydrate"]}, &WrapStats{
// 			// 	WebVersion: "react",
// 			// 	Bundler:    fmt.Sprintf("blaah"),
// 			// })
// 		} else {
// 			return NewJSWrapGroup([]JSWebWrapper{j["react_csr"]}, &WrapStats{
// 				WebVersion: "react",
// 				Bundler:    "webpack",
// 			})
// 		}
// 	}

// 	return NewJSWrapGroup([]JSWebWrapper{j["javascript"]}, &WrapStats{
// 		WebVersion: "javascript",
// 		Bundler:    "webpack",
// 	})
// }

func (j JSWebWrapperList) FindFirst(page jsparse.JSDocument) JSWebWrapper {
	for _, r := range j {
		if r.DoesSatisfyConstraints(page) {
			return r
		}
	}

	return nil
}

func NewActiveMap(bundler *BaseBundler) JSWebWrapperList {

	baseList := map[string]JSWebWrapper{
		"react": NewReactHydrate(bundler),
	}

	return baseList
}

func (l JSWebWrapperList) VerifyAll() error {
	for _, w := range l {
		err := w.VerifyRequirements()
		if err != nil {
			return err
		}
	}
	return nil
}

type BaseWebWrapper struct {
	WebDir string
}

type BundlerKey string

const (
	BundlerID BundlerKey = "bundlerID"
)

type BundlerMode string

const (
	ProductionBundle  BundlerMode = "production"
	DevelopmentBundle BundlerMode = "development"
)

type BaseBundler struct {
	Mode           BundlerMode
	WebDir         string
	PageOutputDir  string
	NodeModulesDir string
	Logger         log.Logger
}

type BundleOpts struct {
	FileName  string
	BundleKey string
	Name      string
}

type BundleConfigurator struct {
	// ConfiguratorPage represents a bundler setup file
	Page     jsparse.JSDocument
	FilePath string
}

type BundledResource struct {
	Configurators          []BundleConfigurator
	BundleOpFileDescriptor map[string]string
}

const (
	BundlerModeKey string = "bundler-mode"
)

type CacheDOMOpts struct {
	CacheDir  string
	WebPrefix string
}

//go:embed embed/*
var embedFiles embed.FS

type embedFileReader struct {
	fileName string
}

func (r *embedFileReader) Read() (fs.File, error) {
	fpath := path.Join("embed", r.fileName)

	return embedFiles.Open(fpath)
}

func (c *CacheDOMOpts) CacheWebRequest(uris []string) ([]string, error) {
	final := make([]string, len(uris))
	for i, f := range uris {
		sum := md5.Sum([]byte(f))
		hash := hex.EncodeToString(sum[:])

		extensions := strings.Split(f, ".")
		extension := extensions[len(extensions)-1]

		filepath := fmt.Sprintf("%s/%s.%s", c.CacheDir, hash, extension)

		_, err := os.Stat(filepath)

		// file path exists
		if err == nil {
			final[i] = fmt.Sprintf("%s%s.%s", c.WebPrefix, hash, extension)
			continue
		}

		// a local cached instance of the file does not exist so a request is
		// made to the endpoint, then the response is saved to a file
		if errors.Is(err, os.ErrNotExist) {
			res, err := http.Get(f)
			if err != nil {
				final[i] = uris[i]
				continue
			}

			outFile, err := os.Create(filepath)
			if err != nil {
				final[i] = uris[i]
				continue
			}

			defer outFile.Close()
			_, err = io.Copy(outFile, res.Body)
			if err != nil {
				final[i] = uris[i]
				continue
			}
		}

		final[i] = fmt.Sprintf("%s%s.%s", c.WebPrefix, hash, extension)
	}

	return final, nil
}
