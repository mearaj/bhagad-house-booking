package noto_sans_gujarati

import (
	_ "embed"
	"gioui.org/font/opentype"
	"gioui.org/text"
)

//go:embed NotoSansGujarati-VariableFont_wdth,wght.ttf

var notoSansGujarati []byte

var gujaratiFace, _ = opentype.Parse(notoSansGujarati)

var gujaratiFont = text.Font{Weight: text.Normal, Style: text.Regular, Typeface: "Gujarati"}

var Collection = []text.FontFace{
	{Font: gujaratiFont, Face: gujaratiFace},
}
