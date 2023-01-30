package view

import (
	"fmt"
	giokey "gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/frontend/assets/fonts"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/helper"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"image/color"
	"strings"
	"time"
)

type BookingDetailsForm struct {
	Manager
	Theme                 *material.Theme
	initialized           bool
	BtnSubmit             widget.Clickable
	BtnManageTransactions widget.Clickable
	BookingForm           BookingDateForm
	Details               component.TextField
	TotalPrice            float64
	CustomerName          component.TextField
	RatePerDay            component.TextField
	BtnText               string
	booking               service.Booking
	layout.List
}

// NewBookingDetailsForm
func NewBookingDetailsForm(manager Manager, booking service.Booking) BookingDetailsForm {
	inActiveTheme := fonts.NewTheme()
	inActiveTheme.ContrastBg = color.NRGBA(colornames.Grey500)
	contForm := BookingDetailsForm{
		Manager:     manager,
		Theme:       user.Theme(),
		BookingForm: NewBookingForm(manager, booking, false),
		List:        layout.List{Axis: layout.Vertical},
		booking:     booking,
	}
	contForm.Details.SetText(booking.CustomerName)
	contForm.CustomerName.SetText(booking.Details)
	contForm.RatePerDay.InputHint = giokey.HintNumeric
	contForm.RatePerDay.SetText(fmt.Sprintf("%.2f", booking.RatePerDay))
	return contForm
}

func (bf *BookingDetailsForm) Layout(gtx Gtx) Dim {
	if !bf.initialized {
		if bf.Theme == nil {
			bf.Theme = fonts.NewTheme()
		}
		ModalContentInstance.OnCloseClick = func() {
			bf.Modal().Dismiss(nil)
		}
	}

	flex := layout.Flex{Axis: layout.Vertical}
	return bf.List.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
		return flex.Layout(gtx,

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return bf.BookingForm.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				inset := layout.UniformInset(16)
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					labelText := i18n.Get(key.CustomerName)
					return DrawFormFieldRowWithLabel(gtx, bf.Theme, labelText, labelText, &bf.CustomerName, nil)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				inset := layout.UniformInset(16)
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					labelText := i18n.Get(key.Details)
					return DrawFormFieldRowWithLabel(gtx, bf.Theme, labelText, labelText, &bf.Details, nil)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				inset := layout.UniformInset(16)
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					labelText := i18n.Get(key.RatePerDay)
					return DrawFormFieldRowWithLabel(gtx, bf.Theme, labelText, labelText, &bf.RatePerDay, nil)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				bf.setTotalPrice(gtx)
				inset := layout.UniformInset(16)
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					flex := layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}
					return flex.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							flex := layout.Flex{Axis: layout.Vertical}
							return flex.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									labelText := i18n.Get(key.TotalPrice)
									label := material.Label(bf.Theme, 16, labelText)
									label.Font.Weight = text.Bold
									return label.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Width: 16}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									totalPrice := fmt.Sprintf("%.2f", bf.TotalPrice)
									label := material.Label(bf.Theme, 16, totalPrice)
									return label.Layout(gtx)
								}),
							)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							btnManage := &bf.BtnManageTransactions
							isBookingNew := helper.IsNilObjectID(bf.booking.ID)
							bgColor := bf.Theme.ContrastBg
							if isBookingNew {
								btnManage = &widget.Clickable{}
								bgColor.A = 150
							}
							btnText := i18n.Get(key.ManageTransactions)
							btn := material.Button(bf.Theme, btnManage, btnText)
							btn.Background = bgColor
							return btn.Layout(gtx)
						}),
					)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				inset := layout.UniformInset(16)
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					bf.BtnText = i18n.Get(key.Key(bf.BtnText))
					if bf.BtnText == "" {
						bf.BtnText = i18n.Get(key.Apply)
					}
					btn := material.Button(bf.Theme, &bf.BtnSubmit, bf.BtnText)
					return btn.Layout(gtx)
				})
			}),
		)
	})
}

func (bf *BookingDetailsForm) drawDateField(gtx Gtx, label string, btnDate *widget.Clickable, btnClearDate *IconButton, t time.Time) Dim {
	fieldVal := i18n.Get(key.TapToEnterADate)
	labelWidth := gtx.Dp(100)
	flx := layout.Flex{Axis: layout.Vertical}
	inset := layout.UniformInset(unit.Dp(16))
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return flx.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X, gtx.Constraints.Max.X = labelWidth, labelWidth
				inset := layout.Inset{Bottom: unit.Dp(4)}
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return material.Label(bf.Theme, bf.Theme.TextSize, label).Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				flx := layout.Flex{Alignment: layout.Middle}
				return flx.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return btnDate.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return widget.Border{
								Color:        bf.Theme.ContrastBg,
								CornerRadius: 0,
								Width:        unit.Dp(1),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								inset := layout.UniformInset(unit.Dp(12))
								return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									lbl := material.Label(bf.Theme, bf.Theme.TextSize, fieldVal)
									lbl.Color = color.NRGBA(colornames.Grey500)
									lbl.Font.Weight = text.Normal
									if !t.IsZero() {
										lbl.Text = fmt.Sprintf("%d %s %d", t.Day(), t.Month().String(), t.Year())
										lbl.Color = bf.Theme.Fg
										lbl.Font.Weight = text.Bold
									}
									return lbl.Layout(gtx)
								})
							})
						})
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btnClearDate.Inset = layout.UniformInset(unit.Dp(8))
						return btnClearDate.Layout(gtx)
					}),
				)
			}),
		)
	})
}

func (bf *BookingDetailsForm) setTotalPrice(gtx fwk.Gtx) {
	ratePerDayStr := strings.TrimSpace(bf.RatePerDay.Text())
	startDate := bf.BookingForm.StartDate
	endDate := bf.BookingForm.EndDate
	bf.TotalPrice = helper.BookingTotalPriceFromRateStr(ratePerDayStr, startDate, endDate)
}
