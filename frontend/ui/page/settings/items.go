package settings

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"time"
)

type items struct {
	fwk.Manager
	Theme *material.Theme
	widget.Clickable
	animation         view.Animation
	List              layout.List
	Items             []fwk.View
	loginUserResponse service.UserResponse
	subscription      service.Subscriber
}

func newItems(manager fwk.Manager) *items {
	bookingsIcon, _ := widget.NewIcon(icons.SocialGroup)
	themeIcon, _ := widget.NewIcon(icons.ImagePalette)
	notificationsIcon, _ := widget.NewIcon(icons.SocialNotifications)
	helpIcon, _ := widget.NewIcon(icons.ActionHelp)
	aboutIcon, _ := widget.NewIcon(icons.ActionInfo)
	pageItems := []fwk.View{
		&pageItem{
			Manager: manager,
			Theme:   manager.Theme(),
			Title:   "Bookings",
			Icon:    bookingsIcon,
			url:     fwk.BookingsPageURL,
		},
		&pageItem{
			Manager: manager,
			Theme:   manager.Theme(),
			Title:   "Theme",
			Icon:    themeIcon,
			url:     fwk.ThemePageURL,
		},
		&pageItem{
			Manager: manager,
			Theme:   manager.Theme(),
			Title:   "Notifications",
			Icon:    notificationsIcon,
			url:     fwk.NotificationsPageURL,
		},
		&pageItem{
			Manager: manager,
			Theme:   manager.Theme(),
			Title:   "Help",
			Icon:    helpIcon,
			url:     fwk.HelpPageURL,
		},
		&pageItem{
			Manager: manager,
			Theme:   manager.Theme(),
			Title:   "About",
			Icon:    aboutIcon,
			url:     fwk.AboutPageURL,
		},
		view.NewUserForm(manager),
	}
	p := items{
		List:         layout.List{Axis: layout.Vertical},
		Manager:      manager,
		subscription: manager.Service().Subscribe(service.TopicUserLoggedInOut),
		animation: view.Animation{
			Duration: time.Millisecond * 100,
			State:    component.Invisible,
			Started:  time.Time{},
		},
		Theme: manager.Theme(),
		Items: pageItems,
	}
	p.subscription.SubscribeWithCallback(p.OnServiceStateChange)
	return &p
}

func (i *items) Layout(gtx fwk.Gtx) fwk.Dim {
	flex := layout.Flex{Axis: layout.Vertical}
	return flex.Layout(gtx,
		layout.Rigid(i.drawItems),
	)
}

func (i *items) drawItems(gtx fwk.Gtx) fwk.Dim {
	isAuthorized := i.loginUserResponse.IsLoggedIn() && i.loginUserResponse.IsAdmin()
	return i.List.Layout(gtx, len(i.Items), func(gtx fwk.Gtx, index int) (d fwk.Dim) {
		switch i := i.Items[index].(type) {
		case *pageItem:
			switch i.Title {
			case "Customers":
				if !isAuthorized {
					return view.Dim{}
				}
			}
		}
		inset := layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(0)}
		return inset.Layout(gtx, func(gtx fwk.Gtx) fwk.Dim {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(i.Items[index].Layout),
				layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
					size := image.Pt(gtx.Constraints.Max.X, gtx.Dp(1))
					bounds := image.Rectangle{Max: size}
					bgColor := color.NRGBA(colornames.Grey500)
					bgColor.A = 75
					paint.FillShape(gtx.Ops, bgColor, clip.UniformRRect(bounds, 0).Op(gtx.Ops))
					return fwk.Dim{Size: image.Pt(size.X, size.Y)}
				}),
			)
		})
	})
}
func (p *items) OnServiceStateChange(event service.Event) {
	switch userResponse := event.Data.(type) {
	case service.UserResponse:
		p.loginUserResponse = userResponse
	}
}
