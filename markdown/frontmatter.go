package markdown

import (
	"bytes"

	"github.com/adrg/frontmatter"
)

type Matter map[string]string

type MarkdownDoc struct {
	Frontmatter Matter
	Body        []byte
}

func (md *MarkdownDoc) Extract(source []byte) error {
	r := bytes.NewReader(source)
	rest, err := frontmatter.Parse(r, &md.Frontmatter)
	if err != nil {
		return err
	}
	md.Body = rest
	return nil
}
