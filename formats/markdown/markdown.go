package markdown

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	gotmpl "text/template"
	"time"

	"git.icyphox.sh/vite/config"
	"git.icyphox.sh/vite/template"
	"git.icyphox.sh/vite/types"
	"github.com/adrg/frontmatter"

	bf "git.icyphox.sh/grayfriday"
)

var (
	bfFlags = bf.UseXHTML | bf.Smartypants | bf.SmartypantsFractions |
		bf.SmartypantsDashes | bf.NofollowLinks | bf.FootnoteReturnLinks
	bfExts = bf.NoIntraEmphasis | bf.Tables | bf.FencedCode | bf.Autolink |
		bf.Strikethrough | bf.SpaceHeadings | bf.BackslashLineBreak |
		bf.AutoHeadingIDs | bf.HeadingIDs | bf.Footnotes | bf.NoEmptyLineBeforeBlock
)

type Markdown struct {
	body        []byte
	frontmatter map[string]string
	Path        string
}

func (*Markdown) Ext() string { return ".md" }

func (md *Markdown) Basename() string {
	return filepath.Base(md.Path)
}

// mdToHtml renders source markdown to html
func mdToHtml(source []byte) []byte {
	return bf.Run(
		source,
		bf.WithNoExtensions(),
		bf.WithRenderer(bf.NewHTMLRenderer(bf.HTMLRendererParameters{Flags: bfFlags})),
		bf.WithExtensions(bfExts),
	)
}

// template checks the frontmatter for a specified template or falls back
// to the default template -- to which it, well, templates whatever is in
// data and writes it to dest.
func (md *Markdown) template(dest, tmplDir string, data interface{}) error {
	metaTemplate := md.frontmatter["template"]
	if metaTemplate == "" {
		metaTemplate = config.Config.DefaultTemplate
	}

	tmpl := template.NewTmpl()
	tmpl.SetFuncs(gotmpl.FuncMap{
		"parsedate": func(s string) time.Time {
			date, _ := time.Parse("2006-01-02", s)
			return date
		},
	})
	if err := tmpl.Load(tmplDir); err != nil {
		return err
	}

	return tmpl.Write(dest, metaTemplate, data)
}

// extract takes the source markdown page, extracts the frontmatter
// and body. The body is converted from markdown to html here.
func (md *Markdown) extractFrontmatter(source []byte) error {
	r := bytes.NewReader(source)
	rest, err := frontmatter.Parse(r, &md.frontmatter)
	if err != nil {
		return err
	}
	md.body = mdToHtml(rest)
	return nil
}

func (md *Markdown) Frontmatter() map[string]string {
	return md.frontmatter
}

func (md *Markdown) Body() string {
	return string(md.body)
}

type templateData struct {
	Cfg   config.ConfigYaml
	Meta  map[string]string
	Body  string
	Extra interface{}
}

func (md *Markdown) Render(dest string, data interface{}) error {
	source, err := os.ReadFile(md.Path)
	if err != nil {
		return fmt.Errorf("markdown: error reading file: %w", err)
	}

	err = md.extractFrontmatter(source)
	if err != nil {
		return fmt.Errorf("markdown: error extracting frontmatter: %w", err)
	}

	if md.frontmatter["draft"] == "true" {
		fmt.Printf("vite: skipping draft %s\n", md.Path)
		return nil
	}

	err = md.template(dest, types.TemplatesDir, templateData{
		config.Config,
		md.frontmatter,
		string(md.body),
		data,
	})
	if err != nil {
		return fmt.Errorf("markdown: failed to render to destination %s: %w", dest, err)
	}
	return nil
}
