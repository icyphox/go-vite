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
	DefaultTemplate string `yaml:"-"`
	Author          struct {
		Name  string `yaml:"name"`
		Email string `yaml:"email"`
	} `yaml:"author"`
	URL       string   `yaml:"url"`
	PreBuild  []string `yaml:"preBuild"`
	PostBuild []string `yaml:"postBuild"`
}

// For backward compat with `default-template`
func (c *ConfigYaml) UnmarshalYAML(value *yaml.Node) error {
	type Alias ConfigYaml // Create an alias to avoid recursion

	var aux Alias

	if err := value.Decode(&aux); err != nil {
		return err
	}

	// Handle the DefaultTemplate field
	var temp struct {
		DefaultTemplate1 string `yaml:"default-template"`
		DefaultTemplate2 string `yaml:"defaultTemplate"`
	}
	if err := value.Decode(&temp); err != nil {
		return err
	}

	if temp.DefaultTemplate1 != "" {
		aux.DefaultTemplate = temp.DefaultTemplate1
	} else {
		aux.DefaultTemplate = temp.DefaultTemplate2
	}

	*c = ConfigYaml(aux) // Assign the unmarshalled values back to the original struct

	return nil
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
