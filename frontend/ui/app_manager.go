package ui

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"gioui.org/x/notify"
	"github.com/mearaj/bhagad-house-booking/common/alog"
	"github.com/mearaj/bhagad-house-booking/common/assets/fonts"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	. "github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/page/about"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/page/add_edit_booking"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/page/bookings"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/page/help"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/page/notifications"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/page/settings"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/page/theme"
	view2 "github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"image"
	"strconv"
	"strings"
	"sync"
	"time"
)

// AppManager Always call NewAppManager function to create AppManager instance
type AppManager struct {
	window *app.Window
	view2.Greetings
	settingsSideBar Page
	theme           *material.Theme
	service         service.Service
	Constraints     layout.Constraints
	Metric          unit.Metric
	notifier        notify.Notifier
	system.Insets
	// isStageRunning, true value indicates app is running in foreground,
	// false indicates running in background
	isStageRunning   bool
	modal            Modal
	pagesStack       []Page
	prePushedView    Page
	poppedUpView     Page
	pageAnimation    Animation
	afterPageChange  func()
	snackbar         Snackbar
	initialized      bool
	initializedMutex sync.RWMutex
}

func (m *AppManager) Theme() *material.Theme {
	return m.theme
}

func (m *AppManager) Service() service.Service {
	return m.service
}
func (m *AppManager) SystemInsets() system.Insets {
	return m.Insets
}

func (m *AppManager) Window() *app.Window {
	return m.window
}

func (m *AppManager) Notifier() notify.Notifier {
	return m.notifier
}

func (m *AppManager) Snackbar() Snackbar {
	return m.snackbar
}

func (m *AppManager) Initialized() bool {
	m.initializedMutex.RLock()
	initialized := m.initialized
	m.initializedMutex.RUnlock()
	return initialized
}
func (m *AppManager) setInitialized(initialized bool) {
	m.initializedMutex.Lock()
	m.initialized = initialized
	m.initializedMutex.Unlock()
}

func (m *AppManager) init() {
	m.theme = fonts.NewTheme()
	settingsPage := settings.New(m)
	m.settingsSideBar = settingsPage
	var err error
	m.notifier, err = notify.NewNotifier()
	if err != nil {
		alog.Logger().Errorln(err)
	}
	m.SetModal(view2.NewModalStack())
	m.snackbar = view2.NewSnackBar(m)
	m.pageAnimation = component.VisibilityAnimation{
		Duration: time.Millisecond * 250,
		State:    component.Invisible,
		Started:  time.Time{},
	}
	m.setInitialized(true)
}

func (m *AppManager) Layout(gtx Gtx) Dim {
	d := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx Gtx) Dim {
			return component.Rect{
				Color: m.Theme().ContrastBg,
				Size:  image.Point{X: gtx.Constraints.Max.X, Y: gtx.Dp(m.Insets.Top)},
				Radii: 0,
			}.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx Gtx) Dim {
			size := image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y - gtx.Dp(m.Insets.Bottom)}
			bounds := image.Rectangle{Max: size}
			paint.FillShape(gtx.Ops, m.Theme().Bg, clip.UniformRRect(bounds, 0).Op(gtx.Ops))
			d := m.drawPage(gtx)
			m.Snackbar().Layout(gtx)
			m.Modal().Layout(gtx)
			return d
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return component.Rect{
				Color: m.Theme().ContrastBg,
				Size:  image.Point{X: gtx.Constraints.Max.X, Y: gtx.Dp(m.Insets.Bottom)},
				Radii: 0,
			}.Layout(gtx)
		}),
	)
	return d
}

func (m *AppManager) CurrentPage() Page {
	stackSize := len(m.pagesStack)
	if stackSize > 0 {
		return m.pagesStack[stackSize-1]
	}
	m.pagesStack = []Page{settings.New(m)}
	return m.settingsSideBar
}
func (m *AppManager) GetWindowWidthInDp() int {
	width := int(float32(m.Constraints.Max.X) / m.Metric.PxPerDp)
	return width
}

func (m *AppManager) GetWindowWidthInPx() int {
	return m.Constraints.Max.X
}

