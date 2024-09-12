package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/fsutils"
	"github.com/GuyARoss/orbit/pkg/htmlparse"
	"github.com/GuyARoss/orbit/pkg/webwrap"
	"github.com/spf13/viper"
)

// SPABuildOpts options used to build a SPA (single-page-application)
type SPABuildOpts struct {
	// PublicHTMLPath is the html path that if set gets used as a base for the
	// output of the web application.
	PublicHTMLPath string
	// SpaOutDir is the directory that the spa will get built to
	SpaOutDir string
}

func BuildSPA(component srcpack.PackComponent, opts *SPABuildOpts) error {
	// this bundle component should get copied from the .orbit/dist to the spa_output
	bundlePath := fmt.Sprintf(".orbit/dist/%s.js", component.BundleKey())
	outDir := viper.GetString("spa_out_dir")

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		if err := os.Mkdir(outDir, os.ModePerm); err != nil {
			return fmt.Errorf("cannot setup spa out directory %s, exiting", err)
		}
	}

	if err := fsutils.CopyFile(bundlePath, fmt.Sprintf("%s/%s.js", outDir, component.BundleKey())); err != nil {
		return fmt.Errorf("cannot copy file %s, exiting", err)
	}

	// @@ add case where public path is not set
	if viper.GetString("public_path") != "" {
		// parse the html page (if exists) and add the javascript to it.
		htmlDoc := htmlparse.DocFromFile(viper.GetString("public_path"))
		wr := component.WebWrapper()
		if wr == nil {
			// TODO(guy): should likely throw an error
			return nil
		}

		body := wr.RequiredBodyDOMElements(context.TODO(), &webwrap.CacheDOMOpts{
			WebPrefix: outDir,
			CacheDir:  "",
		})
		// note: altering the order of the appends will break functionality
		htmlDoc.Body = append(htmlDoc.Body, wr.DocumentTag(component.BundleKey()))

		htmlDoc.Body = append(htmlDoc.Body, body...)
		htmlDoc.Body = append(htmlDoc.Body, fmt.Sprintf(`<script src="./%s.js"></script>`, component.BundleKey()))

		htmlDoc.SaveToFile(fmt.Sprintf("%s/index.html", outDir))
	}

	return nil
}
