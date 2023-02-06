package addedit

import (
	"errors"
	"fmt"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"github.com/mearaj/bhagad-house-booking/frontend/assets/fonts"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/page/transactions"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
	"net/mail"
	"strconv"
	"strings"
)

type page struct {
	initialized        bool
	isCreatingUpdating bool
	service.Booking
	fwk.Manager
	Theme              *material.Theme
	navigationIcon     *widget.Icon
	buttonNavigation   widget.Clickable
	bookingDetailsForm view.BookingDetailsForm
	closeButton        widget.Clickable
}

func New(manager fwk.Manager, booking service.Booking) fwk.Page {
	navIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	th := user.Theme()
	errorTh := *fonts.NewTheme()
	errorTh.ContrastBg = color.NRGBA(colornames.Red500)
	inActiveTh := *fonts.NewTheme()
	inActiveTh.ContrastBg = color.NRGBA(colornames.Grey500)
	s := page{
		Manager:            manager,
		Theme:              th,
		navigationIcon:     navIcon,
		Booking:            booking,
		bookingDetailsForm: view.NewBookingDetailsForm(manager, booking),
	}
	return &s
}

func (p *page) Layout(gtx fwk.Gtx) fwk.Dim {
	if p.Theme == nil {
		p.Theme = fonts.NewTheme()
	}
	if !p.Manager.User().IsAuthorized() {
		return fwk.Dim{}
	}

	flex := layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}
	p.bookingDetailsForm.BtnText = i18n.Get(key.CreateNewBooking)
	if p.Booking.Number != 0 {
		p.bookingDetailsForm.BtnText = i18n.Get(key.UpdateCurrentBooking)
	}
	isBookingNew := p.Booking.Number == 0
	if !isBookingNew && p.bookingDetailsForm.BtnManageTransactions.Clicked() {
		p.Manager.NavigateToPage(transactions.New(p.Manager, p.Booking))
		op.InvalidateOp{}.Add(gtx.Ops)
	}
	p.handleFormSubmit()
	p.handleSendSMSSubmit()
	p.handleSendEmailSubmit()
	d := flex.Layout(gtx,
		layout.Rigid(p.DrawAppBar),
		layout.Rigid(p.bookingDetailsForm.Layout),
	)
	return d
}
func (p *page) DrawAppBar(gtx fwk.Gtx) fwk.Dim {
	gtx.Constraints.Max.Y = gtx.Dp(56)
	th := p.Theme
	if p.buttonNavigation.Clicked() {
		p.PopUp()
	}
	return view.DrawAppBarLayout(gtx, th, func(gtx fwk.Gtx) fwk.Dim {
		return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
						navigationIcon := p.navigationIcon
						button := material.IconButton(th, &p.buttonNavigation, navigationIcon, "Nav Icon Button")
						button.Size = unit.Dp(40)
						button.Background = th.Palette.ContrastBg
						button.Color = th.Palette.ContrastFg
						button.Inset = layout.UniformInset(unit.Dp(8))
						return button.Layout(gtx)
					}),
					layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
						return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx fwk.Gtx) fwk.Dim {
							titleText := fmt.Sprintf("%s %d", i18n.Get(key.EditBooking), p.Booking.Number)
							if p.Booking.Number == 0 {
								titleText = i18n.Get(key.CreateNewBooking)
							}
							title := material.Body1(th, titleText)
							title.Color = th.Palette.ContrastFg
							title.TextSize = unit.Sp(18)
							return title.Layout(gtx)
						})
					}),
				)
			}),
		)
	})
}
func (p *page) URL() fwk.URL {
	return fwk.AddEditBookingPageURL(fmt.Sprintf("%d", p.Booking.Number))
}
func (p *page) handleFormSubmit() {
	if p.bookingDetailsForm.BtnSubmit.Clicked() {
		isValid := !p.bookingDetailsForm.BookingForm.StartDate.IsZero() && !p.bookingDetailsForm.BookingForm.EndDate.IsZero() &&
			(p.bookingDetailsForm.BookingForm.StartDate.Before(p.bookingDetailsForm.BookingForm.EndDate) || p.bookingDetailsForm.BookingForm.StartDate.Equal(p.bookingDetailsForm.BookingForm.EndDate))
		var err error
		var ratePerDay float64
		ratePerDayStr := strings.TrimSpace(p.bookingDetailsForm.RatePerDay.Text())
		if !isValid {
			err = errors.New("make sure start date is before end date")
		}
		if isValid && ratePerDayStr != "" {
			ratePerDay, err = strconv.ParseFloat(ratePerDayStr, 64)
			if err != nil {
				err = errors.New("please enter a valid rate per day")
			}
		}
		if err == nil {
			phoneNumber := p.bookingDetailsForm.CustomerPhone.Text()
			var isPhoneValid bool
			if strings.TrimSpace(phoneNumber) == "" {
				isPhoneValid = true
			} else {
				isPhoneValid = utils.ValidateIndianPhoneNumber(phoneNumber)
			}
			if !isPhoneValid {
				err = errors.New("invalid phone number")
				p.bookingDetailsForm.CustomerPhone.SetError(err.Error())
			}
		}
		if err == nil {
			email := p.bookingDetailsForm.CustomerEmail.Text()
			var isEmailValid bool
			if strings.TrimSpace(email) == "" {
				isEmailValid = true
			} else {
				_, err = mail.ParseAddress(email)
				isEmailValid = err == nil
			}
			if !isEmailValid {
				err = errors.New("invalid email")
				p.bookingDetailsForm.CustomerEmail.SetError(err.Error())
			}
		}
		if err == nil {
			p.isCreatingUpdating = true
			if p.Booking.Number == 0 {
				p.Service().CreateBooking(service.CreateBookingRequest{
					StartDate:     p.bookingDetailsForm.BookingForm.StartDate,
					EndDate:       p.bookingDetailsForm.BookingForm.EndDate,
					Details:       p.bookingDetailsForm.Details.TextField.Text(),
					CustomerName:  p.bookingDetailsForm.CustomerName.TextField.Text(),
					CustomerEmail: p.bookingDetailsForm.CustomerEmail.Text(),
					CustomerPhone: p.bookingDetailsForm.CustomerPhone.Text(),
					RatePerDay:    ratePerDay,
				})
			} else {
				p.Service().UpdateBooking(service.UpdateBookingRequest{
					Number:           p.Booking.Number,
					StartDate:        p.bookingDetailsForm.BookingForm.StartDate,
					EndDate:          p.bookingDetailsForm.BookingForm.EndDate,
					Details:          p.bookingDetailsForm.Details.TextField.Text(),
					CustomerName:     p.bookingDetailsForm.CustomerName.TextField.Text(),
					CustomerEmail:    p.bookingDetailsForm.CustomerEmail.Text(),
					CustomerPhone:    p.bookingDetailsForm.CustomerPhone.Text(),
					RatePerDay:       ratePerDay,
					ConfirmEmailSent: p.Booking.ConfirmEmailSent,
					ConfirmSMSSent:   p.Booking.ConfirmSMSSent,
				})
			}
		}
		if err != nil {
			p.Snackbar().Show(err.Error(), &p.closeButton, color.NRGBA(colornames.White), "Close")
		}
	}
}

