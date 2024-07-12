package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"git.icyphox.sh/vite/atom"
	"git.icyphox.sh/vite/config"
	"git.icyphox.sh/vite/markdown"
	"git.icyphox.sh/vite/util"
	"gopkg.in/yaml.v3"
)

const (
	BuildDir     = "build"
	PagesDir     = "pages"
	TemplatesDir = "templates"
	StaticDir    = "static"
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
		switch filepath.Ext(f) {
		case ".md":
			// ex: pages/about.md
			mdFile := filepath.Join(PagesDir, f)
			var htmlDir string
			// ex: build/index.html (root index)
			if f == "_index.md" {
				htmlDir = BuildDir
			} else {
				htmlDir = filepath.Join(
					BuildDir,
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
			if err = out.RenderMarkdown(fb); err != nil {
				return err
			}
			if err = out.RenderHTML(
				htmlFile,
				TemplatesDir,
				struct {
					Cfg  config.ConfigYaml
					Meta markdown.Matter
					Body string
				}{config.Config, out.Meta, string(out.HTML)},
			); err != nil {
				return err
			}
		case ".yaml":
			// ex: pages/reading.yaml
			yamlFile := filepath.Join(PagesDir, f)
			htmlDir := filepath.Join(BuildDir, strings.TrimSuffix(f, ".yaml"))
			os.Mkdir(htmlDir, 0755)
			htmlFile := filepath.Join(htmlDir, "index.html")

			yb, err := os.ReadFile(yamlFile)
			if err != nil {
				return err
			}

			data := map[string]interface{}{}
			err = yaml.Unmarshal(yb, &data)
			if err != nil {
				return fmt.Errorf("error: unmarshalling yaml file %s: %v", yamlFile, err)
			}

			meta := make(map[string]string)
			for k, v := range data["meta"].(map[string]interface{}) {
				meta[k] = v.(string)
			}

			out := markdown.Output{}
			out.Meta = meta
			if err = out.RenderHTML(
				htmlFile,
				TemplatesDir,
				struct {
					Cfg  config.ConfigYaml
					Meta markdown.Matter
					Yaml map[string]interface{}
					Body string
				}{config.Config, meta, data, ""},
			); err != nil {
				return err
			}
		default:
			src := filepath.Join(PagesDir, f)
			dst := filepath.Join(BuildDir, f)
			if err := util.CopyFile(src, dst); err != nil {
				return err
			}
		}
	}
	return nil
}

func (pgs *Pages) processDirs() error {
	for _, d := range pgs.Dirs {
		// ex: build/blog
		dstDir := filepath.Join(BuildDir, d)
		// ex: pages/blog
		srcDir := filepath.Join(PagesDir, d)
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
				if err := out.RenderMarkdown(fb); err != nil {
					return err
				}
				if err = out.RenderHTML(
					htmlFile,
					TemplatesDir,
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
		if err := out.RenderMarkdown(indexMd); err != nil {
			return err
		}

		out.RenderHTML(indexHTML, TemplatesDir, struct {
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
	if err := preBuild(); err != nil {
		return err
	}
	fmt.Print("vite: building... ")
	pages := Pages{}
	if err := pages.initPages(); err != nil {
		return err
	}

	// Clean the build directory.
	if err := util.Clean(BuildDir); err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	wgDone := make(chan bool)

	ec := make(chan error)

	// Deal with files.
	// ex: pages/{_index,about,etc}.md
	go func() {
		err := pages.processFiles()
		if err != nil {
			ec <- err
		}
		wg.Done()
	}()

	// Deal with dirs -- i.e. dirs of markdown files.
	// ex: pages/{blog,travel}/*.md
	go func() {
		err := pages.processDirs()
		if err != nil {
			ec <- err
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(wgDone)
	}()

	select {
	case <-wgDone:
		break
	case err := <-ec:
		close(ec)
		return err
	}

	// Copy the static directory into build
	// ex: build/static/
	buildStatic := filepath.Join(BuildDir, StaticDir)
	os.Mkdir(buildStatic, 0755)
	if err := util.CopyDir(StaticDir, buildStatic); err != nil {
		return err
	}
	fmt.Print("done\n")

	if err := postBuild(); err != nil {
		return err
	}
	return nil
}

func postBuild() error {
	for _, cmd := range config.Config.PostBuild {
		fmt.Println("vite: running post-build command:", cmd)
		if err := util.RunCmd(cmd); err != nil {
			return err
		}
	}
	return nil
}

func preBuild() error {
	for _, cmd := range config.Config.PreBuild {
		fmt.Println("vite: running pre-build command:", cmd)
		if err := util.RunCmd(cmd); err != nil {
			return err
		}
	}
	return nil
}
