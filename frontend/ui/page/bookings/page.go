package bookings

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/page/add_edit_booking"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
	"strings"
	"time"
)

type page struct {
	layout.List
	fwk.Manager
	Theme              *material.Theme
	title              string
	buttonNavigation   widget.Clickable
	btnBackdrop        widget.Clickable
	btnAddBooking      widget.Clickable
	btnMenuContent     widget.Clickable
	btnYes             widget.Clickable
	btnNo              widget.Clickable
	btnMenuIcon        widget.Clickable
	menuIcon           *widget.Icon
	closeIcon          *widget.Icon
	menuVisibilityAnim component.VisibilityAnimation
	navigationIcon     *widget.Icon
	bookingItems       []*pageItem
	isFetchingBookings bool
	isDeletingBooking  bool
	initialized        bool
	subscription       service.Subscriber
	fetchingBookingsCh chan service.BookingsResponse
	bookingsCount      int64
	bookingForm        view.BookingForm
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
	iconMenu, _ := widget.NewIcon(icons.NavigationMoreVert)
	errorTh := *manager.Theme()
	errorTh.ContrastBg = color.NRGBA(colornames.Red500)
	startDate := utils.GetFirstDayOfMonth(time.Now().Local())
	endDate := utils.GetLastDayOfMonth(time.Now().Local().AddDate(0, 5, 0))
	theme := *manager.Theme()
	p := page{
		Manager:            manager,
		Theme:              &theme,
		title:              "Bookings",
		navigationIcon:     navIcon,
		List:               layout.List{Axis: layout.Vertical},
		bookingItems:       []*pageItem{},
		menuIcon:           iconMenu,
		closeIcon:          closeIcon,
		fetchingBookingsCh: make(chan service.BookingsResponse, 1),
		menuVisibilityAnim: component.VisibilityAnimation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		},
		limit:       1000,
		offset:      0,
		bookingForm: view.NewBookingForm(manager, service.Booking{StartDate: startDate, EndDate: endDate}, false),
	}
	p.subscription = manager.Service().Subscribe(service.TopicBookingsFetched, service.TopicUserLoggedInOut, service.TopicDeleteBooking)
	p.subscription.SubscribeWithCallback(p.OnServiceStateChange)
	return &p
}

func (p *page) Layout(gtx fwk.Gtx) fwk.Dim {
	if !p.initialized {
		if p.Theme == nil {
			p.Theme = p.Manager.Theme()
		}
		p.fetchBookings()
		p.initialized = true
	}

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
		layout.Rigid(p.bookingForm.Layout),
		layout.Rigid(p.drawBookingItems),
	)
	p.drawMenuLayout(gtx)
	return d
}

func (p *page) DrawAppBar(gtx fwk.Gtx) fwk.Dim {
	gtx.Constraints.Max.Y = gtx.Dp(56)
	th := p.Theme
	if p.btnMenuIcon.Clicked() {
		p.menuVisibilityAnim.Appear(gtx.Now)
	}
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
			layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
				if !p.loginUserResponse.IsLoggedIn() || !p.loginUserResponse.IsAdmin() {
					return fwk.Dim{}
				}
				button := material.IconButton(th, &p.btnMenuIcon, p.menuIcon, "Context Menu")
				button.Size = unit.Dp(40)
				button.Background = th.Palette.ContrastBg
				button.Color = th.Palette.ContrastFg
				button.Inset = layout.UniformInset(unit.Dp(8))
				d := button.Layout(gtx)
				return d
			}),
		)
	})
}

func (p *page) drawMenuLayout(gtx fwk.Gtx) fwk.Dim {
	if p.btnBackdrop.Clicked() {
		if !p.btnMenuContent.Pressed() {
			p.menuVisibilityAnim.Disappear(gtx.Now)
		}
	}
	layout.Stack{Alignment: layout.NE}.Layout(gtx,
		layout.Stacked(func(gtx fwk.Gtx) fwk.Dim {
			return p.btnBackdrop.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					progress := p.menuVisibilityAnim.Revealed(gtx)
					gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * progress)
					gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * progress)
					return component.Rect{Size: gtx.Constraints.Max, Color: color.NRGBA{A: 200}}.Layout(gtx)
				},
			)
		}),
		layout.Stacked(func(gtx fwk.Gtx) fwk.Dim {
			progress := p.menuVisibilityAnim.Revealed(gtx)
			macro := op.Record(gtx.Ops)
			d := p.btnMenuContent.Layout(gtx, p.drawMenuItems)
			call := macro.Stop()
			d.Size.X = int(float32(d.Size.X) * progress)
			d.Size.Y = int(float32(d.Size.Y) * progress)
			component.Rect{Size: d.Size, Color: color.NRGBA(colornames.White)}.Layout(gtx)
			clipOp := clip.Rect{Max: d.Size}.Push(gtx.Ops)
			call.Add(gtx.Ops)
			clipOp.Pop()
			return d
		}),
	)
	return fwk.Dim{}
}

