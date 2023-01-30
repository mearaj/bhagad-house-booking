package addedit

import (
	"errors"
	"fmt"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mearaj/bhagad-house-booking/frontend/assets/fonts"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/helper"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/page/transactions"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
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
	subscription       service.Subscriber
	closeButton        widget.Clickable
	loginUserResponse  service.UserResponse
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
	s.subscription = manager.Service().Subscribe()
	s.subscription.SubscribeWithCallback(s.OnServiceStateChange)
	return &s
}

func (p *page) Layout(gtx fwk.Gtx) fwk.Dim {
	if p.Theme == nil {
		p.Theme = fonts.NewTheme()
	}
	if !p.loginUserResponse.IsAuthorized() {
		return fwk.Dim{}
	}

	flex := layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}
	p.bookingDetailsForm.BtnText = i18n.Get(key.CreateNewBooking)
	if p.Booking.ID.Hex() != primitive.NilObjectID.Hex() {
		p.bookingDetailsForm.BtnText = i18n.Get(key.UpdateCurrentBooking)
	}
	isBookingNew := helper.IsNilObjectID(p.Booking.ID)
	if !isBookingNew && p.bookingDetailsForm.BtnManageTransactions.Clicked() {
		p.Manager.NavigateToPage(transactions.New(p.Manager, p.Booking))
		op.InvalidateOp{}.Add(gtx.Ops)
	}
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
			p.isCreatingUpdating = true
			if helper.IsNilObjectID(p.Booking.ID) {
				p.Service().CreateBooking(service.CreateBookingRequest{
					StartDate:    p.bookingDetailsForm.BookingForm.StartDate,
					EndDate:      p.bookingDetailsForm.BookingForm.EndDate,
					Details:      p.bookingDetailsForm.Details.Text(),
					CustomerName: p.bookingDetailsForm.CustomerName.Text(),
					RatePerDay:   ratePerDay,
				})
			} else {
				p.Service().UpdateBooking(service.UpdateBookingRequest{
					ID:           p.Booking.ID,
					StartDate:    p.bookingDetailsForm.BookingForm.StartDate,
					EndDate:      p.bookingDetailsForm.BookingForm.EndDate,
					Details:      p.bookingDetailsForm.Details.Text(),
					CustomerName: p.bookingDetailsForm.CustomerName.Text(),
					RatePerDay:   ratePerDay,
				})
			}
		}
		if err != nil {
			p.Snackbar().Show(err.Error(), &p.closeButton, color.NRGBA(colornames.White), "Close")
		}
	}

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
							startIndex := len(p.Booking.ID.Hex()) - 4
							titleText := fmt.Sprintf("%s %s", i18n.Get(key.EditBooking), p.Booking.ID.Hex()[startIndex:])
							if p.Booking.ID.Hex() == primitive.NilObjectID.Hex() {
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
	return fwk.AddEditBookingPageURL(p.Booking.ID.Hex())
}

func (p *page) OnServiceStateChange(event service.Event) {
	switch eventData := event.Data.(type) {
	case service.UserResponse:
		p.loginUserResponse = eventData
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
			txt = fmt.Sprintf("Successfully created booking %s\n StartDate :- %s.\n EndDate:- %s.\n",
				eventData.Booking.ID.Hex(), eventData.Booking.StartDate.Format("2006-01-02"),
				p.Booking.EndDate.Format("2006-01-02"))
		}
		if evErr == "" {
			p.Manager.NavigateToPage(transactions.New(p.Manager, eventData.Booking))
			p.Window().Invalidate()
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
			txt = fmt.Sprintf("Couldn't update booking %s, error: %s", eventData.Booking.ID.Hex(), errStr)
		}
		if errStr == "" {
			txt = fmt.Sprintf("Successfully updated booking %s\n StartDate :- %s.\n EndDate:- %s.\n",
				eventData.Booking.ID.Hex(), eventData.Booking.StartDate.Format("2006-01-02"),
				p.Booking.EndDate.Format("2006-01-02"))
		}
		if errStr == "" {
			p.Manager.NavigateToPage(transactions.New(p.Manager, eventData.Booking))
			p.Window().Invalidate()
		}
		p.Snackbar().Show(txt, &p.closeButton, color.NRGBA(colornames.White), "CLOSE")
	}
}
