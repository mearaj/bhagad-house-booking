package fonts

import (
	_ "embed"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/mearaj/bhagad-house-booking/frontend/assets/fonts/noto_sans_gujarati"
)

func NewTheme() *material.Theme {
	var collection []text.FontFace
	collection = append(collection, noto_sans_gujarati.Collection...)
	th := material.NewTheme(collection)
	return th
}
