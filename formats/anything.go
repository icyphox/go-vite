package formats

import (
	"path/filepath"

	"git.icyphox.sh/vite/util"
)

// Anything is a stub format for unrecognized files
type Anything struct{ Path string }

func (Anything) Ext() string                    { return "" }
func (Anything) Frontmatter() map[string]string { return nil }
func (Anything) Body() string                   { return "" }
func (a Anything) Basename() string             { return filepath.Base(a.Path) }

func (a Anything) Render(dest string, data interface{}, drafts bool) error {
	return util.CopyFile(a.Path, dest)
}
