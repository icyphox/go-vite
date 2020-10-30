package main

import (
	"bytes"
	"github.com/adrg/frontmatter"
)

type Matter struct {
	Template string `yaml:"template"`
	URL      string `yaml:"url"`
	Title    string `yaml:"title"`
	Subtitle string `yaml:"subtitle"`
	Date     string `yaml:"date"`
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