func (p *page) drawMenuItems(gtx fwk.Gtx) fwk.Dim {
	gtx.Constraints.Max.X = 200
	if p.btnAddBooking.Clicked() {
		p.menuVisibilityAnim.Disappear(gtx.Now)
		p.Manager.NavigateToPage(add_edit_booking.New(p.Manager, sqlc.Booking{
			ID: 0,
		}), nil)
	}
	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx,
		p.drawMenuItem("New Booking", &p.btnAddBooking),
	)
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

func (p *page) drawBookingItems(gtx fwk.Gtx) fwk.Dim {
	if len(p.bookingItems) == 0 {
		return view.Dim{}
	}
	return p.List.Layout(gtx, len(p.bookingItems), func(gtx fwk.Gtx, index int) (d fwk.Dim) {
		showHeader := true
		currItemTime := p.bookingItems[index].Time
		currStartMonth := currItemTime.Month()
		currStartYear := currItemTime.Year()
		if index > 0 {
			prevItemTime := p.bookingItems[index-1].Time
			prevStartMonth := prevItemTime.Month()
			prevStartYear := prevItemTime.Year()
			showHeader = currStartMonth != prevStartMonth || currStartYear != prevStartYear
		}
		if showHeader {
			flex := layout.Flex{Axis: layout.Vertical}
			return flex.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					inset := layout.UniformInset(16)
					return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						val := fmt.Sprintf("%s  %d", currStartMonth, currStartYear)
						b := material.Body1(p.Theme, val)
						b.Font.Weight = text.ExtraBlack
						b.TextSize = unit.Sp(24)
						d := b.Layout(gtx)
						return d
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					inset := layout.Inset{Left: 16, Right: 16, Top: 8, Bottom: 8}
					return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						flex := layout.Flex{}
						layoutRatio := flexUserLayoutRatio
						isAuthorized := p.loginUserResponse.IsLoggedIn() && p.loginUserResponse.IsAdmin()
						if isAuthorized {
							layoutRatio = flexAdminLayoutRatio
						}
						return flex.Layout(gtx,
							layout.Flexed(layoutRatio[0], func(gtx layout.Context) layout.Dimensions {
								b := material.Body1(p.Theme, "ID")
								b.Font.Weight = text.Black
								b.TextSize = unit.Sp(20)
								return b.Layout(gtx)
							}),
							layout.Flexed(layoutRatio[1], func(gtx layout.Context) layout.Dimensions {
								b := material.Body1(p.Theme, "Day")
								b.Font.Weight = text.Black
								b.TextSize = unit.Sp(20)
								return b.Layout(gtx)
							}),
							layout.Flexed(layoutRatio[2], func(gtx layout.Context) layout.Dimensions {
								b := material.Body1(p.Theme, "Weekday")
								b.Font.Weight = text.Black
								b.TextSize = unit.Sp(20)
								return b.Layout(gtx)
							}),
							layout.Flexed(layoutRatio[3], func(gtx layout.Context) layout.Dimensions {
								return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									b := material.Body1(p.Theme, "Availability")
									b.Font.Weight = text.Black
									b.TextSize = unit.Sp(20)
									return b.Layout(gtx)
								})
							}),
							layout.Flexed(layoutRatio[4], func(gtx layout.Context) layout.Dimensions {
								if !isAuthorized {
									return fwk.Dim{}
								}
								return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									b := material.Body1(p.Theme, "Delete")
									b.Font.Weight = text.Black
									b.TextSize = unit.Sp(20)
									return b.Layout(gtx)
								})
							}),
						)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.bookingItems[index].Layout(gtx)
				}),
			)
		}
		return p.bookingItems[index].Layout(gtx)
	})
}

func (p *page) fetchBookings() {
	if !p.isFetchingBookings {
		p.isFetchingBookings = true
		p.Service().Bookings(service.BookingParams{
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
				Booking: sqlc.Booking{
					ID:        0,
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
	case service.UserResponse:
		p.loginUserResponse = eventData
		for _, bookingItem := range p.bookingItems {
			bookingItem.LoginUserResponse = eventData
		}
		if !p.isFetchingBookings {
			p.fetchBookings()
		}
	case service.DeleteBookingResponse:
		var txt string
		if (eventData.Error == "" && eventData.ID == 0) || !p.isDeletingBooking {
			return
		}
		if eventData.Error != "" {
			txt = fmt.Sprintf("couldn't delete booking with ID %d, error: %s", eventData.ID, eventData.Error)
		}
		if eventData.Error == "" {
			txt = fmt.Sprintf("Successfully deleted booking with ID %d", eventData.ID)
		}
		if txt != "" {
			if !p.isFetchingBookings {
				p.fetchBookings()
			}
			p.Snackbar().Show(txt, nil, color.NRGBA(colornames.White), "CLOSE")
		}
		p.isDeletingBooking = false
	}
}
func (p *page) URL() fwk.URL {
	return fwk.BookingsPageURL
}