func (m *AppManager) GetWindowHeightInDp() int {
	width := int(float32(m.Constraints.Max.Y) / m.Metric.PxPerDp)
	return width
}

func (m *AppManager) GetWindowHeightInPx() int {
	return m.Constraints.Max.Y
}

func (m *AppManager) IsStageRunning() bool {
	return m.isStageRunning
}
func (m *AppManager) ShouldDrawSidebar() bool {
	minWidth := 800 // 800 is value in Dp
	winWidth := m.GetWindowWidthInDp()
	return winWidth >= minWidth
}

func (m *AppManager) Modal() Modal {
	return m.modal
}
func (m *AppManager) SetModal(modal Modal) {
	m.modal = modal
}

func (m *AppManager) PageFromUrl(url URL) Page {
	switch url {
	case SettingsPageURL:
		return m.settingsSideBar
	default:
		return m.settingsSideBar
	}
}

func (m *AppManager) drawPage(gtx Gtx) Dim {
	maxDim := gtx.Constraints.Max

	d := layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) (d Dim) {
			if !m.ShouldDrawSidebar() {
				return d
			}
			gtx.Constraints.Max.X = int(float32(maxDim.X) * 0.40)
			gtx.Constraints.Min = gtx.Constraints.Max
			d = m.settingsSideBar.Layout(gtx)
			return d
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if m.ShouldDrawSidebar() {
				gtx.Constraints.Max.X = int(float32(maxDim.X) * 0.60)
			}
			gtx.Constraints.Min = gtx.Constraints.Max
			maxDim := gtx.Constraints.Max
			gtx.Constraints.Min = maxDim
			currentPage := m.CurrentPage()
			prePushedPage := m.prePushedView
			poppedUpPage := m.poppedUpView
			switch {
			case prePushedPage != nil:
				progress := m.pageAnimation.Revealed(gtx)
				if m.pageAnimation.Animating() || m.pageAnimation.State == component.Invisible {
					m.pageAnimation.Appear(gtx.Now)
				}
				if progress == 1 {
					m.prePushedView = nil
					m.pageAnimation.State = component.Invisible
				}

				cl := clip.Rect{Max: maxDim}.Push(gtx.Ops)
				layout.Stack{Alignment: layout.NW}.Layout(gtx,
					layout.Stacked(func(gtx Gtx) Dim {
						gtx.Constraints.Min = maxDim
						offsetOperation := op.Offset(image.Point{
							X: -int(float32(maxDim.X) * progress),
							Y: 0,
						}).Push(gtx.Ops)
						defer offsetOperation.Pop()
						if prePushedPage.URL() == SettingsPageURL && m.ShouldDrawSidebar() {
							m.Greetings.Layout(gtx)
							return Dim{Size: maxDim}
						}
						prePushedPage.Layout(gtx)
						return Dim{Size: maxDim}
					}),
					layout.Stacked(func(gtx Gtx) Dim {
						gtx.Constraints.Min = maxDim
						offsetOperation := op.Offset(image.Point{
							X: int(float32(maxDim.X) * (1 - progress)),
							Y: 0,
						}).Push(gtx.Ops)
						defer offsetOperation.Pop()
						if currentPage.URL() == SettingsPageURL && m.ShouldDrawSidebar() {
							m.Greetings.Layout(gtx)
							return Dim{Size: maxDim}
						}
						currentPage.Layout(gtx)
						return Dim{Size: maxDim}
					}),
				)
				cl.Pop()
				return Dim{Size: maxDim}
			case poppedUpPage != nil:
				progress := m.pageAnimation.Revealed(gtx)
				if m.pageAnimation.Animating() || m.pageAnimation.State == component.Invisible {
					m.pageAnimation.Appear(gtx.Now)
				}
				if progress == 1 {
					m.poppedUpView = nil
					m.pageAnimation.State = component.Invisible
				}

				cl := clip.Rect{Max: maxDim}.Push(gtx.Ops)
				layout.Stack{Alignment: layout.NW}.Layout(gtx,
					layout.Stacked(func(gtx Gtx) Dim {
						offsetOperation := op.Offset(image.Point{
							X: -int(float32(maxDim.X) * (1 - progress)),
							Y: 0,
						}).Push(gtx.Ops)
						defer offsetOperation.Pop()
						gtx.Constraints.Min = maxDim
						if currentPage.URL() == SettingsPageURL && m.ShouldDrawSidebar() {
							m.Greetings.Layout(gtx)
							return Dim{Size: maxDim}
						}
						currentPage.Layout(gtx)
						return Dim{Size: maxDim}
					}),
					layout.Stacked(func(gtx Gtx) Dim {
						offsetOperation := op.Offset(image.Point{
							X: int(float32(maxDim.X) * progress),
							Y: 0,
						}).Push(gtx.Ops)
						defer offsetOperation.Pop()
						gtx.Constraints.Min = maxDim
						if poppedUpPage.URL() == SettingsPageURL && m.ShouldDrawSidebar() {
							m.Greetings.Layout(gtx)
							return Dim{Size: maxDim}
						}
						poppedUpPage.Layout(gtx)
						return Dim{Size: maxDim}
					}),
				)
				cl.Pop()
				return Dim{Size: maxDim}
			default:

			}
			if currentPage.URL() == SettingsPageURL && m.ShouldDrawSidebar() {
				m.Greetings.Layout(gtx)
				return Dim{Size: maxDim}
			}
			return currentPage.Layout(gtx)
		}),
	)
	state := m.pageAnimation.State
	progress := m.pageAnimation.Revealed(gtx)
	shouldCall := state == component.Invisible && progress == 0 && m.afterPageChange != nil &&
		!m.pageAnimation.Animating()
	if shouldCall {
		m.afterPageChange()
		m.afterPageChange = nil
	}
	return d
}

