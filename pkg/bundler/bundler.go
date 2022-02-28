package bundler

import (
	"context"

	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/log"
)

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
	Mode BundlerMode

	WebDir         string
	PageOutputDir  string
	NodeModulesDir string
	Logger         log.Logger
}

type BundleOpts struct {
	FileName  string
	BundleKey string
}

type BundledResource struct {
	BundleFilePath       string
	ConfiguratorFilePath string
	ConfiguratorPage     jsparse.JSDocument
}

type Bundler interface {
	Setup(context.Context, *BundleOpts) (*BundledResource, error)
	Bundle(string) error
	NodeDependencies() map[string]string
}

const (
	BundlerModeKey string = "bundler-mode"
)
