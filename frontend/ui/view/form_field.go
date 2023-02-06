package view

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"strings"
)

type FormField struct {
	FieldName      string
	LabelHintText  string
	TextField      component.TextField
	Theme          *material.Theme
	LabelBtn       *material.LabelStyle
	LabelEndWidget layout.Widget
	IconButton     *IconButton
}

func (f *FormField) Layout(gtx Gtx) Dim {
	flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceStart, Alignment: layout.Baseline}
	if f.Theme == nil {
		f.Theme = user.Theme()
	}
	th := f.Theme
	return flex.Layout(gtx,
		layout.Rigid(func(gtx Gtx) Dim {
			if f.FieldName == "" {
				return Dim{}
			}
			flex := layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}
			return flex.Layout(gtx,
				layout.Rigid(func(gtx Gtx) Dim {
					inset := layout.Inset{Bottom: 2}
					return inset.Layout(gtx, func(gtx Gtx) Dim {
						return material.Label(th, unit.Sp(16.0), f.FieldName).Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx Gtx) Dim {
					flex := layout.Flex{Alignment: layout.Middle}
					return flex.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if f.LabelBtn == nil {
								return Dim{}
							}
							inset := layout.Inset{Bottom: 2}
							return inset.Layout(gtx, func(gtx Gtx) Dim {
								return f.LabelBtn.Layout(gtx)
							})
						}),
						layout.Rigid(layout.Spacer{Width: 16}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if f.LabelEndWidget == nil {
								return Dim{}
							}
							inset := layout.Inset{Bottom: 2}
							return inset.Layout(gtx, func(gtx Gtx) Dim {
								return f.LabelEndWidget(gtx)
							})
						}),
					)
				}),
			)
		}),
		layout.Rigid(func(gtx Gtx) Dim {
			flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEnd, Alignment: layout.Start}
			inset := layout.Inset{Bottom: unit.Dp(16)}
			if f.IconButton == nil {
				inset.Bottom = unit.Dp(0)
			}
			return flex.Layout(gtx,
				layout.Rigid(func(gtx Gtx) Dim {
					return inset.Layout(gtx,
						func(gtx Gtx) Dim {
							th := *th
							origSize := th.TextSize
							if strings.TrimSpace(f.TextField.Text()) == "" && !f.TextField.Focused() {
								th.TextSize = unit.Sp(12)
							} else {
								th.TextSize = origSize
							}
							d := f.TextField.Layout(gtx, &th, f.LabelHintText)
							return d
						})
				}),
				layout.Rigid(func(gtx Gtx) Dim {
					if f.IconButton == nil {
						return Dim{}
					}
					gtx.Constraints.Min.X = 180
					return f.IconButton.Layout(gtx)
				}),
			)
		}),
	)
}