func (p *page) handleSendSMSSubmit() {
	if p.bookingDetailsForm.BtnSendSMS.Button.Clicked() {
		var err error
		if p.bookingDetailsForm.CustomerPhone.Text() != p.Booking.CustomerPhone {
			err = errors.New("phone number mismatch")
		}
		if err == nil {
			phoneNumber := p.bookingDetailsForm.CustomerPhone.Text()
			isValid := utils.ValidateIndianPhoneNumber(phoneNumber)
			if !isValid {
				err = errors.New("invalid phone number")
			}
		}
		if err != nil {
			p.bookingDetailsForm.CustomerPhone.SetError(err.Error())
		}
		if err == nil {
			p.bookingDetailsForm.IsSendingSMS = true
			p.Service().SendNewBookingSMS(p.Booking.Number, p)
		}
	}
}
func (p *page) handleSendEmailSubmit() {
	if p.bookingDetailsForm.BtnSendEmail.Button.Clicked() {
		var err error
		if p.bookingDetailsForm.CustomerEmail.Text() != p.Booking.CustomerEmail {
			err = errors.New("email mismatch")
		}
		if err == nil {
			email := p.bookingDetailsForm.CustomerEmail.Text()
			_, err = mail.ParseAddress(email)
			if err != nil {
				err = errors.New("invalid email")
			}
		}
		if err != nil {
			p.bookingDetailsForm.CustomerEmail.SetError(err.Error())
		}
		if err == nil {
			p.bookingDetailsForm.IsSendingEmail = true
			p.Service().SendNewBookingEmail(p.Booking.Number, p)
		}
	}
}
func (p *page) OnServiceStateChange(event service.Event) {
	switch eventData := event.Data.(type) {
	case service.UserResponse:
		p.Window().Invalidate()
	case service.NewBookingSMSResponse:
		p.bookingDetailsForm.IsSendingSMS = false
		var txt string
		evErr := eventData.Error
		if evErr != "" {
			txt = fmt.Sprintf("send sms error: %s", evErr)
		}
		if evErr == "" {
			txt = "sms successfully sent."
		}
		if evErr == "" {
			p.replaceCurrentBooking(eventData.Booking)
			p.Window().Invalidate()
		}
		p.Snackbar().Show(txt, &p.closeButton, color.NRGBA(colornames.White), "CLOSE")
	case service.NewBookingEmailResponse:
		p.bookingDetailsForm.IsSendingEmail = false
		var txt string
		evErr := eventData.Error
		if evErr != "" {
			txt = fmt.Sprintf("send email error: %s", evErr)
		}
		if evErr == "" {
			txt = "email successfully sent."
		}
		if evErr == "" {
			p.replaceCurrentBooking(eventData.Booking)
		}
		p.Snackbar().Show(txt, &p.closeButton, color.NRGBA(colornames.White), "CLOSE")
	case service.CreateBookingResponse:
		if !p.isCreatingUpdating || event.Cached {
			return
		}
		var txt string
		p.isCreatingUpdating = false
		evErr := eventData.Error
		if evErr != "" {
			txt = fmt.Sprintf("Couldn't create booking, error: %s", evErr)
		}
		if evErr == "" {
			txt = fmt.Sprintf("Successfully created booking %d\n StartDate :- %s.\n EndDate:- %s.\n",
				eventData.Booking.Number, eventData.Booking.StartDate.Format("2006-01-02"),
				p.Booking.EndDate.Format("2006-01-02"))
		}
		if evErr == "" {
			p.replaceCurrentBooking(eventData.Booking)
		}
		p.Snackbar().Show(txt, &p.closeButton, color.NRGBA(colornames.White), "CLOSE")
	case service.UpdateBookingResponse:
		if !p.isCreatingUpdating || event.Cached {
			return
		}
		var txt string
		p.isCreatingUpdating = false
		errStr := eventData.Error
		if errStr != "" {
			txt = fmt.Sprintf("Couldn't update booking %d, error: %s", eventData.Booking.Number, errStr)
		}
		if errStr == "" {
			txt = fmt.Sprintf("Successfully updated booking %d\n StartDate :- %s.\n EndDate:- %s.\n",
				eventData.Booking.Number, eventData.Booking.StartDate.Format("2006-01-02"),
				p.Booking.EndDate.Format("2006-01-02"))
		}
		if errStr == "" {
			p.replaceCurrentBooking(eventData.Booking)
		}
		p.Snackbar().Show(txt, &p.closeButton, color.NRGBA(colornames.White), "CLOSE")
	}
}

func (p *page) replaceCurrentBooking(booking service.Booking) {
	p.Booking = booking
	p.bookingDetailsForm = view.NewBookingDetailsForm(p.Manager, booking)
	p.Window().Invalidate()
}
