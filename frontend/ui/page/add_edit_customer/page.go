package add_edit_customer

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mearaj/bhagad-house-booking/common/assets/fonts"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	. "github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	. "github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
)

type page struct {
	Manager
	Theme            *material.Theme
	navigationIcon   *widget.Icon
	buttonNavigation widget.Clickable
	CustomerForm
	initialized bool
}

func New(manager Manager, customer service.Customer) Page {
	navIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	th := manager.Theme()
	errorTh := *fonts.NewTheme()
	errorTh.ContrastBg = color.NRGBA(colornames.Red500)
	inActiveTh := *fonts.NewTheme()
	inActiveTh.ContrastBg = color.NRGBA(colornames.Grey500)
	s := page{
		Manager:        manager,
		Theme:          th,
		navigationIcon: navIcon,
		CustomerForm: NewCustomerForm(manager, customer, func(customerID string) {

		}),
	}
	return &s
}

func (p *page) Layout(gtx Gtx) Dim {
	if p.Theme == nil {
		p.Theme = fonts.NewTheme()
	}

	flex := layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}
	d := flex.Layout(gtx,
		layout.Rigid(p.DrawAppBar),
		layout.Rigid(p.CustomerForm.Layout),
	)
	return d
}
func (p *page) DrawAppBar(gtx Gtx) Dim {
	gtx.Constraints.Max.Y = gtx.Dp(56)
	th := p.Theme
	if p.buttonNavigation.Clicked() {
		p.PopUp()
	}
	return DrawAppBarLayout(gtx, th, func(gtx Gtx) Dim {
		return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx Gtx) Dim {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx Gtx) Dim {
						navigationIcon := p.navigationIcon
						button := material.IconButton(th, &p.buttonNavigation, navigationIcon, "Nav Icon Button")
						button.Size = unit.Dp(40)
						button.Background = th.Palette.ContrastBg
						button.Color = th.Palette.ContrastFg
						button.Inset = layout.UniformInset(unit.Dp(8))
						return button.Layout(gtx)
					}),
					layout.Rigid(func(gtx Gtx) Dim {
						return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx Gtx) Dim {
							titleText := fmt.Sprintf("Edit Customer %d", p.Customer.ID)
							if p.Customer.ID == 0 {
								titleText = "Add New Customer"
							}
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
func (p *page) URL() URL {
	return AddEditCustomerPageURL(p.Customer.ID)
}
