package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Title       string            `yaml:"title"`
	Header      string            `yaml:"header"`
	DateFmt     string			  `yaml:datefmt`
	SiteURL     string            `yaml:"siteurl"`
	Description string            `yaml:"description"`
	Author      map[string]string `yaml:"author"`
	Footer      string            `yaml:"footer"`
	Prebuild    []string          `yaml:"prebuild"`
	Postbuild   []string          `yaml:"postbuild"`
	RSSPrefixURL string			  `yaml:"rssprefixurl"`
}

func parseConfig() Config {
	var config Config
	cf, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		printErr(err)
	}

	err = yaml.Unmarshal(cf, &config)
	if err != nil {
		printErr(err)
	}

	return config
}
