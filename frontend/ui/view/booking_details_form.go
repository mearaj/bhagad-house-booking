package view

import (
	"fmt"
	giokey "gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"github.com/mearaj/bhagad-house-booking/frontend/assets/fonts"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"strings"
)

type BookingDetailsForm struct {
	Manager
	Theme                 *material.Theme
	initialized           bool
	IsSendingEmail        bool
	IsSendingSMS          bool
	BtnSubmit             widget.Clickable
	BtnManageTransactions widget.Clickable
	BtnSendEmail          material.ButtonStyle
	BtnSendSMS            material.ButtonStyle
	BookingForm           BookingDateForm
	TotalPrice            float64
	Details               FormField
	CustomerName          FormField
	CustomerEmail         component.TextField
	CustomerPhone         component.TextField
	RatePerDay            component.TextField
	BtnText               string
	Booking               service.Booking
	IconDone              *widget.Icon
	prevPhoneNumber       string
	prevEmail             string
	layout.List
}

// NewBookingDetailsForm
func NewBookingDetailsForm(manager Manager, booking service.Booking) BookingDetailsForm {
	inActiveTheme := fonts.NewTheme()
	iconDone, _ := widget.NewIcon(icons.ActionDone)
	inActiveTheme.ContrastBg = color.NRGBA(colornames.Grey500)
	contForm := BookingDetailsForm{
		Manager: manager,
		Theme:   user.Theme(),
		BookingForm: NewBookingForm(manager, service.BookingsRequest{
			StartDate: booking.StartDate,
			EndDate:   booking.EndDate,
		}, false),
		List:           layout.List{Axis: layout.Vertical},
		Booking:        booking,
		CustomerEmail:  component.TextField{Editor: widget.Editor{InputHint: giokey.HintEmail, SingleLine: true}},
		CustomerPhone:  component.TextField{Editor: widget.Editor{InputHint: giokey.HintTelephone, SingleLine: true}},
		BtnSendEmail:   material.Button(user.Theme(), &widget.Clickable{}, ""),
		BtnSendSMS:     material.Button(user.Theme(), &widget.Clickable{}, ""),
		IconDone:       iconDone,
		IsSendingEmail: false,
		IsSendingSMS:   false,
	}
	contForm.Details.TextField.SetText(booking.CustomerName)
	contForm.CustomerName.TextField.SetText(booking.Details)
	contForm.CustomerEmail.SetText(booking.CustomerEmail)
	contForm.CustomerPhone.SetText(booking.CustomerPhone)
	contForm.RatePerDay.InputHint = giokey.HintNumeric
	contForm.RatePerDay.SetText(fmt.Sprintf("%.2f", booking.RatePerDay))
	if booking.Number == 0 {
		contForm.RatePerDay.SetText(fmt.Sprintf("%.2f", user.BookingRate()))
	}
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
	if bf.prevPhoneNumber != bf.CustomerPhone.Text() {
		bf.CustomerPhone.ClearError()
		bf.prevPhoneNumber = bf.CustomerPhone.Text()
	}
	if bf.prevEmail != bf.CustomerEmail.Text() {
		bf.CustomerEmail.ClearError()
		bf.prevEmail = bf.CustomerEmail.Text()
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
					bf.CustomerName.LabelHintText = labelText
					bf.CustomerName.FieldName = labelText
					bf.CustomerName.Theme = bf.Theme
					return bf.CustomerName.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				inset := layout.UniformInset(16)
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					labelText := i18n.Get(key.CustomerEmail)
					sendEmailText := i18n.Get(key.SendEmail)
					isBookingNew := bf.Booking.Number == 0
					btn := &bf.BtnSendEmail
					iconDone := func(gtx fwk.Gtx) Dim {
						iconColor := color.NRGBA(colornames.Green500)
						if bf.IsSendingEmail {
							loader := Loader{
								Theme: bf.Theme,
								Size:  image.Point{X: gtx.Dp(32), Y: gtx.Dp(32)},
							}
							return loader.Layout(gtx)
						}

						return bf.IconDone.Layout(gtx, iconColor)
					}
					if isBookingNew {
						btn = nil
						iconDone = nil
					}
					if !bf.Booking.ConfirmEmailSent && !bf.IsSendingEmail {
						iconDone = nil
					}
					bf.BtnSendEmail.Text = sendEmailText
					return DrawFormField(gtx, bf.Theme, labelText, labelText, &bf.CustomerEmail, nil, btn, iconDone)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				inset := layout.UniformInset(16)
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					labelText := i18n.Get(key.CustomerPhone)
					sendSMSText := i18n.Get(key.SendSMS)
					isBookingNew := bf.Booking.Number == 0
					btn := &bf.BtnSendSMS
					iconDone := func(gtx fwk.Gtx) Dim {
						iconColor := color.NRGBA(colornames.Green500)
						if bf.IsSendingSMS {
							loader := Loader{
								Theme: bf.Theme,
								Size:  image.Point{X: gtx.Dp(32), Y: gtx.Dp(32)},
							}
							return loader.Layout(gtx)
						}

						return bf.IconDone.Layout(gtx, iconColor)
					}
					if isBookingNew {
						btn = nil
						iconDone = nil
					}
					if !bf.Booking.ConfirmSMSSent && !bf.IsSendingSMS {
						iconDone = nil
					}
					bf.BtnSendSMS.Text = sendSMSText
					return DrawFormField(gtx, bf.Theme, labelText, labelText, &bf.CustomerPhone, nil, btn, iconDone)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				inset := layout.UniformInset(16)
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					labelText := i18n.Get(key.Details)
					bf.Details.FieldName = labelText
					bf.Details.LabelHintText = labelText
					return bf.Details.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				inset := layout.UniformInset(16)
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					labelText := i18n.Get(key.RatePerDay)
					return DrawFormField(gtx, bf.Theme, labelText, labelText, &bf.RatePerDay, nil, nil, nil)
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
							isBookingNew := bf.Booking.Number == 0
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

func (bf *BookingDetailsForm) setTotalPrice(gtx fwk.Gtx) {
	ratePerDayStr := strings.TrimSpace(bf.RatePerDay.Text())
	startDate := bf.BookingForm.StartDate
	endDate := bf.BookingForm.EndDate
	bf.TotalPrice = utils.BookingTotalPriceFromRateStr(ratePerDayStr, startDate, endDate)
}
