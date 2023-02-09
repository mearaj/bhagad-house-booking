package view

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"time"
)

type BookingDetails struct {
	*material.Theme
	service.Booking
	layout.Inset
}

const BookingDetailsHeadFontSize = unit.Sp(18)
const BookingDetailsBodyFontSize = unit.Sp(18)
const BookingDetailsHeadFontWeight = text.ExtraBold
const BookingDetailsBodyFontWeight = text.Medium
const BookingDetailsLabelWidth = unit.Dp(160)
const BookingDetailsFormInset = unit.Dp(16)

func (p *BookingDetails) Layout(gtx fwk.Gtx) fwk.Dim {
	if p.Theme == nil {
		p.Theme = user.Theme()
	}
	return p.layoutContent(gtx)
}
func (p *BookingDetails) layoutContent(gtx fwk.Gtx) fwk.Dim {
	d := p.Inset.Layout(gtx, func(gtx fwk.Gtx) fwk.Dim {
		flex := layout.Flex{Axis: layout.Vertical}
		d := flex.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				labelTxt := fmt.Sprintf("%s %s", i18n.Get(key.Booking), i18n.Get(key.NumberShort))
				valueTxt := fmt.Sprintf("%d", p.Booking.Number)
				return p.drawBookingField(gtx, labelTxt, valueTxt)
			}),
			layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
				labelTxt := i18n.Get(key.CustomerName)
				valueTxt := p.Booking.CustomerName
				return p.drawBookingField(gtx, labelTxt, valueTxt)
			}),
			layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
				labelTxt := i18n.Get(key.CustomerPhone)
				valueTxt := p.Booking.CustomerPhone
				return p.drawBookingField(gtx, labelTxt, valueTxt)
			}),
			layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
				labelTxt := i18n.Get(key.CustomerEmail)
				valueTxt := p.Booking.CustomerEmail
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
			layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
				labelTxt := i18n.Get(key.Period)
				valueInt := utils.BookingTotalNumberOfDays(p.StartDate, p.EndDate)
				valueTxt := fmt.Sprintf("%d", valueInt)
				if valueInt > 1 {
					valueTxt += " Days"
				} else {
					valueTxt += " Day"
				}
				return p.drawBookingField(gtx, labelTxt, valueTxt)
			}),
			layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
				labelTxt := i18n.Get(key.RatePerDay)
				valueTxt := fmt.Sprintf("%.2f", p.Booking.RatePerDay)
				return p.drawBookingField(gtx, labelTxt, valueTxt)
			}),
			layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
				labelTxt := i18n.Get(key.TotalPrice)
				startDate := p.Booking.StartDate
				endDate := p.Booking.EndDate
				ratePerDay := p.Booking.RatePerDay
				valueTxt := fmt.Sprintf("%.2f",
					utils.BookingTotalPrice(ratePerDay, startDate, endDate),
				)
				return p.drawBookingField(gtx, labelTxt, valueTxt)
			}),
		)
		return d
	})
	return d
}

func (p *BookingDetails) drawLabelField(gtx fwk.Gtx, labelText string) [2]layout.FlexChild {
	return [2]layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = gtx.Dp(BookingDetailsLabelWidth)
			gtx.Constraints.Min.X = gtx.Dp(BookingDetailsLabelWidth)
			label := material.Label(p.Theme, BookingDetailsHeadFontSize, labelText)
			label.Font.Weight = BookingDetailsHeadFontWeight
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Label(p.Theme, BookingDetailsHeadFontSize, ": ").Layout(gtx)
		}),
	}
}

func (p *BookingDetails) drawBookingField(gtx fwk.Gtx, labelText string, valueText string) fwk.Dim {
	flex := layout.Flex{Alignment: layout.Start}
	labelField := p.drawLabelField(gtx, labelText)
	return flex.Layout(gtx,
		labelField[0],
		labelField[1],
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			ins := layout.Inset{}
			return ins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				b := material.Body1(p.Theme, valueText)
				b.Font.Weight = BookingDetailsBodyFontWeight
				b.TextSize = BookingDetailsBodyFontSize
				return b.Layout(gtx)
			})
		}),
	)
}

func (p *BookingDetails) drawBookingDate(gtx fwk.Gtx, t time.Time, labelStr string) fwk.Dim {
	flex := layout.Flex{Alignment: layout.Start}
	return flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = gtx.Dp(BookingDetailsLabelWidth)
			gtx.Constraints.Min.X = gtx.Dp(BookingDetailsLabelWidth)
			label := material.Label(p.Theme, BookingDetailsHeadFontSize, labelStr)
			label.Font.Weight = BookingDetailsHeadFontWeight
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Label(p.Theme, BookingDetailsHeadFontSize, ": ").Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			ins := layout.Inset{}
			return ins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				b := material.Body1(p.Theme, utils.GetFormattedDate(t))
				b.Font.Weight = BookingDetailsBodyFontWeight
				b.TextSize = BookingDetailsBodyFontSize
				return b.Layout(gtx)
			})
		}),
	)
}
func DrawBookingDetailsLabelField(gtx fwk.Gtx, theme *material.Theme, labelText string) [2]layout.FlexChild {
	return [2]layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = gtx.Dp(BookingDetailsLabelWidth)
			gtx.Constraints.Min.X = gtx.Dp(BookingDetailsLabelWidth)
			label := material.Label(theme, BookingDetailsHeadFontSize, labelText)
			label.Font.Weight = BookingDetailsHeadFontWeight
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Label(theme, BookingDetailsHeadFontSize, ": ").Layout(gtx)
		}),
	}
}
