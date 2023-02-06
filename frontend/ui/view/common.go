package view

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"image"
	"image/color"
	"strings"
)

func DrawFormField(gtx Gtx, th *material.Theme, labelText, labelHintText string, textField *component.TextField, button *IconButton, labelBtn *material.ButtonStyle, labelEndWidget layout.Widget) Dim {
	flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceStart, Alignment: layout.Baseline}
	return flex.Layout(gtx,
		layout.Rigid(func(gtx Gtx) Dim {
			if labelText == "" {
				return Dim{}
			}
			flex := layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}
			return flex.Layout(gtx,
				layout.Rigid(func(gtx Gtx) Dim {
					inset := layout.Inset{Bottom: 2}
					return inset.Layout(gtx, func(gtx Gtx) Dim {
						return material.Label(th, unit.Sp(16.0), labelText).Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx Gtx) Dim {
					flex := layout.Flex{Alignment: layout.Middle}
					return flex.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if labelBtn == nil {
								return Dim{}
							}
							inset := layout.Inset{Bottom: 2}
							return inset.Layout(gtx, func(gtx Gtx) Dim {
								return labelBtn.Layout(gtx)
							})
						}),
						layout.Rigid(layout.Spacer{Width: 16}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if labelEndWidget == nil {
								return Dim{}
							}
							inset := layout.Inset{Bottom: 2}
							return inset.Layout(gtx, func(gtx Gtx) Dim {
								return labelEndWidget(gtx)
							})
						}),
					)
				}),
			)
		}),
		layout.Rigid(func(gtx Gtx) Dim {
			flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEnd, Alignment: layout.Start}
			inset := layout.Inset{Bottom: unit.Dp(16)}
			if button == nil {
				inset.Bottom = unit.Dp(0)
			}
			return flex.Layout(gtx,
				layout.Rigid(func(gtx Gtx) Dim {
					return inset.Layout(gtx,
						func(gtx Gtx) Dim {
							th := *th
							origSize := th.TextSize
							if strings.TrimSpace(textField.Text()) == "" && !textField.Focused() {
								th.TextSize = unit.Sp(12)
							} else {
								th.TextSize = origSize
							}
							d := textField.Layout(gtx, &th, labelHintText)
							return d
						})
				}),
				layout.Rigid(func(gtx Gtx) Dim {
					if button == nil {
						return Dim{}
					}
					gtx.Constraints.Min.X = 180
					return button.Layout(gtx)
				}),
			)
		}),
	)
}

func DrawAvatar(gtx Gtx, initials string, bgColor color.NRGBA, textTheme *material.Theme) Dim {
	d := component.Rect{
		Color: bgColor,
		Size:  image.Point{X: gtx.Dp(48), Y: gtx.Dp(48)},
		Radii: gtx.Dp(48) / 2,
	}.Layout(gtx)
	macro2 := op.Record(gtx.Ops)
	d2 := material.Label(textTheme, unit.Sp(20), initials).Layout(gtx)
	macro2.Stop()
	op.Offset(image.Point{
		X: d.Size.X - d2.Size.X/2,
		Y: d.Size.Y - d2.Size.Y/2,
	}).Add(gtx.Ops)
	material.Label(textTheme, unit.Sp(20), initials).Layout(gtx)
	return d
}
