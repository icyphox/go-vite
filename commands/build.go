package commands

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"git.icyphox.sh/vite/atom"
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

func (pgs *Pages) processFiles() error {
	for _, f := range pgs.Files {
		if filepath.Ext(f) == ".md" {
			// ex: pages/about.md
			mdFile := filepath.Join(PAGES, f)
			var htmlDir string
			// ex: build/index.html (root index)
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
					Meta markdown.Matter
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
}

func (pgs *Pages) processDirs() error {
	for _, d := range pgs.Dirs {
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
			// foo-bar.md -> foo-bar
			slug := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))

			// ex: build/blog/foo-bar/
			os.Mkdir(filepath.Join(dstDir, slug), 0755)
			// ex: build/blog/foo-bar/index.html
			htmlFile := filepath.Join(dstDir, slug, "index.html")

			if e.Name() != "_index.md" {
				ePath := filepath.Join(srcDir, e.Name())
				fb, err := os.ReadFile(ePath)
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
						Meta markdown.Matter
						Body string
					}{config.Config, out.Meta, string(out.HTML)},
				); err != nil {
					return err
				}
				posts = append(posts, out)
			}

			// Sort posts slice by date
			sort.Slice(posts, func(i, j int) bool {
				dateStr1 := posts[j].Meta["date"]
				dateStr2 := posts[i].Meta["date"]
				date1, _ := time.Parse("2006-01-02", dateStr1)
				date2, _ := time.Parse("2006-01-02", dateStr2)
				return date1.Before(date2)
			})
		}

		// Render index using posts slice.
		// ex: build/blog/index.html
		indexHTML := filepath.Join(dstDir, "index.html")
		// ex: pages/blog/_index.md
		indexMd, err := os.ReadFile(filepath.Join(srcDir, "_index.md"))
		if err != nil {
			return err
		}
		out := markdown.Output{}
		out.RenderMarkdown(indexMd)

		out.RenderHTML(indexHTML, TEMPLATES, struct {
			Cfg   config.ConfigYaml
			Meta  markdown.Matter
			Body  string
			Posts []markdown.Output
		}{config.Config, out.Meta, string(out.HTML), posts})

		// Create feeds
		// ex: build/blog/feed.xml
		xml, err := atom.NewAtomFeed(d, posts)
		if err != nil {
			return err
		}
		feedFile := filepath.Join(dstDir, "feed.xml")
		os.WriteFile(feedFile, xml, 0755)
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
	err = pages.processFiles()
	if err != nil {
		return err
	}

	// Deal with dirs -- i.e. dirs of markdown files.
	// ex: pages/{blog,travel}/*.md
	err = pages.processDirs()
	if err != nil {
		return err
	}

	return nil
}
