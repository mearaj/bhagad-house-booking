package view

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mearaj/bhagad-house-booking/assets"
	. "github.com/mearaj/bhagad-house-booking/ui/fwk"
)

func DrawAppIconImageCenter(gtx Gtx, theme *material.Theme) Dim {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.20)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.20)
	imgOps := paint.NewImageOp(assets.AppIconImage)
	imgWidget := widget.Image{Src: imgOps, Fit: widget.Contain, Position: layout.Center, Scale: 0}
	return imgWidget.Layout(gtx)
}

func DrawAppImageForNav(gtx Gtx, theme *material.Theme) Dim {
	gtx.Constraints.Max.X = gtx.Dp(56)
	gtx.Constraints.Max.Y = gtx.Dp(56)
	inset := layout.UniformInset(8)
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		imgOps := paint.NewImageOp(assets.AppIconImage)
		imgWidget := widget.Image{Src: imgOps, Fit: widget.Contain, Position: layout.Center, Scale: 0}
		return imgWidget.Layout(gtx)
	})
}
