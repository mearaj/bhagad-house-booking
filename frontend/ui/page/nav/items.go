package nav

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
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
	settings, _ := widget.NewIcon(icons.ActionSettings)
	search, _ := widget.NewIcon(icons.ActionSearch)
	// notificationsIcon, _ := widget.NewIcon(icons.SocialNotifications)
	// helpIcon, _ := widget.NewIcon(icons.ActionHelp)
	// aboutIcon, _ := widget.NewIcon(icons.ActionInfo)
	pageItems := []fwk.View{
		&pageItem{
			Manager: manager,
			Theme:   user.Theme(),
			Title:   key.Bookings,
			Icon:    bookingsIcon,
			url:     fwk.BookingsPageURL,
		},
		&pageItem{
			Manager: manager,
			Theme:   user.Theme(),
			Title:   key.SearchBookings,
			Icon:    search,
			url:     fwk.SearchBookingsPageURL,
		},
		&pageItem{
			Manager: manager,
			Theme:   user.Theme(),
			Title:   key.Settings,
			Icon:    settings,
			url:     fwk.SettingsPageURL,
		},
		//&pageItem{
		//	Manager: manager,
		//	Theme:   user.Theme(),
		//	Title:   "Notifications",
		//	Icon:    notificationsIcon,
		//	url:     fwk.NotificationsPageURL,
		//},
		//&pageItem{
		//	Manager: manager,
		//	Theme:   user.Theme(),
		//	Title:   "Help",
		//	Icon:    helpIcon,
		//	url:     fwk.HelpPageURL,
		//},
		//&pageItem{
		//	Manager: manager,
		//	Theme:   user.Theme(),
		//	Title:   "About",
		//	Icon:    aboutIcon,
		//	url:     fwk.AboutPageURL,
		//},
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
		Theme: user.Theme(),
		Items: pageItems,
	}
	for _, i := range p.Items {
		if v, ok := i.(*pageItem); ok {
			v.parentPage = &p
		}
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
	return i.List.Layout(gtx, len(i.Items), func(gtx fwk.Gtx, index int) (d fwk.Dim) {
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
func (i *items) OnServiceStateChange(event service.Event) {
	if userResponse, ok := event.Data.(service.UserResponse); ok {
		i.loginUserResponse = userResponse
	}
}
