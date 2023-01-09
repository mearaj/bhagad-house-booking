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
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"image/color"
)

type Manager interface {
	NavigateToPage(page Page)
	NavigateToUrl(pageURL URL)
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
	PageFromUrl(url URL) Page
	SystemInsets() system.Insets
	ShouldDrawSidebar() bool
	Snackbar() Snackbar
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

//type ServiceListener interface {
//	OnServiceStateChange(event service.Event)
//}

type URL string

const (
	NavPageUrl            URL = "/nav"
	BookingsPageURL       URL = "/bookings"
	SearchBookingsPageURL URL = "/search"
	SettingsPageURL       URL = "/settings"
	NotificationsPageURL  URL = "/notifications"
	HelpPageURL           URL = "/help"
	AboutPageURL          URL = "/about"
)

func AddEditBookingPageURL(bookingID int64) URL {
	return URL(fmt.Sprintf("%s/%d", BookingsPageURL, bookingID))
}

type (
	Gtx       = layout.Context
	Dim       = layout.Dimensions
	Animation = component.VisibilityAnimation
)
