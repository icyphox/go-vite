package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigYaml struct {
	Title  string `yaml:"title"`
	Desc   string `yaml:"description"`
	Author struct {
		Name  string `yaml:"name"`
		Email string `yaml:"email"`
	} `yaml:"author"`
	URL       string   `yaml:"url"`
	Prebuild  []string `yaml:"prebuild"`
	Postbuild []string `yaml:"postbuild"`
}

var Config ConfigYaml

func init() {
	cf, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Errorf("error: %+v\n", err)
		os.Exit(1)
	}
	if err = Config.ParseConfig(cf); err != nil {
		fmt.Errorf("error: %+v\n", err)
	}
}

func (c *ConfigYaml) ParseConfig(cf []byte) error {
	err := yaml.Unmarshal(cf, c)
	return err
}
