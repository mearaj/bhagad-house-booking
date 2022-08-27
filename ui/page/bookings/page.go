package bookings

import (
	"fmt"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/service"
	. "github.com/mearaj/bhagad-house-booking/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/ui/page/add_edit_booking"
	"github.com/mearaj/bhagad-house-booking/ui/view"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
	"time"
)

type page struct {
	layout.List
	Manager
	Theme                   *material.Theme
	title                   string
	iconNewChat             *widget.Icon
	btnBackdrop             widget.Clickable
	buttonNavigation        widget.Clickable
	btnMenuIcon             widget.Clickable
	btnMenuContent          widget.Clickable
	btnAddBooking           widget.Clickable
	btnDeleteBookings       widget.Clickable
	btnCloseSelection       widget.Clickable
	btnYes                  widget.Clickable
	btnNo                   widget.Clickable
	btnSelectAll            widget.Clickable
	btnDeleteAll            widget.Clickable
	btnSelectionMode        widget.Clickable
	menuIcon                *widget.Icon
	closeIcon               *widget.Icon
	menuVisibilityAnim      component.VisibilityAnimation
	navigationIcon          *widget.Icon
	bookingsView            []*pageItem
	NoBooking               View
	ModalContent            *view.ModalContent
	SelectionMode           bool
	isFetchingBookings      bool
	isFetchingBookingsCount bool
	initialized             bool
	subscription            service.Subscriber
	fetchingBookingsCh      chan []service.Booking
	fetchingBookingsCountCh chan int64
	bookingsCount           int64
}

func New(manager Manager) Page {
	navIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	closeIcon, _ := widget.NewIcon(icons.ContentClear)
	iconNewChat, _ := widget.NewIcon(icons.ContentCreate)
	iconMenu, _ := widget.NewIcon(icons.NavigationMoreVert)
	errorTh := *manager.Theme()
	errorTh.ContrastBg = color.NRGBA(colornames.Red500)
	theme := *manager.Theme()
	p := page{
		Manager:                 manager,
		Theme:                   &theme,
		title:                   "Bookings",
		navigationIcon:          navIcon,
		iconNewChat:             iconNewChat,
		List:                    layout.List{Axis: layout.Vertical},
		bookingsView:            []*pageItem{},
		menuIcon:                iconMenu,
		closeIcon:               closeIcon,
		fetchingBookingsCh:      make(chan []service.Booking, 10),
		fetchingBookingsCountCh: make(chan int64, 10),
		menuVisibilityAnim: component.VisibilityAnimation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		},
	}
	p.ModalContent = view.NewModalContent(func() {
		p.Modal().Dismiss(nil)
	})
	p.NoBooking = view.NewNoBooking(manager)
	p.subscription = manager.Service().Subscribe(service.BookingsChangedEventTopic)
	return &p
}

func (p *page) Layout(gtx Gtx) Dim {
	if !p.initialized {
		if p.Theme == nil {
			p.Theme = p.Manager.Theme()
		}
		p.fetchBookings()
		p.fetchBookingsCount()
		p.initialized = true
	}

	p.listenToFetchBookings()
	p.listenToFetchBookingsCount()
	p.handleSelectionMode()
	p.handleAddBookingClick(gtx)

	flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEnd, Alignment: layout.Start}

	d := flex.Layout(gtx,
		layout.Rigid(p.DrawAppBar),
		layout.Rigid(p.drawIdentitiesItems),
	)
	p.drawMenuLayout(gtx)
	p.handleEvents(gtx)
	return d
}

