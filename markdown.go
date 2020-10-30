package main

import (
	bfc "github.com/Depado/bfchroma"
	bf "github.com/russross/blackfriday/v2"
)

var bfFlags = bf.UseXHTML | bf.Smartypants | bf.SmartypantsFractions |
	bf.SmartypantsDashes | bf.NofollowLinks | bf.FootnoteReturnLinks

func mdRender(input []byte) []byte {
	return bf.Run(
		input,
		bf.WithRenderer(
			bfc.NewRenderer(
				bfc.Style("bw"),
				bfc.Extend(
					bf.NewHTMLRenderer(bf.HTMLRendererParameters{
						Flags: bfFlags,
					}),
				),
			),
		),
	)
}
