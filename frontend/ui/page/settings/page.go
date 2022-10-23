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
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/page/add_edit_booking"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"time"
)

type page struct {
	layout.List
	fwk.Manager
	buttonNavIcon      widget.Clickable
	btnAddBooking      widget.Clickable
	btnShowBookings    widget.Clickable
	menuIcon           *widget.Icon
	items              []*pageItem
	BookingsView       fwk.View
	menuVisibilityAnim component.VisibilityAnimation
	*view.ModalContent
}

func New(manager fwk.Manager) fwk.Page {
	menuIcon, _ := widget.NewIcon(icons.ContentAddCircle)
	bookingsIcon, _ := widget.NewIcon(icons.SocialGroup)
	contactsIcon, _ := widget.NewIcon(icons.CommunicationContacts)
	themeIcon, _ := widget.NewIcon(icons.ImagePalette)
	notificationsIcon, _ := widget.NewIcon(icons.SocialNotifications)
	helpIcon, _ := widget.NewIcon(icons.ActionHelp)
	aboutIcon, _ := widget.NewIcon(icons.ActionInfo)
	p := page{
		Manager:  manager,
		List:     layout.List{Axis: layout.Vertical},
		menuIcon: menuIcon,
		menuVisibilityAnim: component.VisibilityAnimation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		},
		items: []*pageItem{
			{
				Manager: manager,
				Theme:   manager.Theme(),
				Title:   "Bookings",
				Icon:    bookingsIcon,
				url:     fwk.BookingsPageURL,
			},
			{
				Manager: manager,
				Theme:   manager.Theme(),
				Title:   "Customers",
				Icon:    contactsIcon,
				url:     fwk.CustomersPageURL,
			},
			{
				Manager: manager,
				Theme:   manager.Theme(),
				Title:   "Theme",
				Icon:    themeIcon,
				url:     fwk.ThemePageURL,
			},
			{
				Manager: manager,
				Theme:   manager.Theme(),
				Title:   "Notifications",
				Icon:    notificationsIcon,
				url:     fwk.NotificationsPageURL,
			},
			{
				Manager: manager,
				Theme:   manager.Theme(),
				Title:   "Help",
				Icon:    helpIcon,
				url:     fwk.HelpPageURL,
			},
			{
				Manager: manager,
				Theme:   manager.Theme(),
				Title:   "About",
				Icon:    aboutIcon,
				url:     fwk.AboutPageURL,
			},
		},
	}
	p.ModalContent = view.NewModalContent(func() { p.Modal().Dismiss(nil) })
	return &p
}
func (p *page) Layout(gtx fwk.Gtx) (d fwk.Dim) {
	if p.items == nil {
		p.items = make([]*pageItem, 0)
	}
	if p.btnAddBooking.Clicked() {
		addEditBookingPage := add_edit_booking.New(p.Manager, sqlc.Booking{})
		p.Manager.NavigateToPage(addEditBookingPage, func() {

		})
	}

	if p.btnShowBookings.Clicked() {
		p.menuVisibilityAnim.Disappear(gtx.Now)
		p.Modal().Show(p.drawShowBookingsModal, nil, fwk.Animation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		})
	}

	flex := layout.Flex{Axis: layout.Vertical,
		Spacing:   layout.SpaceEnd,
		Alignment: layout.Start,
	}
	d = flex.Layout(gtx,
		layout.Rigid(p.DrawAppBar),
		layout.Rigid(p.drawItems),
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
				var img image.Image
				//var err error
				//a := p.Service().Booking()
				//if a.PublicKey != "" && len(a.PublicImage) != 0 {
				//	img, _, err = image.Decode(bytes.NewReader(a.PublicImage))
				//	if err != nil {
				//		alog.Logger().Error(err)
				//	}
				//}
				if img != nil {
					return p.btnShowBookings.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						radii := gtx.Dp(20)
						gtx.Constraints.Max.X, gtx.Constraints.Max.Y = radii*2, radii*2
						bounds := image.Rect(0, 0, radii*2, radii*2)
						clipOp := clip.UniformRRect(bounds, radii).Push(gtx.Ops)
						imgOps := paint.NewImageOp(img)
						imgWidget := widget.Image{Src: imgOps, Fit: widget.Contain, Position: layout.Center, Scale: 0}
						d := imgWidget.Layout(gtx)
						clipOp.Pop()
						return d
					})
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
func (p *page) drawItems(gtx fwk.Gtx) fwk.Dim {
	return p.List.Layout(gtx, len(p.items), func(gtx fwk.Gtx, index int) (d fwk.Dim) {
		inset := layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(0)}
		return inset.Layout(gtx, func(gtx fwk.Gtx) fwk.Dim {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(p.items[index].Layout),
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
func (p *page) onAddBookingSuccess() {
	p.Modal().Dismiss(func() {
		p.NavigateToUrl(fwk.SettingsPageURL, nil)
	})
}

func (p *page) drawShowBookingsModal(gtx fwk.Gtx) fwk.Dim {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.85)
	return p.ModalContent.DrawContent(gtx, p.Theme(), p.BookingsView.Layout)
}

func (p *page) onBookingChange() {
	p.Modal().Dismiss(p.afterBookingsModalDismissed)
}
func (p *page) afterBookingsModalDismissed() {
	p.NavigateToUrl(fwk.SettingsPageURL, func() {
		a := p.Service().Booking()
		txt := fmt.Sprintf("Switched to %d booking", a.ID)
		p.Snackbar().Show(txt, nil, color.NRGBA{}, "")
	})
}

func (p *page) URL() fwk.URL {
	return fwk.SettingsPageURL
}
