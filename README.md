vite
----

A fast (this time, actually) and minimal static site generator forked from https://git.icyphox.sh/vite to provide additional functionality.

INSTALLING

    go install github.com/toozej/go-vite@latest


USAGE

    usage: vite [options]

    A simple and minimal static site generator.

    options:
        init PATH                   create vite project at PATH
        build                       builds the current project
        new PATH                    create a new markdown post
        serve [HOST:PORT]           serves the 'build' directory


CONFIG

The configuration is unmarshalled from a config.yaml file, into the
below struct:

    type ConfigYaml struct {
        Title           string `yaml:"title"`
        Desc            string `yaml:"description"`
        DefaultTemplate string `yaml:"default-template"`
        Author          struct {
            Name  string `yaml:"name"`
            Email string `yaml:"email"`
        } `yaml:"author"`
        URL string `yaml:"url"`
    }

Example config: https://git.icyphox.sh/site/tree/config.yaml


SYNTAX HIGHLIGHTING

vite uses chroma (https://github.com/alecthomas/chroma) for syntax
highlighting. Note that CSS is not provided, and will have to be
included by the user in the templates.


TEMPLATING

Non-index templates have access to the below objects:
• Cfg: object of ConfigYaml
• Meta: map[string]string of the page's frontmatter metadata
• Body: Contains the HTML

Index templates have access to everything above, and a Posts object,
which is a slice containing HTML and Meta. This is useful for iterating
through to generate an index page.
Example: https://git.icyphox.sh/site/tree/templates/index.html

Templates are written as standard Go templates (ref:
https://godocs.io/text/template), and can be loaded recursively.
Consider the below template structure:

    templates/
    |-- blog.html
    |-- index.html
    |-- project/
        |-- index.html
        `-- project.html

The templates under project/ are referenced as project/index.html.
This deserves mention because Go templates don't recurse into
subdirectories by default (template.ParseGlob uses filepath.Glob, and
doesn't support deep-matching, i.e. **).

More templating examples can be found at:
https://git.icyphox.sh/site/tree/templates


FEEDS

Atom feeds are generated for all directories under pages/. So
pages/foo will have a Atom feed at build/foo/feed.xml.


FILE TREE

    .
    |-- build/
    |-- config.yaml
    |-- pages/
    |-- static/
    |-- templates/

The entire static/ directory gets copied over to build/, and can be
used to reference static assets -- css, images, etc. pages/ supports
only nesting one directory deep; for example: pages/blog/*.md will
render, but pages/blog/foo/*.md will not.
