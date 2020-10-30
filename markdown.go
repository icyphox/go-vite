package main

import (
	bf "github.com/russross/blackfriday/v2"
)

var bfFlags = bf.UseXHTML | bf.Smartypants | bf.SmartypantsFractions |
	bf.SmartypantsDashes | bf.SmartypantsAngledQuotes | bf.NofollowLinks |
	bf.FootnoteReturnLinks

func mdRender(input []byte) []byte {
	return bf.Run(
		input,
		bf.WithRenderer(bf.NewHTMLRenderer(bf.HTMLRendererParameters{
			Flags: bfFlags,
		})),
	)
}
