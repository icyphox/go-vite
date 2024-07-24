package yaml

import (
	"fmt"
	"os"
	"path/filepath"
	gotmpl "text/template"
	"time"

	"git.icyphox.sh/vite/config"
	"git.icyphox.sh/vite/template"
	"git.icyphox.sh/vite/types"
	"gopkg.in/yaml.v3"
)

type YAML struct {
	Path string

	meta map[string]string
}

func (*YAML) Ext() string        { return ".yaml" }
func (*YAML) Body() string       { return "" }
func (y *YAML) Basename() string { return filepath.Base(y.Path) }

func (y *YAML) Frontmatter() map[string]string {
	return y.meta
}

type templateData struct {
	Cfg  config.ConfigYaml
	Meta map[string]string
	Yaml map[string]interface{}
	Body string
}

func (y *YAML) template(dest, tmplDir string, data interface{}) error {
	metaTemplate := y.meta["template"]
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

func (y *YAML) Render(dest string, data interface{}) error {
	yamlBytes, err := os.ReadFile(y.Path)
	if err != nil {
		return fmt.Errorf("yaml: failed to read file: %s: %w", y.Path, err)
	}

	yamlData := map[string]interface{}{}
	err = yaml.Unmarshal(yamlBytes, yamlData)
	if err != nil {
		return fmt.Errorf("yaml: failed to unmarshal yaml file: %s: %w", y.Path, err)
	}

	metaInterface := yamlData["meta"].(map[string]interface{})

	meta := make(map[string]string)
	for k, v := range metaInterface {
		vStr := convertToString(v)
		meta[k] = vStr
	}

	y.meta = meta

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

func convertToString(value interface{}) string {
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
