package bookings

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mearaj/bhagad-house-booking/common/utils"
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
	"time"
)

type page struct {
	ViewList         layout.List
	ViewListPosition layout.Position
	fwk.Manager
	Theme              *material.Theme
	title              string
	buttonNavigation   widget.Clickable
	btnAddBooking      widget.Clickable
	btnMenuContent     widget.Clickable
	btnYes             widget.Clickable
	btnNo              widget.Clickable
	closeIcon          *widget.Icon
	navigationIcon     *widget.Icon
	bookingItems       []*pageItem
	isFetchingBookings bool
	isDeletingBooking  bool
	initialized        bool
	subscription       service.Subscriber
	fetchingBookingsCh chan service.BookingsResponse
	bookingsCount      int64
	bookingForm        view.BookingDateForm
	limit              int32
	offset             int32
	loginUserResponse  service.UserResponse
	closeSnapBar       widget.Clickable
}

var flexUserLayoutRatio = [5]float32{0.16, 0.28, 0.28, 0.28, 0}
var flexAdminLayoutRatio = [5]float32{0.16, 0.21, 0.21, 0.21, 0.21}

func New(manager fwk.Manager) fwk.Page {
	navIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	closeIcon, _ := widget.NewIcon(icons.ContentClear)
	errorTh := *user.Theme()
	errorTh.ContrastBg = color.NRGBA(colornames.Red500)
	startDate := utils.GetFirstDayOfMonth(time.Now().Local())
	endDate := utils.GetLastDayOfMonth(time.Now().Local().AddDate(0, 5, 0))
	theme := *user.Theme()
	p := page{
		Manager:            manager,
		Theme:              &theme,
		title:              i18n.Get(key.Bookings),
		navigationIcon:     navIcon,
		ViewList:           layout.List{Axis: layout.Vertical},
		bookingItems:       []*pageItem{},
		closeIcon:          closeIcon,
		fetchingBookingsCh: make(chan service.BookingsResponse, 1),
		limit:              1000,
		offset:             0,
		bookingForm:        view.NewBookingForm(manager, service.Booking{StartDate: startDate, EndDate: endDate}, true),
	}
	p.subscription = manager.Service().Subscribe(service.TopicFetchBookings, service.TopicLoggedInOut, service.TopicDeleteBooking)
	p.subscription.SubscribeWithCallback(p.OnServiceStateChange)
	return &p
}

func (p *page) Layout(gtx fwk.Gtx) fwk.Dim {
	if !p.initialized {
		if p.Theme == nil {
			p.Theme = user.Theme()
		}
		p.fetchBookings()
		p.initialized = true
	}
	p.title = i18n.Get(key.Bookings)

	if p.bookingForm.ButtonSubmit.Clicked() {
		shouldFetch := !p.isFetchingBookings &&
			!p.bookingForm.StartDate.IsZero() &&
			!p.bookingForm.EndDate.IsZero() &&
			(p.bookingForm.EndDate.After(p.bookingForm.StartDate) ||
				p.bookingForm.EndDate.Equal(p.bookingForm.StartDate))
		if shouldFetch {
			p.fetchBookings()
		}
	}

	flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEnd, Alignment: layout.Start}
	d := flex.Layout(gtx,
		layout.Rigid(p.DrawAppBar),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(gtx, layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				position := p.ViewListPosition.First
				if position >= 1 && len(p.bookingItems) >= position {
					currItemTime := p.bookingItems[position-1].Time
					currStartMonth := currItemTime.Month()
					currStartYear := currItemTime.Year()
					flex := layout.Flex{Axis: layout.Vertical}
					return flex.Layout(gtx, p.drawBookingItemHeader(gtx, currStartMonth, currStartYear)...)
				}
				return fwk.Dim{}
			}))
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.ViewList.Layout(gtx, len(p.bookingItems)+1, func(gtx layout.Context, index int) layout.Dimensions {
				if index == 0 {
					return p.bookingForm.Layout(gtx)
				}
				return p.drawBookingItem(gtx, index-1)
			})
		}),
	)
	p.ViewListPosition = p.ViewList.Position
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