func (p *page) DrawAppBar(gtx Gtx) Dim {
	gtx.Constraints.Max.Y = gtx.Dp(56)
	if p.btnMenuIcon.Clicked() {
		p.menuVisibilityAnim.Appear(gtx.Now)
	}
	if p.SelectionMode {
		return p.DrawSelectionAppBar(gtx)
	}
	return p.DrawNormalAppBar(gtx)
}
func (p *page) DrawNormalAppBar(gtx Gtx) Dim {
	gtx.Constraints.Max.Y = gtx.Dp(56)
	th := p.Theme
	if p.buttonNavigation.Clicked() {
		p.PopUp()
	}

	return view.DrawAppBarLayout(gtx, th, func(gtx Gtx) Dim {
		return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx Gtx) Dim {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx Gtx) Dim {
						navigationIcon := p.navigationIcon
						button := material.IconButton(th, &p.buttonNavigation, navigationIcon, "Nav Icon Button")
						button.Size = unit.Dp(40)
						button.Background = th.Palette.ContrastBg
						button.Color = th.Palette.ContrastFg
						button.Inset = layout.UniformInset(unit.Dp(8))
						return button.Layout(gtx)
					}),
					layout.Rigid(func(gtx Gtx) Dim {
						return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx Gtx) Dim {
							titleText := p.title
							title := material.Body1(th, titleText)
							title.Color = th.Palette.ContrastFg
							title.TextSize = unit.Sp(18)
							return title.Layout(gtx)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx Gtx) Dim {
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
func (p *page) DrawSelectionAppBar(gtx Gtx) Dim {
	gtx.Constraints.Max.Y = gtx.Dp(56)
	th := p.Theme
	if p.btnCloseSelection.Clicked() {
		p.clearAllSelection()
		p.menuVisibilityAnim.Disappear(gtx.Now)
	}
	return view.DrawAppBarLayout(gtx, th, func(gtx Gtx) Dim {
		return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx Gtx) Dim {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx Gtx) Dim {
						closeIcon := p.closeIcon
						button := material.IconButton(th, &p.btnCloseSelection, closeIcon, "Nav Icon Button")
						button.Size = unit.Dp(40)
						button.Background = th.Palette.ContrastBg
						button.Color = th.Palette.ContrastFg
						button.Inset = layout.UniformInset(unit.Dp(8))
						return button.Layout(gtx)
					}),
					layout.Rigid(func(gtx Gtx) Dim {
						return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx Gtx) Dim {
							var txt string
							count := p.getSelectionCount()
							if count == 0 {
								txt = "None Selected"
							} else {
								txt = fmt.Sprintf("(%d) Selected", count)
							}
							title := material.Body1(th, txt)
							title.Color = th.Palette.ContrastFg
							title.TextSize = unit.Sp(18)
							return title.Layout(gtx)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx Gtx) Dim {
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

func (p *page) drawIdentitiesItems(gtx Gtx) Dim {
	if len(p.bookingsView) == 0 {
		return p.NoBooking.Layout(gtx)
	}
	return p.List.Layout(gtx, len(p.bookingsView), func(gtx Gtx, index int) (d Dim) {
		inset := layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(0)}
		return inset.Layout(gtx, func(gtx Gtx) Dim {
			return p.bookingsView[index].Layout(gtx)
		})
	})
}

func (p *page) drawMenuLayout(gtx Gtx) Dim {
	if p.btnBackdrop.Clicked() {
		if !p.btnMenuContent.Pressed() {
			p.menuVisibilityAnim.Disappear(gtx.Now)
		}
		for _, idView := range p.bookingsView {
			if !idView.btnMenuContent.Pressed() && !idView.Hovered() {
				idView.menuVisibilityAnim.Disappear(gtx.Now)
			}
		}
	}
	layout.Stack{Alignment: layout.NE}.Layout(gtx,
		layout.Stacked(func(gtx Gtx) Dim {
			return p.btnBackdrop.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					progress := p.menuVisibilityAnim.Revealed(gtx)
					gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * progress)
					gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * progress)
					return component.Rect{Size: gtx.Constraints.Max, Color: color.NRGBA{A: 200}}.Layout(gtx)
				},
			)
		}),
		layout.Stacked(func(gtx Gtx) Dim {
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
	return Dim{}
}

func (p *page) drawMenuItems(gtx Gtx) Dim {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) / 1.5)
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	if p.SelectionMode {
		return p.drawSelectionMenuItems(gtx)
	}
	return p.drawNormalMenuItems(gtx)
}

func (p *page) drawNormalMenuItems(gtx Gtx) Dim {
	if p.btnSelectAll.Clicked() {
		p.selectAll()
		p.menuVisibilityAnim.Disappear(gtx.Now)
	}
	if p.btnSelectionMode.Clicked() {
		p.SelectionMode = true
		p.menuVisibilityAnim.Disappear(gtx.Now)
	}
	if p.btnDeleteAll.Clicked() {
		p.selectAll()
		p.Modal().Show(p.drawDeleteBookingsModal, nil, Animation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		})
		p.menuVisibilityAnim.Disappear(gtx.Now)
	}

	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx,
		p.drawMenuItem("Add Booking", &p.btnAddBooking),
		p.drawMenuItem("Selection Mode", &p.btnSelectionMode),
		p.drawMenuItem("Select All Bookings", &p.btnSelectAll),
		p.drawMenuItem("Delete All Bookings", &p.btnDeleteAll),
	)
}
func (p *page) drawSelectionMenuItems(gtx Gtx) Dim {
	if p.btnDeleteBookings.Clicked() {
		p.Modal().Show(p.drawDeleteBookingsModal, nil, Animation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		})
		p.menuVisibilityAnim.Disappear(gtx.Now)
	}

	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx,
		p.drawMenuItem("Delete Selected Bookings", &p.btnDeleteBookings),
		p.drawMenuItem("Clear Selection", &p.btnCloseSelection),
	)
}
func (p *page) drawMenuItem(txt string, btn *widget.Clickable) layout.FlexChild {
	inset := layout.UniformInset(unit.Dp(12))
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		btnStyle := material.ButtonLayoutStyle{Button: btn}
		btnStyle.Background = color.NRGBA(colornames.White)
		return btnStyle.Layout(gtx,
			func(gtx Gtx) Dim {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				inset := inset
				return inset.Layout(gtx, func(gtx Gtx) Dim {
					return layout.Flex{Spacing: layout.SpaceEnd}.Layout(gtx,
						layout.Rigid(func(gtx Gtx) Dim {
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

func (p *page) drawDeleteBookingsModal(gtx Gtx) Dim {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.85)
	if p.btnYes.Clicked() {
		bookings := make([]service.Booking, 0)
		bookingsViewSize := len(p.bookingsView)
		for _, eachView := range p.bookingsView {
			if eachView.Selected {
				bookings = append(bookings, eachView.Booking)
			}
		}
		<-p.Service().DeleteBookings(bookings)
		p.Modal().Dismiss(func() {
			p.Window().Invalidate()
			p.clearAllSelection()
			var txtTmp string
			if len(bookings) == bookingsViewSize {
				txtTmp = "all bookings."
			} else {
				txtTmp = fmt.Sprintf("%d bookings.", len(bookings))
			}
			if len(bookings) == 1 {
				txtTmp = "1 booking."
			}
			txt := fmt.Sprintf("Successfully deleted %s", txtTmp)
			p.Snackbar().Show(txt, nil, color.NRGBA{}, "")
		})
	}
	if p.btnNo.Clicked() {
		p.Modal().Dismiss(func() {
			p.clearAllSelection()
		})
	}
	count := p.getSelectionCount()
	bookingsSize := len(p.bookingsView)
	var txt string
	if count == bookingsSize {
		txt = "all bookings"
	} else {
		txt = fmt.Sprintf("%d selected bookings", count)
	}
	if count == 1 {
		txt = "the selected booking"
	}
	promptContent := view.NewPromptContent(p.Theme,
		"Booking Deletion!",
		fmt.Sprintf("Are you sure you want to delete %s?", txt),
		&p.btnYes, &p.btnNo)
	return p.ModalContent.DrawContent(gtx, p.Theme, promptContent.Layout)
}

func (p *page) onSuccess() {
	p.Modal().Dismiss(func() {
		a := p.Service().Booking()
		txt := fmt.Sprintf("Successfully created %d", a.ID)
		p.Window().Invalidate()
		p.Snackbar().Show(txt, nil, color.NRGBA{}, "")
	})
}
func (p *page) getSelectionCount() (count int) {
	for _, item := range p.bookingsView {
		if item.Selected {
			count++
		}
	}
	return count
}
func (p *page) clearAllSelection() {
	p.SelectionMode = false
	for _, item := range p.bookingsView {
		item.Selected = false
		item.SelectionMode = false
	}
}
func (p *page) selectAll() {
	p.SelectionMode = true
	for _, item := range p.bookingsView {
		item.Selected = true
		item.SelectionMode = true
	}
}

func (p *page) fetchBookings() {
	if !p.isFetchingBookings {
		p.isFetchingBookings = true
		go func() {
			p.fetchingBookingsCh <- <-p.Service().Bookings()
			p.Window().Invalidate()
		}()
	}
}
func (p *page) fetchBookingsCount() {
	if !p.isFetchingBookingsCount {
		p.isFetchingBookingsCount = true
		go func() {
			p.fetchingBookingsCountCh <- <-p.Service().BookingsCount()
			p.Window().Invalidate()
		}()
	}
}

func (p *page) listenToFetchBookings() {
	shouldBreak := false
	for {
		select {
		case bookings := <-p.fetchingBookingsCh:
			bookingViews := make([]*pageItem, len(bookings))
			for i, eachCustomer := range bookings {
				bookingViews[i] = &pageItem{
					Theme:        p.Theme,
					Manager:      p.Manager,
					Booking:      eachCustomer,
					ModalContent: p.ModalContent,
				}
			}
			//pos := p.Position.First
			p.bookingsView = bookingViews
			//p.Position.First = pos + len(bookings)
			p.isFetchingBookings = false
		default:
			shouldBreak = true
		}
		if shouldBreak {
			break
		}
	}

}

func (p *page) listenToFetchBookingsCount() {
	shouldBreak := false
	for {
		select {
		case bookingsCount := <-p.fetchingBookingsCountCh:
			if bookingsCount != p.bookingsCount {
				p.bookingsCount = bookingsCount
				if !p.isFetchingBookings {
					p.fetchBookings()
				}
			}
			p.isFetchingBookingsCount = false
		default:
			shouldBreak = true
		}
		if shouldBreak {
			break
		}
	}
}

func (p *page) handleSelectionMode() {
	for _, item := range p.bookingsView {
		if p.SelectionMode {
			item.SelectionMode = p.SelectionMode
		} else {
			if item.SelectionMode {
				p.SelectionMode = item.SelectionMode
				break
			}
		}
	}
	if p.SelectionMode {
		p.Theme.ContrastBg = color.NRGBA{A: 255}
	} else {
		p.Theme.ContrastBg = p.Manager.Theme().ContrastBg
	}
}

func (p *page) handleAddBookingClick(gtx Gtx) {
	if p.btnAddBooking.Clicked() {
		addEditBookingPage := add_edit_booking.New(p.Manager, service.Booking{})
		p.Manager.NavigateToPage(addEditBookingPage, func() {
			p.menuVisibilityAnim.Disappear(gtx.Now)
		})
	}
}

func (p *page) handleEvents(gtx Gtx) {
	for _, e := range gtx.Queue.Events(p) {
		switch e := e.(type) {
		case pointer.Event:
			switch e.Type {
			case pointer.Press:
				if !p.btnMenuContent.Pressed() {
					p.menuVisibilityAnim.Disappear(gtx.Now)
				}
				for _, idView := range p.bookingsView {
					if !idView.btnMenuContent.Pressed() && !idView.Hovered() {
						idView.menuVisibilityAnim.Disappear(gtx.Now)
					}
				}
			}
		}
	}
}

func (p *page) OnDatabaseChange(event service.Event) {
	switch e := event.Data.(type) {
	case service.BookingsChangedEventData:
		_ = e
		p.fetchBookings()
		p.fetchBookingsCount()
	}
}
func (p *page) URL() URL {
	return BookingsPageURL
}
