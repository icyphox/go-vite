package markdown

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Matter map[string]string

type MarkdownDoc struct {
	Frontmatter Matter
	Body        []byte
}

const (
	yamlDelim = "---"
)

func (md *MarkdownDoc) Extract(source []byte) error {
	bufsize := 1024 * 1024
	buf := make([]byte, bufsize)

	input := bytes.NewReader(source)
	s := bufio.NewScanner(input)
	s.Buffer(buf, bufsize)

	matter := []byte{}
	body := []byte{}

	s.Split(splitFunc)
	n := 0
	for s.Scan() {
		if n == 0 {
			matter = s.Bytes()
		} else if n == 1 {
			body = s.Bytes()
		}
		n++
	}
	if err := s.Err(); err != nil {
		return fmt.Errorf("error: failed to scan text")
	}
	if err := yaml.Unmarshal(matter, &md.Frontmatter); err != nil {
		return fmt.Errorf("error: failed to parse yaml")
	}
	md.Body = body
	return nil
}

func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	delim, err := sniffDelim(data)
	if err != nil {
		return 0, nil, err
	}
	if delim != yamlDelim {
		return 0, nil, fmt.Errorf("error: %s is not a supported delimiter", delim)
	}
	if x := bytes.Index(data, []byte(delim)); x >= 0 {
		if next := bytes.Index(data[x+len(delim):], []byte(delim)); next > 0 {
			return next + len(delim), bytes.TrimSpace(data[:next+len(delim)]), nil
		}
		return len(data), bytes.TrimSpace(data[x+len(delim):]), nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func sniffDelim(input []byte) (string, error) {
	if len(input) < 4 {
		return "", errors.New("error: input is empty")
	}
	return string(input[:3]), nil
}
