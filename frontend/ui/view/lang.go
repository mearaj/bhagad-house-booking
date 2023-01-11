package view

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mearaj/bhagad-house-booking/frontend/assets/fonts"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/code"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"image"
)

type LanguageForm struct {
	langEnum     widget.Enum
	languageCode code.Code
	showTitle    bool
	layout.Axis
	fwk.Manager
}

func NewLanguageForm(m fwk.Manager, axis layout.Axis, showTitle bool) *LanguageForm {
	langForm := LanguageForm{Manager: m, Axis: axis, showTitle: showTitle}
	langForm.langEnum.Value = string(*user.LanguageCode())
	langForm.languageCode = *user.LanguageCode()
	return &langForm
}

func (p *LanguageForm) Layout(gtx fwk.Gtx) fwk.Dim {
	th := fonts.NewTheme()
	paint.FillShape(gtx.Ops, th.Bg,
		clip.UniformRRect(image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y), 0).Op(gtx.Ops))
	inset := layout.UniformInset(unit.Dp(16))
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	if p.langEnum.Changed() {
		p.languageCode = code.Code(p.langEnum.Value)
		*user.LanguageCode() = p.languageCode
		user.SaveSettings()
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		flex := layout.Flex{Axis: layout.Vertical}
		return flex.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if !p.showTitle {
					return fwk.Dim{}
				}
				h := material.H4(th, i18n.GetFromCode(key.Language, p.languageCode))
				h.Alignment = text.Middle
				return h.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if p.showTitle {
					return layout.Spacer{Height: unit.Dp(16)}.Layout(gtx)
				}
				return fwk.Dim{}
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				flex := layout.Flex{Axis: p.Axis}
				return flex.Layout(gtx,
					layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
						inset := layout.Inset{Right: unit.Dp(16)}
						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return material.RadioButton(
								th,
								&p.langEnum,
								string(code.English),
								i18n.GetFromCode(key.English, p.languageCode),
							).Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
						inset := layout.Inset{Right: unit.Dp(16)}
						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return material.RadioButton(
								th,
								&p.langEnum,
								string(code.Gujarati),
								i18n.GetFromCode(key.Gujarati, p.languageCode),
							).Layout(gtx)
						})
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
				)
			}),
		)
	})
}
