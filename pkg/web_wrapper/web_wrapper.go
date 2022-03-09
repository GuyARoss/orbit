// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// LICENSE file in the root directory of this source tree.
package webwrapper

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/GuyARoss/orbit/pkg/jsparse"
)

type BaseWebWrapper struct {
	WebDir string
}

type CacheDOMOpts struct {
	CacheDir  string
	WebPrefix string
}

func (c *CacheDOMOpts) CacheWebRequest(uris []string) ([]string, error) {
	final := make([]string, 0)
	for _, f := range uris {
		sum := md5.Sum([]byte(f))
		hash := hex.EncodeToString(sum[:])

		extensions := strings.Split(f, ".")
		extension := extensions[len(extensions)-1]

		filepath := fmt.Sprintf("%s/%s.%s", c.CacheDir, hash, extension)

		// a local cached instance of the file does not exist so a request is
		// made to the endpoint, then the response is saved to a file
		if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
			res, err := http.Get(f)
			if err != nil {
				return final, err
			}

			outFile, err := os.Create(filepath)
			if err != nil {
				return final, err
			}

			defer outFile.Close()
			_, err = io.Copy(outFile, res.Body)
			if err != nil {
				return final, err
			}
		}

		final = append(final, fmt.Sprintf("%s%s.%s", c.WebPrefix, hash, extension))
	}

	return final, nil
}

type JSWebWrapper interface {
	Apply(jsparse.JSDocument, string) jsparse.JSDocument
	NodeDependencies() map[string]string
	DoesSatisfyConstraints(string) bool
	Version() string
	RequiredBodyDOMElements(context.Context, *CacheDOMOpts) []string
}

type JSWebWrapperMap []JSWebWrapper

func NewActiveMap() JSWebWrapperMap {
	return []JSWebWrapper{
		&ReactWebWrapper{},
	}
}

func (j *JSWebWrapperMap) FirstMatch(fileExtension string) JSWebWrapper {
	for _, f := range *j {
		if f.DoesSatisfyConstraints(fileExtension) {
			return f
		}
	}

	return nil
}
