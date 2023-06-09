package cmd

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/lwydyby/mr-chglog/config"
)

// ConfigLoader ...
type ConfigLoader interface {
	Load(string) (*config.MRChLogConfig, error)
}

type configLoaderImpl struct {
}

// NewConfigLoader ...
func NewConfigLoader() ConfigLoader {
	return &configLoaderImpl{}
}

func (loader *configLoaderImpl) Load(path string) (*config.MRChLogConfig, error) {
	fp := filepath.Clean(path)
	bytes, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	c := &config.MRChLogConfig{}
	err = yaml.Unmarshal(bytes, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
