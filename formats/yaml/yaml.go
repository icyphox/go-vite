package yaml

import (
	"fmt"
	"os"
	"path/filepath"
	gotmpl "text/template"
	"time"

	"gopkg.in/yaml.v3"
	"tangled.sh/icyphox.sh/vite/config"
	"tangled.sh/icyphox.sh/vite/template"
	"tangled.sh/icyphox.sh/vite/types"
)

type YAML struct {
	Path string

	meta map[string]any
}

func (*YAML) Ext() string        { return ".yaml" }
func (*YAML) Body() string       { return "" }
func (y *YAML) Basename() string { return filepath.Base(y.Path) }

func (y *YAML) Frontmatter() map[string]any {
	return y.meta
}

type templateData struct {
	Cfg  config.ConfigYaml
	Meta map[string]any
	Yaml map[string]any
	Body string
}

func (y *YAML) template(dest, tmplDir string, data any) error {
	var metaTemplate string
	if templateVal, ok := y.meta["template"]; ok {
		if strVal, isStr := templateVal.(string); isStr {
			metaTemplate = strVal
		}
	}
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

func (y *YAML) Render(dest string, data any, drafts bool) error {
	yamlBytes, err := os.ReadFile(y.Path)
	if err != nil {
		return fmt.Errorf("yaml: failed to read file: %s: %w", y.Path, err)
	}

	yamlData := map[string]any{}
	err = yaml.Unmarshal(yamlBytes, yamlData)
	if err != nil {
		return fmt.Errorf("yaml: failed to unmarshal yaml file: %s: %w", y.Path, err)
	}

	metaInterface, ok := yamlData["meta"].(map[string]any)
	if !ok {
		return fmt.Errorf("yaml: meta section is not a map: %s", y.Path)
	}

	y.meta = metaInterface

	err = y.template(dest, types.TemplatesDir, templateData{
		config.Config,
		y.meta,
		yamlData,
		"",
	})
	if err != nil {
		return fmt.Errorf("yaml: failed to render to destination %s: %w", dest, err)
	}

	return nil
}

func convertToString(value any) string {
	// Infer type and convert to string
	switch v := value.(type) {
	case string:
		return v
	case time.Time:
		return v.Format("2006-01-02")
	default:
		return fmt.Sprintf("%v", v)
	}
}
