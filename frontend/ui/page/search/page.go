package search

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
	"strings"
)

type page struct {
	isSearchingBookings    bool
	isDeletingBooking      bool
	initialized            bool
	bookingsCount          int64
	limit                  int32
	offset                 int32
	query                  string
	title                  string
	ViewList               layout.List
	Theme                  *material.Theme
	buttonNavigation       widget.Clickable
	btnYes                 widget.Clickable
	btnNo                  widget.Clickable
	closeSnapBar           widget.Clickable
	navigationIcon         *widget.Icon
	bookingItems           []*pageItem
	subscription           service.Subscriber
	loginUserResponse      service.UserResponse
	searchBookingsChannels chan service.BookingsResponse
	formField              view.FormField
	fwk.Manager
}

func New(manager fwk.Manager) fwk.Page {
	navIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	errorTh := *user.Theme()
	errorTh.ContrastBg = color.NRGBA(colornames.Red500)
	theme := *user.Theme()
	p := page{
		Manager:                manager,
		Theme:                  &theme,
		navigationIcon:         navIcon,
		ViewList:               layout.List{Axis: layout.Vertical},
		bookingItems:           []*pageItem{},
		searchBookingsChannels: make(chan service.BookingsResponse, 1),
		limit:                  1000,
		offset:                 0,
		formField:              view.FormField{FieldName: i18n.Get(key.SearchBookings)},
	}
	p.subscription = manager.Service().Subscribe()
	p.subscription.SubscribeWithCallback(p.OnServiceStateChange)
	return &p
}

func (p *page) Layout(gtx fwk.Gtx) fwk.Dim {
	if !p.initialized {
		if p.Theme == nil {
			p.Theme = user.Theme()
		}
		p.searchBookings()
		p.initialized = true
	}
	if !p.loginUserResponse.IsAuthorized() {
		return fwk.Dim{}
	}
	val := strings.TrimSpace(p.formField.Text())
	if val != p.query && !p.isSearchingBookings {
		p.query = p.formField.Text()
		p.searchBookings()
	}
	p.title = i18n.Get(key.SearchBookings)

	flex := layout.Flex{Axis: layout.Vertical}
	d := flex.Layout(gtx,
		layout.Rigid(p.DrawAppBar),
		layout.Rigid(p.drawFormField),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return p.ViewList.Layout(gtx, len(p.bookingItems), func(gtx layout.Context, index int) layout.Dimensions {
				return p.drawBookingItem(gtx, index)
			})
		}),
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
							titleText := p.title
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

func (p *page) drawFormField(gtx fwk.Gtx) fwk.Dim {
	inset := layout.UniformInset(unit.Dp(16))
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		flex := layout.Flex{Axis: layout.Vertical}
		return flex.Layout(gtx,
			layout.Rigid(layout.Spacer{Height: 8}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return view.DrawFormFieldRowWithLabel(gtx, p.Theme, "", p.formField.FieldName, &p.formField.TextField, nil)
			}),
			layout.Rigid(layout.Spacer{Height: 8}.Layout),
		)
	})
}

func (p *page) drawBookingItem(gtx fwk.Gtx, index int) fwk.Dim {
	if len(p.bookingItems) == 0 {
		return view.Dim{}
	}
	return p.bookingItems[index].Layout(gtx)
}

func (p *page) searchBookings() {
	if !p.isSearchingBookings {
		p.isSearchingBookings = true
		p.Service().SearchBookings(p.formField.Text())
	}
}

func (p *page) OnServiceStateChange(event service.Event) {
	switch eventData := event.Data.(type) {
	case service.UpdateBookingResponse:
		if !p.isSearchingBookings {
			p.searchBookings()
		}
	case service.SearchBookingsResponse:
		p.isSearchingBookings = false
		if eventData.Error != "" {
			errStr := eventData.Error
			if strings.Contains(errStr, "connection refused") {
				errStr = "connection refused"
			}
			p.Snackbar().Show(errStr, &p.closeSnapBar, color.NRGBA{R: 255, A: 255}, i18n.Get(key.Close))
			p.bookingItems = []*pageItem{}
			break
		}
		var bookingItems []*pageItem
		for _, eachBooking := range eventData.Bookings {
			bookingItem := &pageItem{
				Booking:        eachBooking,
				BookingDetails: view.BookingDetails{Theme: p.Theme, Booking: eachBooking},
				parentPage:     p,
			}
			bookingItems = append(bookingItems, bookingItem)
		}
		p.bookingItems = bookingItems
	case service.UserResponse:
		p.loginUserResponse = eventData
		if !p.isSearchingBookings {
			p.searchBookings()
		}
	case service.DeleteBookingResponse:
		var txt string
		if (eventData.ID.Hex() == primitive.NilObjectID.Hex()) || !p.isDeletingBooking {
			return
		}
		if eventData.Error != "" {
			txt = fmt.Sprintf("couldn't delete booking with ID %s, error: %s", eventData.ID.Hex(), eventData.Error)
		}
		if eventData.Error == "" {
			txt = fmt.Sprintf("Successfully deleted booking with ID %s", eventData.ID.Hex())
		}
		if txt != "" {
			if !p.isSearchingBookings {
				p.searchBookings()
			}
			p.Snackbar().Show(txt, nil, color.NRGBA(colornames.White), "CLOSE")
		}
		p.isDeletingBooking = false
	}
}
func (p *page) URL() fwk.URL {
	return fwk.SearchBookingsPageURL
}
