package types

const (
	BuildDir     = "build"
	PagesDir     = "pages"
	TemplatesDir = "templates"
	StaticDir    = "static"
)

type File interface {
	Ext() string
	// Render takes any arbitrary data and combines that with the global config,
	// page frontmatter and the body, as template params. Templates are read
	// from types.TemplateDir and the final html is written to dest,
	// with necessary directories being created.
	Render(dest string, data interface{}, drafts bool) error

	// Frontmatter will not be populated if Render hasn't been called.
	Frontmatter() map[string]string
	// Body will not be populated if Render hasn't been called.
	Body() string
	Basename() string
}

// Only used for building indexes and Atom feeds
type Post struct {
	Meta map[string]string
	// HTML-formatted body of post
	Body string
}
