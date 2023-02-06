package nav

import (
	"gioui.org/layout"
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
	"image/color"
	"time"
)

type items struct {
	fwk.Manager
	Theme *material.Theme
	widget.Clickable
	animation view.Animation
	List      layout.List
	*view.UserForm
	Items []*pageItem
	*view.LanguageForm
}

func newItems(manager fwk.Manager) *items {
	bookingsIcon, _ := widget.NewIcon(icons.SocialGroup)
	settings, _ := widget.NewIcon(icons.ActionSettings)
	search, _ := widget.NewIcon(icons.ActionSearch)
	// notificationsIcon, _ := widget.NewIcon(icons.SocialNotifications)
	// helpIcon, _ := widget.NewIcon(icons.ActionHelp)
	// aboutIcon, _ := widget.NewIcon(icons.ActionInfo)
	pageItems := []*pageItem{
		{
			Manager: manager,
			Theme:   user.Theme(),
			Title:   key.Bookings,
			Icon:    bookingsIcon,
			url:     fwk.BookingsPageURL,
		},
		{
			Manager: manager,
			Theme:   user.Theme(),
			Title:   key.SearchBookings,
			Icon:    search,
			url:     fwk.SearchBookingsPageURL,
		},
		{
			Manager: manager,
			Theme:   user.Theme(),
			Title:   key.Settings,
			Icon:    settings,
			url:     fwk.SettingsPageURL,
		},
		//{
		//	Manager: manager,
		//	Theme:   user.Theme(),
		//	Title:   "Notifications",
		//	Icon:    notificationsIcon,
		//	url:     fwk.NotificationsPageURL,
		//},
		//{
		//	Manager: manager,
		//	Theme:   user.Theme(),
		//	Title:   "Help",
		//	Icon:    helpIcon,
		//	url:     fwk.HelpPageURL,
		//},
		//{
		//	Manager: manager,
		//	Theme:   user.Theme(),
		//	Title:   "About",
		//	Icon:    aboutIcon,
		//	url:     fwk.AboutPageURL,
		//},
	}
	p := items{
		List:    layout.List{Axis: layout.Vertical},
		Manager: manager,
		animation: view.Animation{
			Duration: time.Millisecond * 100,
			State:    component.Invisible,
			Started:  time.Time{},
		},
		Theme:        user.Theme(),
		Items:        pageItems,
		UserForm:     view.NewUserForm(manager),
		LanguageForm: view.NewLanguageForm(manager, layout.Horizontal, false),
	}
	for _, i := range p.Items {
		i.parentPage = &p
	}
	return &p
}

func (i *items) Layout(gtx fwk.Gtx) fwk.Dim {
	return i.drawItems(gtx)
}

func (i *items) drawItems(gtx fwk.Gtx) fwk.Dim {
	flex := layout.Flex{Axis: layout.Vertical}
	gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	return flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			height := gtx.Constraints.Max.Y
			d := i.List.Layout(gtx, 1, func(gtx fwk.Gtx, _ int) (d fwk.Dim) {
				inset := layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(0)}
				return inset.Layout(gtx, func(gtx fwk.Gtx) fwk.Dim {
					gtx.Constraints.Min.Y = height
					return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							flex := layout.Flex{Axis: layout.Vertical}
							itemOne := i.Items[0]
							itemTwo := i.Items[1]
							itemThree := i.Items[2]
							return flex.Layout(gtx,
								layout.Rigid(itemOne.Layout),
								layout.Rigid(view.HorizontalRule{Color: color.NRGBA(colornames.Grey500)}.Layout),
								layout.Rigid(itemTwo.Layout),
								layout.Rigid(view.HorizontalRule{Color: color.NRGBA(colornames.Grey500)}.Layout),
								layout.Rigid(itemThree.Layout),
								layout.Rigid(view.HorizontalRule{Color: color.NRGBA(colornames.Grey500)}.Layout),
								layout.Rigid(i.UserForm.Layout),
								layout.Rigid(view.HorizontalRule{Color: color.NRGBA(colornames.Grey500)}.Layout),
							)
						}),
						layout.Rigid(i.LanguageForm.Layout),
					)
				})
			})
			return d
		}),
	)
}
func (i *items) OnServiceStateChange(event service.Event) {
	if _, ok := event.Data.(service.UserResponse); ok {
		i.UserForm.OnServiceStateChange(event)
	}
}
