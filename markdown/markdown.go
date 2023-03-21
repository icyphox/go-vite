package markdown

import (
	"fmt"
	"os"
	gotmpl "text/template"
	"time"

	"git.icyphox.sh/vite/config"
	"git.icyphox.sh/vite/markdown/template"

	bf "git.icyphox.sh/grayfriday"
)

var (
	bfFlags = bf.UseXHTML | bf.Smartypants | bf.SmartypantsFractions |
		bf.SmartypantsDashes | bf.NofollowLinks | bf.FootnoteReturnLinks
	bfExts = bf.NoIntraEmphasis | bf.Tables | bf.FencedCode | bf.Autolink |
		bf.Strikethrough | bf.SpaceHeadings | bf.BackslashLineBreak |
		bf.AutoHeadingIDs | bf.Footnotes | bf.NoEmptyLineBeforeBlock
)

type Output struct {
	HTML []byte
	Meta Matter
}

// Renders markdown to html, and fetches metadata.
func (out *Output) RenderMarkdown(source []byte) error {
	md := MarkdownDoc{}
	if err := md.Extract(source); err != nil {
		return fmt.Errorf("markdown: %w", err)
	}

	out.HTML = bf.Run(
		md.Body,
		bf.WithNoExtensions(),
		bf.WithRenderer(bf.NewHTMLRenderer(bf.HTMLRendererParameters{Flags: bfFlags})),
		bf.WithExtensions(bfExts),
	)
	out.Meta = md.Frontmatter
	return nil
}

// Renders out.HTML into dst html file, using the template specified
// in the frontmatter. data is the template struct.
func (out *Output) RenderHTML(dst, tmplDir string, data interface{}) error {
	metaTemplate := out.Meta["template"]
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

	w, err := os.Create(dst)
	if err != nil {
		return err
	}

	if err = tmpl.ExecuteTemplate(w, metaTemplate, data); err != nil {
		return err
	}
	return nil
}
