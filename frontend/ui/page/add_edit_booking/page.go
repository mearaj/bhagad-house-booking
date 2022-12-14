package add_edit_booking

import (
	"errors"
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"github.com/mearaj/bhagad-house-booking/frontend/assets/fonts"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
	"strconv"
	"strings"
	"time"
)

type page struct {
	fwk.Manager
	Theme            *material.Theme
	navigationIcon   *widget.Icon
	buttonNavigation widget.Clickable
	initialized      bool
	bookingForm      view.BookingForm
	service.Booking
	subscription service.Subscriber
	isCreating   bool
	isUpdating   bool
	closeButton  widget.Clickable
}

func New(manager fwk.Manager, booking sqlc.Booking) fwk.Page {
	navIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	th := user.Theme()
	errorTh := *fonts.NewTheme()
	errorTh.ContrastBg = color.NRGBA(colornames.Red500)
	inActiveTh := *fonts.NewTheme()
	inActiveTh.ContrastBg = color.NRGBA(colornames.Grey500)
	s := page{
		Manager:        manager,
		Theme:          th,
		navigationIcon: navIcon,
		Booking:        booking,
		bookingForm: view.NewBookingForm(manager, booking, view.BookingFormParams{
			ShowDetails:      true,
			ShowCustomerName: true,
			ShowTotalPrice:   true,
		}),
	}
	s.subscription = manager.Service().Subscribe()
	s.subscription.SubscribeWithCallback(s.OnServiceStateChange)
	return &s
}

func (p *page) Layout(gtx fwk.Gtx) fwk.Dim {
	if p.Theme == nil {
		p.Theme = fonts.NewTheme()
	}

	flex := layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}
	p.bookingForm.ButtonText = i18n.Get(key.CreateNewBooking)
	if p.Booking.ID != 0 {
		p.bookingForm.ButtonText = i18n.Get(key.UpdateCurrentBooking)
	}
	if p.bookingForm.ButtonSubmit.Clicked() {
		isValid := !p.bookingForm.StartDate.IsZero() && !p.bookingForm.EndDate.IsZero() &&
			(p.bookingForm.StartDate.Before(p.bookingForm.EndDate) || p.bookingForm.StartDate.Equal(p.bookingForm.EndDate))
		var err error
		var totalPrice float64
		totalPriceStr := strings.TrimSpace(p.bookingForm.TotalPrice.Text())
		if !isValid {
			err = errors.New("make sure start date is before end date")
		}
		if isValid && totalPriceStr != "" {
			totalPrice, err = strconv.ParseFloat(totalPriceStr, 64)
			if err != nil {
				err = errors.New("please enter a valid Total Price")
			}
		}
		if err == nil {
			if p.Booking.ID == 0 {
				p.Service().CreateBooking(sqlc.CreateBookingParams{
					StartDate:    p.bookingForm.StartDate,
					EndDate:      p.bookingForm.EndDate,
					Details:      p.bookingForm.Details.Text(),
					CustomerName: p.bookingForm.CustomerName.Text(),
					TotalPrice:   totalPrice,
				})
				p.isCreating = true
			}
			if p.Booking.ID != 0 {
				p.Service().UpdateBooking(sqlc.UpdateBookingParams{
					ID:           p.Booking.ID,
					StartDate:    p.bookingForm.StartDate,
					EndDate:      p.bookingForm.EndDate,
					Details:      p.bookingForm.Details.Text(),
					CustomerName: p.bookingForm.CustomerName.Text(),
					TotalPrice:   totalPrice,
				})
				p.isUpdating = true
			}
		}
		if err != nil {
			p.Snackbar().Show(err.Error(), &p.closeButton, color.NRGBA(colornames.White), "Close")
		}
	}
	d := flex.Layout(gtx,
		layout.Rigid(p.DrawAppBar),
		layout.Rigid(p.bookingForm.Layout),
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
							titleText := fmt.Sprintf("Edit Booking %d", p.Booking.ID)
							if p.Booking.ID == 0 {
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
	return fwk.AddEditBookingPageURL(p.Booking.ID)
}
func (p *page) OnServiceStateChange(event service.Event) {
	startDate := utils.GetFirstDayOfMonth(time.Now().Local())
	endDate := utils.GetLastDayOfMonth(time.Now().Local().AddDate(0, 5, 0))
	bookingParams := service.BookingParams{
		StartDate: startDate,
		EndDate:   endDate,
	}
	switch eventData := event.Data.(type) {
	case service.CreateBookingResponse:
		if !p.isCreating {
			return
		}
		var txt string
		p.isCreating = false
		evErr := eventData.Error
		if evErr != "" {
			txt = fmt.Sprintf("Couldn't create booking, error: %s", evErr)
		}
		if evErr == "" {
			txt = fmt.Sprintf("Successfully created booking %d\n StartDate :- %s.\n EndDate:- %s.\n",
				eventData.Booking.ID, eventData.Booking.StartDate.Format("2006-01-02"),
				p.Booking.EndDate.Format("2006-01-02"))
		}
		if evErr == "" {
			p.Service().Bookings(bookingParams)
			p.PopUp()
		}
		p.Snackbar().Show(txt, &p.closeButton, color.NRGBA(colornames.White), "CLOSE")
	case service.UpdateBookingResponse:
		if !p.isUpdating {
			return
		}
		var txt string
		p.isUpdating = false
		errStr := eventData.Error
		if errStr != "" {
			txt = fmt.Sprintf("Couldn't update booking %d, error: %s", eventData.Booking.ID, errStr)
		}
		if errStr == "" {
			txt = fmt.Sprintf("Successfully updated booking %d\n StartDate :- %s.\n EndDate:- %s.\n",
				eventData.Booking.ID, eventData.Booking.StartDate.Format("2006-01-02"),
				p.Booking.EndDate.Format("2006-01-02"))
		}
		if errStr == "" {
			p.Service().Bookings(bookingParams)
			p.PopUp()
		}
		p.Snackbar().Show(txt, &p.closeButton, color.NRGBA(colornames.White), "CLOSE")
	}
}
