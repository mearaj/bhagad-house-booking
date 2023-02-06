// Package fwk stands for framework
package fwk

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/x/component"
	"gioui.org/x/notify"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"image/color"
)

type Manager interface {
	NavigateToPage(page Page)
	NavigateToURL(pageURL URL)
	PopUp()
	CurrentPage() Page
	GetWindowWidthInDp() int
	GetWindowWidthInPx() int
	GetWindowHeightInDp() int
	GetWindowHeightInPx() int
	IsStageRunning() bool
	Service() service.Service
	Window() *app.Window
	Notifier() notify.Notifier
	Modal() Modal
	SystemInsets() system.Insets
	ShouldDrawSidebar() bool
	Snackbar() Snackbar
	User() response.LoginUser
}

type Modal interface {
	Show(widget layout.Widget, onBackdropClickCallback func(), animation Animation)
	Dismiss(afterDismiss func())
	View
}
type Snackbar interface {
	Show(txt string, actionButton *widget.Clickable, actionColor color.NRGBA, actionText string)
	View
}

type ViewWidget interface {
	Layout(gtx Gtx, widget layout.Widget) Dim
}

type View interface {
	Layout(gtx Gtx) Dim
}

type Page interface {
	View
	URL() URL
}

// ServiceListener is also responsible for firing event to it's child
type ServiceListener interface {
	OnServiceStateChange(event service.Event)
}

// PagePostPopUp is a page which is active after previous page is popped up
type PagePostPopUp interface {
	Page
	OnPopUpPreviousPage()
}

type URL string

const (
	NavPageURL            URL = "/nav"
	BookingsPageURL       URL = "/bookings"
	TransactionsPageURL   URL = "/transactions"
	SearchBookingsPageURL URL = "/search"
	SettingsPageURL       URL = "/settings"
	NotificationsPageURL  URL = "/notifications"
	HelpPageURL           URL = "/help"
	AboutPageURL          URL = "/about"
)

func AddEditBookingPageURL(bookingID string) URL {
	return URL(fmt.Sprintf("%s/%s", BookingsPageURL, bookingID))
}

func AddEditTransactionsPageURL(bookingID string) URL {
	return URL(fmt.Sprintf("%s/%s", TransactionsPageURL, bookingID))
}

type (
	Gtx       = layout.Context
	Dim       = layout.Dimensions
	Animation = component.VisibilityAnimation
)
