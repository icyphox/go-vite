package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cross-cpm/go-shutil"
)

var cfg = parseConfig()

type Post struct {
	Fm Matter
}

var posts []Post

type NewFm struct {
	Template string
	URL      string
	Title    string
	Subtitle string
	Date     string
	Body     string
}

func execute(cmds []string) {
	for _, cmd := range cmds {
		out, err := exec.Command(cmd).Output()
		printMsg("running:", cmd)
		if err != nil {
			printErr(err)
			fmt.Println(string(out))
		}
	}
}

func processTemplate(tmplPath string) *template.Template {
	tmplFile := filepath.Join("templates", tmplPath)
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		printErr(err)
	}

	return tmpl
}

func handleMd(mdPath string) {
	content, err := ioutil.ReadFile(mdPath)
	if err != nil {
		printErr(err)
	}

	restContent, fm := parseFrontmatter(content)
	bodyHtml := mdRender(restContent)
	relPath, _ := filepath.Rel("pages/", mdPath)

	var buildPath string
	if strings.HasSuffix(relPath, "_index.md") {
		dir, _ := filepath.Split(relPath)
		buildPath = filepath.Join("build", dir)
	} else {
		buildPath = filepath.Join(
			"build",
			strings.TrimSuffix(relPath, filepath.Ext(relPath)),
		)
	}

	os.MkdirAll(buildPath, 0755)

	fm.Body = string(bodyHtml)

	if strings.Contains(relPath, "blog/") {
		posts = append(posts, Post{fm})
	}

	var newFm = NewFm{
		fm.Template,
		fm.URL,
		fm.Title,
		fm.Subtitle,
		fm.Date.Time.Format(cfg.DateFmt),
		fm.Body,
	}
	// combine config and matter structs
	combined := struct {
		Cfg Config
		Fm  NewFm
	}{cfg, newFm}

	htmlFile, err := os.Create(filepath.Join(buildPath, "index.html"))
	if err != nil {
		printErr(err)
		return
	}
	if fm.Template == "" {
		fm.Template = "text.html"
	}
	tmpl := processTemplate(fm.Template)
	err = tmpl.Execute(htmlFile, combined)
	if err != nil {
		printErr(err)
		return
	}
	htmlFile.Close()
}

func renderIndex(posts []Post) {
	indexTmpl := processTemplate("index.html")
	path := filepath.Join("pages", "_index.md")

	// Sort posts by date
	sort.Slice(posts, func(i, j int) bool {
		return posts[j].Fm.Date.Time.Before(posts[i].Fm.Date.Time)
	})

	content, err := ioutil.ReadFile(path)
	if err != nil {
		printErr(err)
	}

	restContent, fm := parseFrontmatter(content)
	bodyHtml := mdRender(restContent)
	fm.Body = string(bodyHtml)

	var newFm = NewFm{
		fm.Template,
		fm.URL,
		fm.Title,
		fm.Subtitle,
		fm.Date.Time.Format(cfg.DateFmt),
		fm.Body,
	}

	combined := struct {
		Fm    NewFm
		Posts []Post
		Cfg   Config
	}{newFm, posts, cfg}

	htmlFile, err := os.Create(filepath.Join("build", "index.html"))
	err = indexTmpl.Execute(htmlFile, combined)
	if err != nil {
		printErr(err)
		return
	}
	htmlFile.Close()
}

func viteBuild() {
	if len(cfg.Prebuild) != 0 {
		printMsg("executing pre-build actions...")
		execute(cfg.Prebuild)
	}
	err := filepath.Walk("./pages", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			printErr(err)
			return err
		}
		if filepath.Ext(path) == ".md" && path != filepath.Join("pages", "_index.md") {
			handleMd(path)
		} else {
			f, err := os.Stat(path)
			if err != nil {
				printErr(err)
			}
			mode := f.Mode()
			if mode.IsRegular() {
				options := shutil.CopyOptions{}
				relPath, _ := filepath.Rel("pages/", path)
				options.FollowSymlinks = true
				shutil.CopyFile(
					path,
					filepath.Join("build", relPath),
					&options,
				)
			}
		}
		return nil
	})

	if err != nil {
		printErr(err)
	}

	// Deal with the special snowflake '_index.md'
	renderIndex(posts)

	_, err = shutil.CopyTree("static", filepath.Join("build", "static"), nil)
	if err != nil {
		printErr(err)
	}
	printMsg("site build complete")
	printMsg("generating feeds...")
	generateRSS(posts, cfg)
	if len(cfg.Postbuild) != 0 {
		printMsg("executing post-build actions...")
		execute(cfg.Postbuild)
	}
}