func (p *page) drawMenuItem(txt string, btn *widget.Clickable) layout.FlexChild {
	inset := layout.UniformInset(unit.Dp(12))
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		btnStyle := material.ButtonLayoutStyle{Button: btn}
		btnStyle.Background = color.NRGBA(colornames.White)
		return btnStyle.Layout(gtx,
			func(gtx fwk.Gtx) fwk.Dim {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				inset := inset
				return inset.Layout(gtx, func(gtx fwk.Gtx) fwk.Dim {
					return layout.Flex{Spacing: layout.SpaceEnd}.Layout(gtx,
						layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
							bd := material.Body1(p.Theme, txt)
							bd.Color = color.NRGBA(colornames.Black)
							bd.Alignment = text.Start
							return bd.Layout(gtx)
						}),
					)
				})
			},
		)
	})
}

func (p *page) drawBookingItem(gtx fwk.Gtx, index int) fwk.Dim {
	if len(p.bookingItems) == 0 {
		return view.Dim{}
	}
	showHeader := true
	currItemTime := p.bookingItems[index].Time
	currStartMonth := currItemTime.Month()
	currStartYear := currItemTime.Year()
	position := p.ViewListPosition.First
	if index > 0 {
		prevItemTime := p.bookingItems[index-1].Time
		prevStartMonth := prevItemTime.Month()
		prevStartYear := prevItemTime.Year()
		showHeader = currStartMonth != prevStartMonth || currStartYear != prevStartYear
	}
	if position >= 1 {
		bookingItemTime := p.bookingItems[position-1].Time
		showHeader = showHeader &&
			(bookingItemTime.Month() != currStartMonth || bookingItemTime.Year() != currStartYear) &&
			position != 1
	}
	if showHeader {
		flex := layout.Flex{Axis: layout.Vertical}
		return flex.Layout(gtx,
			p.drawBookingItemHeader(gtx, currStartMonth, currStartYear)[0],
			p.drawBookingItemHeader(gtx, currStartMonth, currStartYear)[1],
			layout.Rigid(p.bookingItems[index].Layout),
		)
	}
	return p.bookingItems[index].Layout(gtx)
}

func (p *page) drawBookingItemHeader(_ fwk.Gtx, month time.Month, year int) []layout.FlexChild {
	dateHeader := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		inset := layout.UniformInset(16)
		return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			val := fmt.Sprintf("%s  %d", month, year)
			b := material.Body1(p.Theme, val)
			b.Font.Weight = text.ExtraBlack
			b.TextSize = unit.Sp(24)
			d := b.Layout(gtx)
			return d
		})
	})
	columnHeader := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		inset := layout.Inset{Left: 16, Right: 16, Top: 8, Bottom: 8}
		return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			flex := layout.Flex{}
			layoutRatio := flexUserLayoutRatio
			if p.loginUserResponse.IsAuthorized() {
				layoutRatio = flexAdminLayoutRatio
			}
			return flex.Layout(gtx,
				layout.Flexed(layoutRatio[0], func(gtx layout.Context) layout.Dimensions {
					title := i18n.Get(key.ID)
					b := material.Body1(p.Theme, title)
					b.Font.Weight = text.Black
					b.TextSize = unit.Sp(20)
					return b.Layout(gtx)
				}),
				layout.Flexed(layoutRatio[1], func(gtx layout.Context) layout.Dimensions {
					title := i18n.Get(key.Day)
					b := material.Body1(p.Theme, title)
					b.Font.Weight = text.Black
					b.TextSize = unit.Sp(20)
					return b.Layout(gtx)
				}),
				layout.Flexed(layoutRatio[2], func(gtx layout.Context) layout.Dimensions {
					title := i18n.Get(key.Weekday)
					b := material.Body1(p.Theme, title)
					b.Font.Weight = text.Black
					b.TextSize = unit.Sp(20)
					return b.Layout(gtx)
				}),
				layout.Flexed(layoutRatio[3], func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						title := i18n.Get(key.Availability)
						b := material.Body1(p.Theme, title)
						b.Font.Weight = text.Black
						b.TextSize = unit.Sp(20)
						return b.Layout(gtx)
					})
				}),
				layout.Flexed(layoutRatio[4], func(gtx layout.Context) layout.Dimensions {
					if !p.loginUserResponse.IsAuthorized() {
						return fwk.Dim{}
					}
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						title := i18n.Get(key.Delete)
						b := material.Body1(p.Theme, title)
						b.Font.Weight = text.Black
						b.TextSize = unit.Sp(20)
						return b.Layout(gtx)
					})
				}),
			)
		})
	})
	return []layout.FlexChild{dateHeader, columnHeader}
}

