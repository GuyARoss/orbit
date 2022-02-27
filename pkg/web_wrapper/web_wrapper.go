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

func (c *CacheDOMOpts) CacheWebRequest(uris []string) []string {
	final := make([]string, 0)
	for _, f := range uris {
		sum := md5.Sum([]byte(f))
		hash := hex.EncodeToString(sum[:])

		extensions := strings.Split(f, ".")
		extension := extensions[len(extensions)-1]

		filepath := fmt.Sprintf("%s/%s.%s", c.CacheDir, hash, extension)
		if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
			res, err := http.Get(f)
			fmt.Println(err)

			// @@todo; return error stating that required wrap data cannot be found
			outFile, err := os.Create(filepath)
			fmt.Println(err)

			defer outFile.Close()
			_, err = io.Copy(outFile, res.Body)
			fmt.Println(err)

			final = append(final, fmt.Sprintf("%s%s", c.WebPrefix, hash))
		}
	}

	return final
}

type JSWebWrapper interface {
	Apply(page jsparse.JSDocument, toFilePath string) jsparse.JSDocument
	NodeDependencies() map[string]string
	DoesSatisfyConstraints(fileExtension string) bool
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
