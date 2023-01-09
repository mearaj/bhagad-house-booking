package search_bookings

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
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
	"image/color"
	"time"
)

type pageItem struct {
	*material.Theme
	btnRow    widget.Clickable
	btnDelete widget.Clickable
	btnYes    widget.Clickable
	btnNo     widget.Clickable
	fwk.Manager
	sqlc.Booking
	LoginUserResponse service.UserResponse
	ModalContent      *view.ModalContent
	parentPage        *page
}

const headFontSize = unit.Sp(18)
const bodyFontSize = unit.Sp(18)
const headFontWeight = text.ExtraBold
const bodyFontWeight = text.Medium
const labelWidth = unit.Dp(160)

func (p *pageItem) Layout(gtx fwk.Gtx) fwk.Dim {
	if p.Theme == nil {
		p.Theme = user.Theme()
	}
	if p.ModalContent == nil {
		p.ModalContent = view.NewModalContent(func() {
			p.Modal().Dismiss(nil)
		})
	}
	return p.layoutContent(gtx)
}
func (p *pageItem) layoutContent(gtx fwk.Gtx) fwk.Dim {
	inset := layout.UniformInset(16)
	isAuthorized := p.LoginUserResponse.IsLoggedIn() && p.LoginUserResponse.IsAdmin()
	if isAuthorized && p.btnDelete.Clicked() {
		p.btnRow.Clicked() // discard row click
		if p.Booking.ID != 0 {
			p.Modal().Show(p.drawDeleteBookingsModal, nil, fwk.Animation{
				Duration: time.Millisecond * 250,
				State:    component.Invisible,
				Started:  time.Time{},
			})
		}
	}
	adminClicked := isAuthorized && p.btnRow.Clicked()
	if adminClicked {
		p.Manager.NavigateToPage(add_edit_booking.New(p.Manager, p.Booking))
		p.Window().Invalidate()
	}
	d := p.btnRow.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		d := inset.Layout(gtx, func(gtx fwk.Gtx) fwk.Dim {
			flex := layout.Flex{Axis: layout.Vertical}
			d := flex.Layout(gtx,
				layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
					flex := layout.Flex{Alignment: layout.Middle}
					return flex.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Max.X = gtx.Constraints.Max.X - 80
							labelTxt := fmt.Sprintf("%s %s", i18n.Get(key.Booking), i18n.Get(key.ID))
							valueTxt := fmt.Sprintf("%d", p.Booking.ID)
							return p.drawBookingField(gtx, labelTxt, valueTxt)
						}),
						layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
							gtx.Constraints.Max.X = 80
							gtx.Constraints.Min.X = 80
							delLabel := i18n.Get(key.Delete)
							btn := material.Button(p.Theme, &p.btnDelete, delLabel)
							btn.Background = color.NRGBA(colornames.Red500)
							return btn.Layout(gtx)
						}),
					)
				}),
				layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
					labelTxt := i18n.Get(key.CustomerName)
					valueTxt := p.Booking.CustomerName
					return p.drawBookingField(gtx, labelTxt, valueTxt)
				}),
				layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
					labelTxt := i18n.Get(key.Booking) + " " + i18n.Get(key.Details)
					valueTxt := p.Booking.Details
					return p.drawBookingField(gtx, labelTxt, valueTxt)
				}),
				layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
					label := i18n.Get(key.StartDate)
					return p.drawBookingDate(gtx, p.StartDate, label)
				}),
				layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
					label := i18n.Get(key.EndDate)
					return p.drawBookingDate(gtx, p.EndDate, label)
				}),
			)
			return d
		})
		return d
	})
	return d
}
func (p *pageItem) drawDeleteBookingsModal(gtx fwk.Gtx) fwk.Dim {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.85)
	if p.btnYes.Clicked() {
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
			startDateStr, startDate.Day(), startDate.Month(), startDate.Year(),
			endDateStr, endDate.Day(), endDate.Month().String(), endDate.Year(),
		),
		&p.btnYes,
		&p.btnNo,
	)
	return p.ModalContent.DrawContent(gtx, p.Theme, promptContent.Layout)
}

func (p *pageItem) drawBookingField(gtx fwk.Gtx, labelText string, valueText string) fwk.Dim {
	flex := layout.Flex{Alignment: layout.Start}
	return flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = gtx.Dp(labelWidth)
			gtx.Constraints.Min.X = gtx.Dp(labelWidth)
			label := material.Label(p.Theme, headFontSize, labelText)
			label.Font.Weight = headFontWeight
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Label(p.Theme, headFontSize, ": ").Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			ins := layout.Inset{}
			return ins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				b := material.Body1(p.Theme, valueText)
				b.Font.Weight = bodyFontWeight
				b.TextSize = bodyFontSize
				return b.Layout(gtx)
			})
		}),
	)
}

func (p *pageItem) drawBookingDate(gtx fwk.Gtx, t time.Time, labelStr string) fwk.Dim {
	flex := layout.Flex{Alignment: layout.Start}
	return flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = gtx.Dp(labelWidth)
			gtx.Constraints.Min.X = gtx.Dp(labelWidth)
			label := material.Label(p.Theme, headFontSize, labelStr)
			label.Font.Weight = headFontWeight
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Label(p.Theme, headFontSize, ": ").Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			ins := layout.Inset{}
			return ins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				b := material.Body1(p.Theme, p.GetFormattedDate(t))
				b.Font.Weight = bodyFontWeight
				b.TextSize = bodyFontSize
				return b.Layout(gtx)
			})
		}),
	)
}

func (p *pageItem) GetFormattedDate(t time.Time) string {
	var bookingDate string
	switch t.Day() {
	case 1, 21, 31:
		bookingDate = fmt.Sprintf("%dst", t.Day())
	case 2, 22:
		bookingDate = fmt.Sprintf("%dnd", t.Day())
	case 3, 23:
		bookingDate = fmt.Sprintf("%drd", t.Day())
	default:
		bookingDate = fmt.Sprintf("%dth", t.Day())
	}
	month := t.Month().String()
	year := t.Year()
	day := t.Weekday().String()
	bookingDate = fmt.Sprintf("%s %s, %s, %d", bookingDate, month, day, year)
	return bookingDate
}
