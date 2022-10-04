package experiments

import "github.com/GuyARoss/orbit/pkg/log"

type Features struct {
	PreferSSR         bool
	PreferSWCCompiler bool
}

var GlobalExperimentalFeatures *Features = &Features{}

func LoadSingleton(logger log.Logger, features []string) error {
	for _, e := range features {
		switch e {
		case "ssr":
			GlobalExperimentalFeatures.PreferSSR = true
			logger.Warn("experimental feature 'prefer ssr' enabled\n")
		case "swc":
			GlobalExperimentalFeatures.PreferSWCCompiler = true
			logger.Warn("experimental feature 'prefer swc compiler' enabled\n")
		}
	}

	return nil
}