func (p *page) fetchBookings() {
	if !p.isFetchingBookings {
		p.isFetchingBookings = true
		p.Service().Bookings(service.BookingsRequest{
			StartDate: p.bookingForm.StartDate,
			EndDate:   p.bookingForm.EndDate,
		})
	}
}

func (p *page) OnServiceStateChange(event service.Event) {
	switch eventData := event.Data.(type) {
	case service.BookingsResponse:
		p.isFetchingBookings = false
		if eventData.Error != "" {
			errStr := eventData.Error
			if strings.Contains(errStr, "connection refused") {
				errStr = "connection refused"
			}
			p.Snackbar().Show(errStr, &p.closeSnapBar, color.NRGBA{R: 255, A: 255}, "Close")
			p.bookingItems = []*pageItem{}
			break
		}
		startDate := p.bookingForm.StartDate
		endDate := p.bookingForm.EndDate
		var bookingItems []*pageItem
		for startDate.Before(endDate) || startDate.Equal(endDate) {
			pg := &pageItem{
				Theme:             p.Theme,
				Manager:           p.Manager,
				Time:              startDate,
				LoginUserResponse: p.loginUserResponse,
				Booking: service.Booking{
					StartDate: startDate,
					EndDate:   startDate,
				},
				parentPage: p,
			}
			bookingItems = append(bookingItems, pg)
			startDate = startDate.AddDate(0, 0, 1)
		}
		for _, eachBooking := range eventData.Bookings {
			for _, bookingItem := range bookingItems {
				bookingTime := bookingItem.Time
				bkYear, bkMonth, bkDay := bookingTime.Date()
				startYear, startMonth, startDay := eachBooking.StartDate.Date()
				endYear, endMonth, endDay := eachBooking.EndDate.Date()
				bkDateStr := fmt.Sprintf("%d%d%d", bkYear, bkMonth, bkDay)
				startDateStr := fmt.Sprintf("%d%d%d", startYear, startMonth, startDay)
				endDateStr := fmt.Sprintf("%d%d%d", endYear, endMonth, endDay)
				isBooked := (bookingTime.After(eachBooking.StartDate) || bkDateStr == startDateStr) &&
					(bookingTime.Before(eachBooking.EndDate) || bkDateStr == endDateStr)
				if isBooked {
					bookingItem.Booking = eachBooking
				}
				bookingItem.LoginUserResponse = p.loginUserResponse
				bookingItem.parentPage = p
			}
		}

		p.bookingItems = bookingItems
		p.Window().Invalidate()
	case service.UserResponse:
		p.loginUserResponse = eventData
		for _, bookingItem := range p.bookingItems {
			bookingItem.LoginUserResponse = eventData
		}
		if !p.isFetchingBookings {
			p.fetchBookings()
		}
		p.Window().Invalidate()
	case service.DeleteBookingResponse:
		var txt string
		if (eventData.Error == "" && eventData.ID.Hex() == primitive.NilObjectID.Hex()) || !p.isDeletingBooking {
			return
		}
		if eventData.Error != "" {
			txt = fmt.Sprintf("couldn't delete booking with ID %s, error: %s", eventData.ID.Hex(), eventData.Error)
		}
		if eventData.Error == "" {
			txt = fmt.Sprintf("Successfully deleted booking with ID %s", eventData.ID.Hex())
		}
		if txt != "" {
			if !p.isFetchingBookings {
				p.fetchBookings()
			}
			p.Snackbar().Show(txt, nil, color.NRGBA(colornames.White), "CLOSE")
		}
		p.isDeletingBooking = false
		p.Window().Invalidate()
	}
}
func (p *page) URL() fwk.URL {
	return fwk.BookingsPageURL
}
