// style.go: generate css
package style

import (
	"os"
	"path/filepath"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
)

var syntax = chroma.MustNewStyle("syntax", chroma.StyleEntries{
	chroma.CommentMultiline:     "italic #999988",
	chroma.CommentPreproc:       "bold #999999",
	chroma.CommentSingle:        "italic #999988",
	chroma.CommentSpecial:       "bold italic #999999",
	chroma.Comment:              "italic #999988",
	chroma.Error:                "bg:#e3d2d2 #a61717",
	chroma.GenericDeleted:       "bg:#ffdddd #000000",
	chroma.GenericEmph:          "italic #000000",
	chroma.GenericError:         "#aa0000",
	chroma.GenericHeading:       "#999999",
	chroma.GenericInserted:      "bg:#ddffdd #000000",
	chroma.GenericOutput:        "#888888",
	chroma.GenericPrompt:        "#555555",
	chroma.GenericStrong:        "bold",
	chroma.GenericSubheading:    "#aaaaaa",
	chroma.GenericTraceback:     "#aa0000",
	chroma.GenericUnderline:     "underline",
	chroma.KeywordType:          "bold #222222",
	chroma.Keyword:              "bold #000000",
	chroma.LiteralNumber:        "#009999",
	chroma.LiteralStringRegex:   "#009926",
	chroma.LiteralStringSymbol:  "#990073",
	chroma.LiteralString:        "#509c93",
	chroma.NameAttribute:        "#008080",
	chroma.NameBuiltinPseudo:    "#999999",
	chroma.NameBuiltin:          "#509c93",
	chroma.NameClass:            "bold #666666",
	chroma.NameConstant:         "#008080",
	chroma.NameDecorator:        "bold #3c5d5d",
	chroma.NameEntity:           "#509c93",
	chroma.NameException:        "bold #444444",
	chroma.NameFunction:         "bold #444444",
	chroma.NameLabel:            "bold #444444",
	chroma.NameNamespace:        "#555555",
	chroma.NameTag:              "#000080",
	chroma.NameVariableClass:    "#008080",
	chroma.NameVariableGlobal:   "#008080",
	chroma.NameVariableInstance: "#008080",
	chroma.NameVariable:         "#008080",
	chroma.Operator:             "bold #000000",
	chroma.TextWhitespace:       "#bbbbbb",
	chroma.Background:           " bg:#ffffff",
})

func GenerateStyleFiles() error {
	// generate syntax.css
	formatter := html.New(html.WithClasses(true))
	syntaxStyleFilePath := filepath.Join("static", "syntax.css")
	syntaxStyleFile, err := os.Create(syntaxStyleFilePath)
	if err != nil {
		return err
	}
	defer syntaxStyleFile.Close()
	if err = formatter.WriteCSS(syntaxStyleFile, syntax); err != nil {
		return err
	}

	return nil
}
