package template

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Tmpl struct {
	*template.Template
}

func NewTmpl() *Tmpl {
	tmpl := Tmpl{}
	tmpl.Template = template.New("")
	return &tmpl
}

func (t *Tmpl) SetFuncs(funcMap template.FuncMap) {
	t.Template = t.Template.Funcs(funcMap)
}

func (t *Tmpl) Load(dir string) (err error) {
	if dir, err = filepath.Abs(dir); err != nil {
		return err
	}

	root := t.Template

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) (_ error) {
		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".html" {
			return
		}

		var rel string
		if rel, err = filepath.Rel(dir, path); err != nil {
			return err
		}

		rel = strings.Join(strings.Split(rel, string(os.PathSeparator)), "/")
		newTmpl := root.New(rel)

		var b []byte
		if b, err = os.ReadFile(path); err != nil {
			return err
		}

		_, err = newTmpl.Parse(string(b))
		return err
	}); err != nil {
		return err
	}
	return nil
}

func (t *Tmpl) Write(dest string, name string, data interface{}) error {
	w, err := os.Create(dest)
	if err != nil {
		return err
	}

	return t.ExecuteTemplate(w, name, data)
}
