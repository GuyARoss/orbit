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
	"github.com/GuyARoss/orbit/pkg/experiments"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/log"
)

type JSWebWrapper interface {
	RequiredBodyDOMElements(context.Context, *CacheDOMOpts) []string
	Setup(context.Context, *BundleOpts) (*BundledResource, error)
	Apply(jsparse.JSDocument) (jsparse.JSDocument, error)
	DoesSatisfyConstraints(string) bool
	Version() string
	Stats() *WrapStats
	Bundle(configuratorFile string, originalFilePath string) error
	HydrationFile() []embedutils.FileReader
	VerifyRequirements() error
}

type WrapStats struct {
	WebVersion string
	Bundler    string
}

type JSWebWrapperList []JSWebWrapper

// FirstMatch finds the first js web wrapper in the currently list that satisfies the file extension constraints
func (j *JSWebWrapperList) FirstMatch(fileExtension string) JSWebWrapper {
	for _, f := range *j {
		if f.DoesSatisfyConstraints(fileExtension) {
			return f
		}
	}

	return nil
}

func NewActiveMap(bundler *BaseBundler) JSWebWrapperList {
	sourcedoc := jsparse.NewEmptyDocument()
	initdoc := jsparse.NewEmptyDocument()

	baseList := []JSWebWrapper{
		&JavascriptWrapper{
			BaseBundler: bundler,
		},
	}

	if experiments.GlobalExperimentalFeatures.PreferSSR {
		baseList = append(baseList, NewReactSSR(&NewReactSSROpts{
			Bundler:      bundler,
			SourceMapDoc: sourcedoc,
			InitDoc:      initdoc,
		}))
	} else {
		baseList = append(baseList, NewReactWebWrap(bundler))
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
	Configurators  []BundleConfigurator
	BundleFilePath string
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
