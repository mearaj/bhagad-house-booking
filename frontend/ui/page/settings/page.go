package settings

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/page/add_edit_booking"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"golang.org/x/exp/shiny/materialdesign/icons"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type page struct {
	layout.List
	items fwk.View
	fwk.Manager
	buttonNavIcon     widget.Clickable
	btnAddBooking     widget.Clickable
	menuIcon          *widget.Icon
	BookingsView      fwk.View
	loginUserResponse service.UserResponse
	subscription      service.Subscriber
}

func New(manager fwk.Manager) fwk.Page {
	menuIcon, _ := widget.NewIcon(icons.ContentAddCircle)
	p := page{
		Manager:      manager,
		menuIcon:     menuIcon,
		List:         layout.List{Axis: layout.Vertical},
		items:        newItems(manager),
		subscription: manager.Service().Subscribe(service.TopicUserLoggedInOut),
	}
	p.subscription.SubscribeWithCallback(p.OnServiceStateChange)
	return &p
}
func (p *page) Layout(gtx fwk.Gtx) (d fwk.Dim) {
	adminClicked := p.loginUserResponse.IsLoggedIn() && p.loginUserResponse.IsAdmin() && p.btnAddBooking.Clicked()
	if adminClicked {
		addEditBookingPage := add_edit_booking.New(p.Manager, sqlc.Booking{})
		p.Manager.NavigateToPage(addEditBookingPage, func() {})
	}

	flex := layout.Flex{Axis: layout.Vertical,
		Spacing:   layout.SpaceEnd,
		Alignment: layout.Start,
	}
	d = flex.Layout(gtx,
		layout.Rigid(p.DrawAppBar),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.items.Layout(gtx)
		}),
	)
	return d
}
func (p *page) DrawAppBar(gtx fwk.Gtx) fwk.Dim {
	if p.buttonNavIcon.Clicked() {
		p.Manager.NavigateToUrl(fwk.SettingsPageURL, nil)
	}

	return view.DrawAppBarLayout(gtx, p.Manager.Theme(), func(gtx fwk.Gtx) fwk.Dim {
		return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
						return material.ButtonLayoutStyle{
							Background:   p.Manager.Theme().ContrastBg,
							Button:       &p.buttonNavIcon,
							CornerRadius: unit.Dp(56 / 2),
						}.Layout(gtx,
							func(gtx fwk.Gtx) fwk.Dim {
								return view.DrawAppImageForNav(gtx, p.Manager.Theme())
							},
						)
					}),
					layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
						gtx.Constraints.Max.X = gtx.Constraints.Max.X - gtx.Dp(56)
						return layout.Inset{Left: unit.Dp(12)}.Layout(gtx, func(gtx fwk.Gtx) fwk.Dim {
							titleText := "Settings"
							label := material.Label(p.Manager.Theme(), unit.Sp(18), titleText)
							label.Color = p.Manager.Theme().Palette.ContrastFg
							return component.TruncatingLabelStyle(label).Layout(gtx)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
				isAuthorized := p.loginUserResponse.IsLoggedIn() && p.loginUserResponse.IsAdmin()
				if !isAuthorized {
					return fwk.Dim{}
				}
				button := material.IconButton(p.Manager.Theme(), &p.btnAddBooking, p.menuIcon, "Context Menu")
				button.Size = unit.Dp(40)
				button.Background = p.Manager.Theme().Palette.ContrastBg
				button.Color = p.Manager.Theme().Palette.ContrastFg
				button.Inset = layout.UniformInset(unit.Dp(8))
				d := button.Layout(gtx)
				return d
			}),
		)
	})
}

func (p *page) onAddBookingSuccess() {
	p.Modal().Dismiss(func() {
		p.NavigateToUrl(fwk.SettingsPageURL, nil)
	})
}

func (p *page) URL() fwk.URL {
	return fwk.SettingsPageURL
}
func (p *page) OnServiceStateChange(event service.Event) {
	switch userResponse := event.Data.(type) {
	case service.UserResponse:
		p.loginUserResponse = userResponse
	}
}
