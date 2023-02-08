package search

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/page/addedit"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/page/transactions"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"image/color"
	"time"
)

type pageItem struct {
	btnEdit         widget.Clickable
	btnDelete       widget.Clickable
	btnTransactions widget.Clickable
	btnYes          widget.Clickable
	btnNo           widget.Clickable
	service.Booking
	ModalContent *view.ModalContent
	parentPage   *page
	view.BookingDetails
}

func (p *pageItem) Layout(gtx fwk.Gtx) fwk.Dim {
	if p.ModalContent == nil {
		p.ModalContent = view.NewModalContent(func() {
			p.parentPage.Modal().Dismiss(nil)
		})
	}
	return p.layoutContent(gtx)
}
func (p *pageItem) layoutContent(gtx fwk.Gtx) fwk.Dim {
	inset := layout.UniformInset(16)
	if p.btnDelete.Clicked() && p.Booking.Number != 0 {
		p.parentPage.Modal().Show(p.drawDeleteBookingsModal, nil, fwk.Animation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		})
	}

	if p.btnEdit.Clicked() {
		p.parentPage.NavigateToPage(addedit.New(p.parentPage, p.Booking))
		op.InvalidateOp{}.Add(gtx.Ops)
	}
	if p.btnTransactions.Clicked() {
		p.parentPage.NavigateToPage(transactions.New(p.parentPage, p.Booking))
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	d := inset.Layout(gtx, func(gtx fwk.Gtx) fwk.Dim {
		flex := layout.Flex{Axis: layout.Vertical}
		d := flex.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return p.BookingDetails.Layout(gtx)
			}),
			layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
				labelField := p.drawLabelField(gtx, i18n.Get(key.SelectAction))
				flex := layout.Flex{Alignment: layout.Middle}
				return flex.Layout(gtx,
					labelField[0],
					labelField[1],
					p.drawActionButtons(gtx),
				)
			}),
			layout.Rigid(layout.Spacer{Height: 24}.Layout),
			layout.Rigid(component.Divider(p.Theme).Layout),
		)
		return d
	})
	return d
}
func (p *pageItem) drawDeleteBookingsModal(gtx fwk.Gtx) fwk.Dim {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.85)
	if p.btnYes.Clicked() {
		p.parentPage.isDeletingBooking = true
		p.parentPage.Service().DeleteBooking(p.Booking.Number)
		p.parentPage.Modal().Dismiss(func() { p.parentPage.Window().Invalidate() })
	}
	if p.btnNo.Clicked() {
		p.parentPage.Modal().Dismiss(func() {})
	}

	delPrompt := i18n.Get(key.BookingDeletePrompt)
	bookingID := fmt.Sprintf("%s %s %d", i18n.Get(key.Booking), i18n.Get(key.NumberShort), p.Booking.Number)
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

func (p *pageItem) drawLabelField(gtx fwk.Gtx, labelText string) [2]layout.FlexChild {
	return view.DrawBookingDetailsLabelField(gtx, p.Theme, labelText)
}
func (p *pageItem) drawActionButtons(gtx fwk.Gtx) layout.FlexChild {
	return layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		flex := layout.Flex{Alignment: layout.Middle}
		return flex.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = 80
				editLabel := i18n.Get(key.Edit)
				btn := material.Button(p.Theme, &p.btnEdit, editLabel)
				return btn.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Width: 16}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				manageLabel := i18n.Get(key.Transactions)
				btn := material.Button(p.Theme, &p.btnTransactions, manageLabel)
				return btn.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Width: 16}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = 80
				delLabel := i18n.Get(key.Delete)
				btn := material.Button(p.Theme, &p.btnDelete, delLabel)
				btn.Background = color.NRGBA(colornames.Red500)
				return btn.Layout(gtx)
			}),
		)
	})
}

func (p *pageItem) drawBookingField(gtx fwk.Gtx, labelText string, valueText string) fwk.Dim {
	flex := layout.Flex{Alignment: layout.Start}
	labelField := p.drawLabelField(gtx, labelText)
	return flex.Layout(gtx,
		labelField[0],
		labelField[1],
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			ins := layout.Inset{}
			return ins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				b := material.Body1(p.Theme, valueText)
				b.Font.Weight = view.BookingDetailsBodyFontWeight
				b.TextSize = view.BookingDetailsBodyFontSize
				return b.Layout(gtx)
			})
		}),
	)
}

func (p *pageItem) drawBookingDate(gtx fwk.Gtx, t time.Time, labelStr string) fwk.Dim {
	flex := layout.Flex{Alignment: layout.Start}
	return flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = gtx.Dp(view.BookingDetailsLabelWidth)
			gtx.Constraints.Min.X = gtx.Dp(view.BookingDetailsLabelWidth)
			label := material.Label(p.Theme, view.BookingDetailsHeadFontSize, labelStr)
			label.Font.Weight = view.BookingDetailsHeadFontWeight
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Label(p.Theme, view.BookingDetailsHeadFontSize, ": ").Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			ins := layout.Inset{}
			return ins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				b := material.Body1(p.Theme, utils.GetFormattedDate(t))
				b.Font.Weight = view.BookingDetailsBodyFontWeight
				b.TextSize = view.BookingDetailsBodyFontSize
				return b.Layout(gtx)
			})
		}),
	)
}
