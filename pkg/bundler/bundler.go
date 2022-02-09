package bundler

import "github.com/GuyARoss/orbit/pkg/jsparse"

type BundlerMode string

const (
	ProductionBundle  BundlerMode = "production"
	DevelopmentBundle BundlerMode = "development"
)

type BundleSettings struct {
	Mode BundlerMode

	WebDir        string
	PageOutputDir string
}

type BundleSetupSettings struct {
	FileName  string
	BundleKey string
}

type BundledResource struct {
	BundleFilePath       string
	ConfiguratorFilePath string
	ConfiguratorPage     jsparse.JSDocument
}

type Bundler interface {
	Setup(settings *BundleSetupSettings) (*BundledResource, error)
	Bundle(configuratorFilePath string) error
}
