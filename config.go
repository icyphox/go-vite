package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Title     string
	Header    string
	Footer    string
	Prebuild  []string
	Postbuild []string
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
