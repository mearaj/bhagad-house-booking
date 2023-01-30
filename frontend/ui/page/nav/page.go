package nav

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/page/addedit"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
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
	p := page{
		Manager:      manager,
		List:         layout.List{Axis: layout.Vertical},
		items:        newItems(manager),
		subscription: manager.Service().Subscribe(service.TopicLoggedInOut),
	}
	p.subscription.SubscribeWithCallback(p.OnServiceStateChange)
	return &p
}
func (p *page) Layout(gtx fwk.Gtx) (d fwk.Dim) {
	adminClicked := p.loginUserResponse.IsLoggedIn() && p.loginUserResponse.IsAdmin() && p.btnAddBooking.Clicked()
	if adminClicked {
		addEditBookingPage := addedit.New(p.Manager, service.Booking{})
		p.Manager.NavigateToPage(addEditBookingPage)
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
		p.Manager.NavigateToURL(fwk.NavPageURL)
	}

	return view.DrawAppBarLayout(gtx, user.Theme(), func(gtx fwk.Gtx) fwk.Dim {
		return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
						return material.ButtonLayoutStyle{
							Background:   user.Theme().ContrastBg,
							Button:       &p.buttonNavIcon,
							CornerRadius: unit.Dp(56 / 2),
						}.Layout(gtx,
							func(gtx fwk.Gtx) fwk.Dim {
								return view.DrawAppImageForNav(gtx, user.Theme())
							},
						)
					}),
					layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
						gtx.Constraints.Max.X = gtx.Constraints.Max.X - gtx.Dp(56)
						return layout.Inset{Left: unit.Dp(12)}.Layout(gtx, func(gtx fwk.Gtx) fwk.Dim {
							titleText := i18n.Get(key.Navigation)
							label := material.Label(user.Theme(), unit.Sp(18), titleText)
							label.Color = user.Theme().Palette.ContrastFg
							return component.TruncatingLabelStyle(label).Layout(gtx)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
				if !p.loginUserResponse.IsAuthorized() {
					return fwk.Dim{}
				}
				btnText := i18n.Get(key.NewBooking)
				button := material.Button(user.Theme(), &p.btnAddBooking, btnText)
				button.TextSize = unit.Sp(18)
				button.Background = user.Theme().Palette.ContrastBg
				button.Color = user.Theme().Palette.ContrastFg
				button.Inset = layout.UniformInset(unit.Dp(8))
				d := button.Layout(gtx)
				return d
			}),
		)
	})
}

func (p *page) URL() fwk.URL {
	return fwk.NavPageURL
}
func (p *page) OnServiceStateChange(event service.Event) {
	switch userResponse := event.Data.(type) {
	case service.UserResponse:
		p.loginUserResponse = userResponse
	}
}