func (m *AppManager) NavigateToPage(page Page, AfterNavCallback func()) {
	pageURL := page.URL()
	m.afterPageChange = AfterNavCallback
	if pageURL == m.CurrentPage().URL() {
		if m.afterPageChange != nil {
			m.afterPageChange()
		}
		m.afterPageChange = nil
		return
	}
	switch pageURL {
	case SettingsPageURL:
		m.pagesStack = []Page{m.settingsSideBar}
	default:
		m.prePushedView = m.CurrentPage()
		m.pagesStack = append(m.pagesStack, page)
	}
	m.pageAnimation.Disappear(time.Now())
}

func (m *AppManager) NavigateToUrl(pageURL URL, AfterNavCallback func()) {
	m.afterPageChange = AfterNavCallback
	if pageURL == m.CurrentPage().URL() {
		if m.afterPageChange != nil {
			m.afterPageChange()
		}
		m.afterPageChange = nil
		return
	}
	var page Page
	lastSegment := strings.Split(string(pageURL), "/")
	var bookingID uint
	if len(lastSegment) > 0 {
		bookingIDStr := lastSegment[len(lastSegment)-1]
		bookingIDUint64, err := strconv.ParseUint(bookingIDStr, 10, 64)
		if err == nil {
			bookingID = uint(bookingIDUint64)
		}
	}
	switch pageURL {
	case SettingsPageURL:
		m.pagesStack = []Page{m.settingsSideBar}
	case BookingsPageURL:
		page = bookings.New(m)
	case AddEditBookingPageURL(int64(bookingID)):
		booking := sqlc.Booking{}
		booking.ID = int64(bookingID)
		page = add_edit_booking.New(m, booking)
	case ThemePageURL:
		page = theme.New(m)
	case NotificationsPageURL:
		page = notifications.New(m)
	case HelpPageURL:
		page = help.New(m)
	case AboutPageURL:
		page = about.New(m)
	}
	if page != nil {
		m.prePushedView = m.CurrentPage()
		m.pagesStack = append(m.pagesStack, page)
	}
	m.pageAnimation.Disappear(time.Now())
}

func (m *AppManager) PopUp() {
	stackLength := len(m.pagesStack)
	if stackLength > 1 {
		m.poppedUpView = m.pagesStack[stackLength-1]
		m.pagesStack = m.pagesStack[0 : stackLength-1]
	}
	m.pageAnimation.Disappear(time.Now())
}
