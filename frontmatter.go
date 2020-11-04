package main

import (
	"bytes"
	"github.com/adrg/frontmatter"
	"time"
)

type Date8601 struct {
	time.Time
}

func (d *Date8601) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var date string

	err := unmarshal(&date)
	if err != nil {
		return err
	}

	d.Time, err = time.Parse("2006-01-02", date)
	return err
}

type Matter struct {
	Template string   `yaml:"template"`
	URL      string   `yaml:"url"`
	Title    string   `yaml:"title"`
	Subtitle string   `yaml:"subtitle"`
	Date     Date8601 `yaml:"date"`
	Body     string
}

// Parses frontmatter, populates the `matter` struct and
// returns the rest
func parseFrontmatter(inputBytes []byte) ([]byte, Matter) {
	m := Matter{}
	input := bytes.NewReader(inputBytes)
	rest, err := frontmatter.Parse(input, &m)

	if err != nil {
		printErr(err)
	}
	return rest, m
}
