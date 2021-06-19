package main

import (
	bfc "github.com/Depado/bfchroma"
	bf "github.com/russross/blackfriday/v2"
)

var bfFlags = bf.UseXHTML | bf.Smartypants | bf.SmartypantsFractions |
	bf.SmartypantsDashes | bf.NofollowLinks
var bfExts = bf.NoIntraEmphasis | bf.Tables | bf.FencedCode | bf.Autolink |
	bf.Strikethrough | bf.SpaceHeadings | bf.BackslashLineBreak |
	bf.HeadingIDs | bf.Footnotes | bf.NoEmptyLineBeforeBlock

func mdRender(input []byte, cfg Config) []byte {
	return bf.Run(
		input,
		bf.WithRenderer(
			bfc.NewRenderer(
				bfc.ChromaStyle(Icy),
				bfc.Extend(
					bf.NewHTMLRenderer(bf.HTMLRendererParameters{
						Flags: bfFlags,
					}),
				),
			),
		),
		bf.WithExtensions(bfExts),
	)
}
