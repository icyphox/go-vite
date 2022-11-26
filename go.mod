module git.icyphox.sh/vite

go 1.15

replace github.com/russross/blackfriday/v2 => git.icyphox.sh/grayfriday v0.0.0-20221126034429-23c704183914

// replace github.com/russross/blackfriday/v2 => ../grayfriday

require (
	github.com/Depado/bfchroma v1.3.0
	github.com/adrg/frontmatter v0.2.0
	github.com/alecthomas/chroma v0.10.0
	github.com/russross/blackfriday/v2 v2.0.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
