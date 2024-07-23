package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"git.icyphox.sh/vite/atom"
	"git.icyphox.sh/vite/config"
	"git.icyphox.sh/vite/formats"
	"git.icyphox.sh/vite/formats/markdown"
	"git.icyphox.sh/vite/formats/yaml"
	"git.icyphox.sh/vite/types"
	"git.icyphox.sh/vite/util"
)

type Dir struct {
	Name     string
	HasIndex bool
	Files    []types.File
}

type Pages struct {
	Dirs  []Dir
	Files []types.File
}

func NewPages() (*Pages, error) {
	pages := &Pages{}

	entries, err := os.ReadDir(types.PagesDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			thingsDir := filepath.Join(types.PagesDir, entry.Name())
			dir := Dir{Name: entry.Name()}
			things, err := os.ReadDir(thingsDir)
			if err != nil {
				return nil, err
			}

			for _, thing := range things {
				if thing.Name() == "_index.md" {
					dir.HasIndex = true
					continue
				}
				switch filepath.Ext(thing.Name()) {
				case ".md":
					path := filepath.Join(thingsDir, thing.Name())
					dir.Files = append(dir.Files, &markdown.Markdown{Path: path})
				case ".yaml":
					path := filepath.Join(thingsDir, thing.Name())
					dir.Files = append(dir.Files, &yaml.YAML{Path: path})
				default:
					fmt.Printf("warn: unrecognized filetype for file: %s\n", thing.Name())
				}
			}

			pages.Dirs = append(pages.Dirs, dir)
		} else {
			path := filepath.Join(types.PagesDir, entry.Name())
			switch filepath.Ext(entry.Name()) {
			case ".md":
				pages.Files = append(pages.Files, &markdown.Markdown{Path: path})
			case ".yaml":
				pages.Files = append(pages.Files, &yaml.YAML{Path: path})
			default:
				pages.Files = append(pages.Files, formats.Anything{Path: path})
			}
		}
	}

	return pages, nil
}

// Build is the core builder function. Converts markdown/yaml
// to html, copies over non-.md/.yaml files, etc.
func Build() error {
	if err := preBuild(); err != nil {
		return err
	}
	fmt.Println("vite: building")

	pages, err := NewPages()
	if err != nil {
		return fmt.Errorf("error: reading 'pages/' %w", err)
	}

	if err := util.Clean(types.BuildDir); err != nil {
		return err
	}

	if err := pages.ProcessFiles(); err != nil {
		return err
	}

	if err := pages.ProcessDirectories(); err != nil {
		return err
	}

	buildStatic := filepath.Join(types.BuildDir, types.StaticDir)
	if err := os.MkdirAll(buildStatic, 0755); err != nil {
		return err
	}
	if err := util.CopyDir(types.StaticDir, buildStatic); err != nil {
		return err
	}
	fmt.Println("done")

	return nil
}

// ProcessFiles handles root level files under 'pages',
// for example: 'pages/_index.md' or 'pages/about.md'.
func (p *Pages) ProcessFiles() error {
	for _, f := range p.Files {
		var htmlDir string
		if f.Basename() == "_index.md" {
			htmlDir = types.BuildDir
		} else {
			htmlDir = filepath.Join(types.BuildDir, strings.TrimSuffix(f.Basename(), f.Ext()))
		}

		destFile := filepath.Join(htmlDir, "index.html")
		if f.Ext() == "" {
			destFile = filepath.Join(types.BuildDir, f.Basename())
		} else {
			if err := os.MkdirAll(htmlDir, 0755); err != nil {
				return err
			}
		}
		if err := f.Render(destFile, nil); err != nil {
			return fmt.Errorf("error: failed to render %s: %w", destFile, err)
		}
	}
	return nil
}

// ProcessDirectories handles directories of posts under 'pages',
// for example: 'pages/photos/foo.md' or 'pages/blog/bar.md'.
func (p *Pages) ProcessDirectories() error {
	for _, dir := range p.Dirs {
		dstDir := filepath.Join(types.BuildDir, dir.Name)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			return fmt.Errorf("error: failed to create directory: %s: %w", dstDir, err)
		}

		posts := []types.Post{}

		for _, file := range dir.Files {
			post := types.Post{}
			// foo-bar.md -> foo-bar
			slug := strings.TrimSuffix(file.Basename(), file.Ext())
			dstFile := filepath.Join(dstDir, slug, "index.html")

			// ex: build/blog/foo-bar/
			if err := os.MkdirAll(filepath.Join(dstDir, slug), 0755); err != nil {
				return fmt.Errorf("error: failed to create directory: %s: %w", dstDir, err)
			}

			if err := file.Render(dstFile, nil); err != nil {
				return fmt.Errorf("error: failed to render %s: %w", dstFile, err)
			}

			post.Meta = file.Frontmatter()
			post.Body = file.Body()
			posts = append(posts, post)
		}

		sort.Slice(posts, func(i, j int) bool {
			dateStr1 := posts[j].Meta["date"]
			dateStr2 := posts[i].Meta["date"]
			date1, _ := time.Parse("2006-01-02", dateStr1)
			date2, _ := time.Parse("2006-01-02", dateStr2)
			return date1.Before(date2)
		})

		if dir.HasIndex {
			indexMd := filepath.Join(types.PagesDir, dir.Name, "_index.md")
			index := markdown.Markdown{Path: indexMd}
			dstFile := filepath.Join(dstDir, "index.html")
			if err := index.Render(dstFile, posts); err != nil {
				return fmt.Errorf("error: failed to render index %s: %w", dstFile, err)
			}
		}

		xml, err := atom.NewAtomFeed(filepath.Join(types.PagesDir, dir.Name), posts)
		if err != nil {
			return fmt.Errorf("error: failed to create atom feed for: %s: %w", dir.Name, err)
		}
		feedFile := filepath.Join(dstDir, "feed.xml")
		os.WriteFile(feedFile, xml, 0755)
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
