package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var Config ConfigYaml

func init() {
	err := Config.parseConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: config: %+v\n", err)
		os.Exit(1)
	}
}

type ConfigYaml struct {
	Title           string `yaml:"title"`
	Desc            string `yaml:"description"`
	DefaultTemplate string `yaml:"default-template"`
	Author          struct {
		Name  string `yaml:"name"`
		Email string `yaml:"email"`
	} `yaml:"author"`
	URL string `yaml:"url"`
	//	Prebuild  []string `yaml:"prebuild"`
	//	Postbuild []string `yaml:"postbuild"`
}

func (c *ConfigYaml) parseConfig() error {
	cf, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(cf, c); err != nil {
		return err
	}
	return nil
}
