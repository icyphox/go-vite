package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"git.icyphox.sh/vite/config"
	"git.icyphox.sh/vite/markdown"
	"git.icyphox.sh/vite/util"
)

const (
	BUILD     = "build"
	PAGES     = "pages"
	TEMPLATES = "templates"
)

type Pages struct {
	Dirs  []string
	Files []string
}

// Populates a Pages object with dirs and files
// found in 'pages/'.
func (pgs *Pages) initPages() error {
	files, err := os.ReadDir("./pages")
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			pgs.Dirs = append(pgs.Dirs, f.Name())
		} else {
			pgs.Files = append(pgs.Files, f.Name())
		}
	}

	return nil
}

// Core builder function. Converts markdown to html,
// copies over non .md files, etc.
func Build() error {
	pages := Pages{}
	err := pages.initPages()
	if err != nil {
		return err
	}

	// Deal with files.
	// ex: pages/{_index,about,etc}.md
	err = func() error {
		for _, f := range pages.Files {
			if filepath.Ext(f) == ".md" {
				// ex: pages/about.md
				mdFile := filepath.Join(PAGES, f)
				var htmlDir string
				if f == "_index.md" {
					htmlDir = BUILD
				} else {
					htmlDir = filepath.Join(
						BUILD,
						strings.TrimSuffix(f, ".md"),
					)
				}
				os.Mkdir(htmlDir, 0755)
				// ex: build/about/index.html
				htmlFile := filepath.Join(htmlDir, "index.html")

				fb, err := os.ReadFile(mdFile)
				if err != nil {
					return err
				}

				out := markdown.Output{}
				out.RenderMarkdown(fb)
				if err = out.RenderHTML(
					htmlFile,
					TEMPLATES,
					struct {
						Cfg  config.ConfigYaml
						Fm   markdown.Matter
						Body string
					}{config.Config, out.Meta, string(out.HTML)},
				); err != nil {
					return err
				}
			} else {
				src := filepath.Join(PAGES, f)
				util.CopyFile(src, BUILD)
			}
		}
		return nil
	}()
	if err != nil {
		return err
	}

	// Deal with dirs -- i.e. of markdown files.
	// ex: pages/{blog,travel}/*.md
	err = func() error {
		for _, d := range pages.Dirs {
			// ex: build/blog
			dstDir := filepath.Join(BUILD, d)
			// ex: pages/blog
			srcDir := filepath.Join(PAGES, d)
			os.Mkdir(dstDir, 0755)

			entries, err := os.ReadDir(srcDir)
			if err != nil {
				return err
			}

			posts := []markdown.Output{}
			// Collect all posts
			for _, e := range entries {
				ePath := filepath.Join(srcDir, e.Name())
				fb, err := os.ReadFile(ePath)
				if err != nil {
					return err
				}

				out := markdown.Output{}
				out.RenderMarkdown(fb)

				slug := strings.TrimSuffix(filepath.Ext(e.Name()), ".md")

				htmlFile := filepath.Join(dstDir, slug)
				out.RenderHTML(
					htmlFile,
					TEMPLATES,
					struct {
						Cfg  config.ConfigYaml
						Fm   markdown.Matter
						Body string
					}{config.Config, out.Meta, string(out.HTML)})
				posts = append(posts, out)
			}

			// Sort posts slice by date
			sort.Slice(posts, func(i, j int) bool {
				date1 := posts[j].Meta["date"].(time.Time)
				date2 := posts[i].Meta["date"].(time.Time)
				return date1.Before(date2)
			})

			for _, p := range posts {
				fmt.Println(p.Meta["date"])
			}
		}
		return nil
	}()
	if err != nil {
		return err
	}

	return nil
}
