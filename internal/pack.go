package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/GuyARoss/orbit/pkg/bundler"
	"github.com/GuyARoss/orbit/pkg/fs"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/log"
	webwrapper "github.com/GuyARoss/orbit/pkg/web_wrapper"
	"github.com/google/uuid"
)

type PackSettings struct {
	Bundler    bundler.Bundler
	WebWrapper webwrapper.WebWrapper

	AssetDir string
	WebDir   string
}

func (s *PackSettings) CopyAssets() ([]*fs.CopyResults, error) {
	results := fs.CopyDir(s.AssetDir, s.AssetDir, ".orbit/assets", false)

	return results, nil
}

func hashKey(name string) string {
	id := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(name))

	return strings.ReplaceAll(id.String(), "-", "")
}

// PackedComponent
// Web-Component that has been succesfully ran, and output from a packing method.
type PackedComponent struct {
	PageName            string
	BundleKey           string
	OriginalFilePath    string
	PackDurationSeconds float64
}

// PackSingle
// Packs the a single file paths into the orbit root directory
// Process includes:
// - Wrapping the component with the specified front-end web framework.
// - Bundling the component with the specified javascript bundler.
func (s *PackSettings) PackSingle(pageFilePath string) (*PackedComponent, error) {
	startTime := time.Now()
	page, err := jsparse.ParsePage(pageFilePath, s.WebDir)
	if err != nil {
		return nil, err
	}

	page = s.WebWrapper.Apply(page, pageFilePath)

	bundleKey := hashKey(page.Name)
	resource, err := s.Bundler.Setup(&bundler.BundleSetupSettings{
		FileName:  pageFilePath,
		BundleKey: bundleKey,
	})
	if err != nil {
		return nil, err
	}

	bundlePageErr := page.WriteFile(resource.BundleFilePath)
	if bundlePageErr != nil {
		return nil, bundlePageErr
	}

	configErr := resource.ConfiguratorPage.WriteFile(resource.ConfiguratorFilePath)
	if configErr != nil {
		return nil, configErr
	}

	bundleErr := s.Bundler.Bundle(resource.ConfiguratorFilePath)
	return &PackedComponent{
		PageName:            page.Name,
		BundleKey:           bundleKey,
		OriginalFilePath:    pageFilePath,
		PackDurationSeconds: time.Since(startTime).Seconds(),
	}, bundleErr
}

// PackHooks
// since the implementation of the generator pattern overly complicated for
// our usecase, we instead allow the passing of "per" & "post" hooks for our
// iterative packing method "PackPages".
type PackHooks interface {
	Pre(filePath string)      // "pre" runs before each component packing iteration
	Post(elaspedTime float64) // "post" runs after each component packing iteration
}

type DefaultPackHook struct{}

func (s *DefaultPackHook) Pre(filePath string) {
	log.Info(fmt.Sprintf("bundling %s → ...", filePath))
}

func (s *DefaultPackHook) Post(elaspedTime float64) {
	log.Success(fmt.Sprintf("completed in %fs\n", elaspedTime))
}

// PackPages
// Packs the provoided file paths into the orbit root directory
// Process includes:
// - Wrapping the component with the specified front-end web framework.
// - Bundling the component with the specified javascript bundler.
func (s *PackSettings) PackMany(pages []string, hooks PackHooks) ([]*PackedComponent, error) {
	packedPages := make([]*PackedComponent, 0)
	for _, dir := range pages {
		if hooks != nil {
			hooks.Pre(dir)
		}

		// @@todo(guy): make this concurrent
		page, err := s.PackSingle(dir)

		// consider adding a flag for skipping in iteration rather than just returning
		if err != nil {
			return packedPages, err
		}
		packedPages = append(packedPages, page)
		if hooks != nil {
			hooks.Post(page.PackDurationSeconds)
		}
	}

	return packedPages, nil
}
