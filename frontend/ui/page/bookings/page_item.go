package bookings

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/page/add_edit_booking"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	color2 "image/color"
	"time"
)

type pageItem struct {
	*material.Theme
	btnRow    widget.Clickable
	btnDelete widget.Clickable
	btnYes    widget.Clickable
	btnNo     widget.Clickable
	fwk.Manager
	time.Time
	sqlc.Booking
	LoginUserResponse service.UserResponse
	ModalContent      *view.ModalContent
	parentPage        *page
}

func (i *pageItem) Layout(gtx fwk.Gtx) fwk.Dim {
	if i.Theme == nil {
		i.Theme = user.Theme()
	}
	if i.ModalContent == nil {
		i.ModalContent = view.NewModalContent(func() {
			i.Modal().Dismiss(nil)
		})
	}
	return i.layoutContent(gtx)
}
func (i *pageItem) layoutContent(gtx fwk.Gtx) fwk.Dim {
	inset := layout.UniformInset(16)
	isAuthorized := i.LoginUserResponse.IsLoggedIn() && i.LoginUserResponse.IsAdmin()
	if isAuthorized && i.btnDelete.Clicked() {
		i.btnRow.Clicked() // discard row click
		if i.Booking.ID != 0 {
			i.Modal().Show(i.drawDeleteBookingsModal, nil, fwk.Animation{
				Duration: time.Millisecond * 250,
				State:    component.Invisible,
				Started:  time.Time{},
			})
		}
	}
	adminClicked := isAuthorized && i.btnRow.Clicked()
	if adminClicked {
		i.Manager.NavigateToPage(add_edit_booking.New(i.Manager, i.Booking))
		i.Window().Invalidate()
	}
	d := i.btnRow.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		d := inset.Layout(gtx, func(gtx fwk.Gtx) fwk.Dim {
			flex := layout.Flex{Spacing: layout.SpaceEnd, Alignment: layout.Middle}
			bookingTime := i.Time
			layoutRatio := flexUserLayoutRatio
			if isAuthorized {
				layoutRatio = flexAdminLayoutRatio
			}
			d := flex.Layout(gtx,
				layout.Flexed(layoutRatio[0], func(gtx fwk.Gtx) fwk.Dim {
					bookingID := fmt.Sprintf("%d", i.Booking.ID)
					if bookingID == "0" {
						bookingID = "N/A"
					}
					b := material.Body1(i.Theme, bookingID)
					b.Font.Weight = text.Normal
					return b.Layout(gtx)
				}),
				layout.Flexed(layoutRatio[1], func(gtx fwk.Gtx) fwk.Dim {
					var bookingDate string
					switch bookingTime.Day() {
					case 1, 21, 31:
						bookingDate = fmt.Sprintf("%dst", bookingTime.Day())
					case 2, 22:
						bookingDate = fmt.Sprintf("%dnd", bookingTime.Day())
					case 3, 23:
						bookingDate = fmt.Sprintf("%drd", bookingTime.Day())
					default:
						bookingDate = fmt.Sprintf("%dth", bookingTime.Day())
					}
					b := material.Body1(i.Theme, fmt.Sprintf("%s", bookingDate))
					b.Font.Weight = text.Medium
					return b.Layout(gtx)
				}),
				layout.Flexed(layoutRatio[2], func(gtx fwk.Gtx) fwk.Dim {
					b := material.Body1(i.Theme, fmt.Sprintf("%s", bookingTime.Weekday()))
					b.Font.Weight = text.Normal
					return b.Layout(gtx)
				}),
				layout.Flexed(layoutRatio[3], func(gtx fwk.Gtx) fwk.Dim {
					icon, _ := widget.NewIcon(icons.NotificationEventBusy)
					color := color2.NRGBA(colornames.Red500)
					isAvailable := i.Booking.ID == 0
					if isAvailable {
						icon, _ = widget.NewIcon(icons.NotificationEventAvailable)
						color = color2.NRGBA(colornames.Green500)
					}
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.X = gtx.Dp(24)
						gtx.Constraints.Max.Y = gtx.Dp(24)
						gtx.Constraints.Min.X = gtx.Dp(24)
						gtx.Constraints.Min.Y = gtx.Dp(24)
						return icon.Layout(gtx, color)
					})
				}),
				layout.Flexed(layoutRatio[4], func(gtx fwk.Gtx) fwk.Dim {
					if !isAuthorized {
						return fwk.Dim{}
					}
					icon, _ := widget.NewIcon(icons.ActionDelete)
					color := color2.NRGBA(colornames.Red500)
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.X = gtx.Dp(24)
						gtx.Constraints.Max.Y = gtx.Dp(24)
						gtx.Constraints.Min.X = gtx.Dp(24)
						gtx.Constraints.Min.Y = gtx.Dp(24)
						if i.Booking.ID == 0 {
							icon, _ = widget.NewIcon(icons.ContentBlock)
							color = i.Theme.ContrastBg
						}
						return i.btnDelete.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return icon.Layout(gtx, color)
						})
					})
				}),
			)
			return d
		})
		return d
	})
	gtx.Constraints.Max.Y = d.Size.Y
	return d
}
func (p *pageItem) drawDeleteBookingsModal(gtx fwk.Gtx) fwk.Dim {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.85)
	isAuthorized := p.LoginUserResponse.IsLoggedIn() && p.LoginUserResponse.IsAdmin()
	if p.btnYes.Clicked() && isAuthorized {
		p.parentPage.isDeletingBooking = true
		p.Service().DeleteBooking(p.Booking.ID)
		p.Modal().Dismiss(func() { p.Window().Invalidate() })
	}
	if p.btnNo.Clicked() {
		p.Modal().Dismiss(func() {})
	}

	delPrompt := i18n.Get(key.BookingDeletePrompt)
	bookingID := fmt.Sprintf("%s %s %d", i18n.Get(key.Booking), i18n.Get(key.ID), p.Booking.ID)
	startDate := p.Booking.StartDate
	endDate := p.Booking.EndDate
	startDateStr := i18n.Get(key.StartDate)
	endDateStr := i18n.Get(key.EndDate)
	promptContent := view.NewPromptContent(p.Theme,
		i18n.Get(key.BookingDeletion),
		fmt.Sprintf(
			"%s\n %s\n %s:- %d %s %d.\n %s:- %d %s %d.\n",
			delPrompt,
			bookingID,
			startDateStr,
			startDate.Day(),
			startDate.Month(),
			startDate.Year(),
			endDateStr,
			endDate.Day(),
			endDate.Month().String(),
			endDate.Year(),
		),
		&p.btnYes,
		&p.btnNo)
	return p.ModalContent.DrawContent(gtx, p.Theme, promptContent.Layout)
}
