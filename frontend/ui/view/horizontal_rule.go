package view

import (
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"image"
	"image/color"
)

type HorizontalRule struct {
	Color color.NRGBA
}

func (h HorizontalRule) Layout(gtx Gtx) Dim {
	size := image.Pt(gtx.Constraints.Max.X, gtx.Dp(1))
	bounds := image.Rectangle{Max: size}
	bgColor := h.Color
	bgColor.A = 75
	paint.FillShape(gtx.Ops, bgColor, clip.UniformRRect(bounds, 0).Op(gtx.Ops))
	return fwk.Dim{Size: image.Pt(size.X, size.Y)}
}
