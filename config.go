package godid

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/go-yaml/yaml"
	homedir "github.com/mitchellh/go-homedir"
)

const (
	configPath = "~/.godid/config.yml"
)

type config struct {
	StorePath string `yaml:"store_path"`
}

func (c *config) GetStorePath() (string, error) {
	return homedir.Expand(c.StorePath)
}

var (
	defaultConfig = config{
		StorePath: "~/.godid/store.db",
	}
)

func getConfig() (*config, error) {
	path, err := homedir.Expand(configPath)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return initDefaultConfig()
	}
	cfgBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg config
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func initDefaultConfig() (*config, error) {
	if os.Getenv("GODID_TEST") != "" {
		cfgPath, err := homedir.Expand(configPath)
		if err != nil {
			return nil, err
		}
		pathDir := path.Dir(cfgPath)
		_, err = os.Stat(pathDir)
		if os.IsNotExist(err) {
			if err := os.MkdirAll(pathDir, os.ModePerm); err != nil { // nolint: vetshadow
				return nil, err
			}
		}
		cfgBytes, err := yaml.Marshal(defaultConfig)
		if err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(cfgPath, cfgBytes, 0600); err != nil {
			return nil, err
		}
	}
	return &defaultConfig, nil
}
