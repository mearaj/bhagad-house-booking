package search_bookings

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
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
	limit              int32
	offset             int32
	loginUserResponse  service.UserResponse
	closeSnapBar       widget.Clickable
	formField          view.FormField
	query              string
}

var flexUserLayoutRatio = [5]float32{0.16, 0.28, 0.28, 0.28, 0}
var flexAdminLayoutRatio = [5]float32{0.16, 0.21, 0.21, 0.21, 0.21}

func New(manager fwk.Manager) fwk.Page {
	navIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	closeIcon, _ := widget.NewIcon(icons.ContentClear)
	errorTh := *user.Theme()
	errorTh.ContrastBg = color.NRGBA(colornames.Red500)
	theme := *user.Theme()
	p := page{
		Manager:            manager,
		Theme:              &theme,
		navigationIcon:     navIcon,
		ViewList:           layout.List{Axis: layout.Vertical},
		bookingItems:       []*pageItem{},
		closeIcon:          closeIcon,
		fetchingBookingsCh: make(chan service.BookingsResponse, 1),
		limit:              1000,
		offset:             0,
		formField:          view.FormField{FieldName: i18n.Get(key.SearchBookings)},
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
		p.fetchBookings()
		p.initialized = true
	}
	val := strings.TrimSpace(p.formField.Text())
	if val != p.query && !p.isFetchingBookings {
		p.query = p.formField.Text()
		p.fetchBookings()
	}
	p.title = i18n.Get(key.SearchBookings)

	flex := layout.Flex{Axis: layout.Vertical}
	d := flex.Layout(gtx,
		layout.Rigid(p.DrawAppBar),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return p.ViewList.Layout(gtx, len(p.bookingItems)+1, func(gtx layout.Context, index int) layout.Dimensions {
				if index == 0 {
					return p.drawFormField(gtx)
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

func (p *page) drawFormField(gtx fwk.Gtx) fwk.Dim {
	inset := layout.UniformInset(unit.Dp(16))
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		flex := layout.Flex{Axis: layout.Vertical}
		return flex.Layout(gtx,
			layout.Rigid(layout.Spacer{Height: 8}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return view.DrawFormFieldRowWithLabel(gtx, p.Theme, "", p.formField.FieldName, &p.formField.TextField, nil)
			}),
			layout.Rigid(layout.Spacer{Height: 16}.Layout),
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
			isAuthorized := p.loginUserResponse.IsLoggedIn() && p.loginUserResponse.IsAdmin()
			if isAuthorized {
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
					if !isAuthorized {
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
		p.Service().SearchBookings(p.formField.Text())
	}
}

func (p *page) OnServiceStateChange(event service.Event) {
	switch eventData := event.Data.(type) {
	case service.UpdateBookingResponse:
		if !p.isFetchingBookings {
			p.fetchBookings()
		}
	case service.SearchBookingsResponse:
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
		var bookingItems []*pageItem
		for _, eachBooking := range eventData.Bookings {
			bookingItem := &pageItem{
				Theme:             p.Theme,
				Manager:           p.Manager,
				Booking:           eachBooking,
				LoginUserResponse: p.loginUserResponse,
				parentPage:        p,
			}
			bookingItems = append(bookingItems, bookingItem)
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
	return fwk.SearchBookingsPageURL
}
