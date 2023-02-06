package settings

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	. "github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"strconv"
)

type page struct {
	layout.List
	Manager
	title         string
	btnNavIcon    widget.Clickable
	navIcon       *widget.Icon
	Theme         *material.Theme
	rateForm      component.TextField
	prevRateValue string
	btnSave       view.IconButton
}

func New(manager Manager) Page {
	navIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	pg := page{
		Manager: manager,
		navIcon: navIcon,
		List:    layout.List{Axis: layout.Vertical},
		Theme:   user.Theme(),
	}
	pg.title = i18n.GetFromCode(key.Theme, *user.LanguageCode())
	pg.btnSave = view.IconButton{Theme: user.Theme()}
	pg.rateForm.SetText(fmt.Sprintf("%.2f", user.BookingRate()))
	return &pg
}

func (p *page) Layout(gtx Gtx) Dim {
	flex := layout.Flex{Axis: layout.Vertical,
		Spacing:   layout.SpaceEnd,
		Alignment: layout.Start,
	}
	p.title = i18n.GetFromCode(key.Settings, *user.LanguageCode())

	d := flex.Layout(gtx,
		layout.Rigid(p.DrawAppBar),
		layout.Rigid(p.drawContentLayout),
	)
	return d
}

func (p *page) DrawAppBar(gtx Gtx) Dim {
	if p.btnNavIcon.Clicked() {
		p.PopUp()
	}
	if p.Theme == nil {
		p.Theme = user.Theme()
	}
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	gtx.Constraints.Max.Y = gtx.Dp(56)
	th := p.Theme
	return view.DrawAppBarLayout(gtx, th, func(gtx Gtx) Dim {
		return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx Gtx) Dim {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx Gtx) Dim {
						navigationIcon := p.navIcon
						button := material.IconButton(th, &p.btnNavIcon, navigationIcon, "Nav Icon Button")
						button.Size = unit.Dp(40)
						button.Background = th.Palette.ContrastBg
						button.Color = th.Palette.ContrastFg
						button.Inset = layout.UniformInset(unit.Dp(8))
						return button.Layout(gtx)
					}),
					layout.Rigid(func(gtx Gtx) Dim {
						return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx Gtx) Dim {
							titleText := p.title
							title := material.Body1(th, titleText)
							title.Color = th.Palette.ContrastFg
							title.TextSize = unit.Sp(18)
							return title.Layout(gtx)
						})
					}),
				)
			}),
		)
	})

}

func (p *page) drawContentLayout(gtx Gtx) Dim {
	th := p.Theme
	paint.FillShape(gtx.Ops, th.Bg,
		clip.UniformRRect(image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y), 0).Op(gtx.Ops))
	inset := layout.UniformInset(unit.Dp(16))
	gtx.Constraints.Min = gtx.Constraints.Max
	if p.prevRateValue != p.rateForm.Text() {
		p.prevRateValue = p.rateForm.Text()
		p.rateForm.ClearError()
	}
	if p.btnSave.Button.Clicked() {
		rate := p.rateForm.Text()
		rateFloat, err := strconv.ParseFloat(rate, 64)
		if err != nil {
			p.rateForm.SetError(err.Error())
		}
		if err == nil {
			user.SetBookingRate(rateFloat)
			p.Snackbar().Show("Settings applied successfully", nil, color.NRGBA{}, "CLOSE")
		}
	}

	return inset.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			p.List.Axis = layout.Vertical
			return p.List.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
				flex := layout.Flex{Axis: layout.Vertical}
				return flex.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						labelText := i18n.Get(key.RatePerDay)
						p.btnSave.Text = i18n.Get(key.Apply)
						return view.DrawFormField(gtx,
							th, labelText, labelText, &p.rateForm, &p.btnSave, nil, nil)
					}),
				)
			})
		},
	)
}

func (p *page) URL() URL {
	return SettingsPageURL
}
