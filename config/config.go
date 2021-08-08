package config

import (
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
	URL string `yaml:"url"`
	//	Prebuild  []string `yaml:"prebuild"`
	//	Postbuild []string `yaml:"postbuild"`
}

func (c *ConfigYaml) ParseConfig() error {
	cf, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(cf, c); err != nil {
		return err
	}
	return nil
}
