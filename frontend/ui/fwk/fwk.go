// Package fwk stands for framework
package fwk

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"gioui.org/x/notify"
	service2 "github.com/mearaj/bhagad-house-booking/frontend/service"
	"image/color"
)

type Manager interface {
	NavigateToPage(page Page, AfterNavCallback func())
	NavigateToUrl(pageURL URL, AfterNavCallback func())
	PopUp()
	CurrentPage() Page
	GetWindowWidthInDp() int
	GetWindowWidthInPx() int
	GetWindowHeightInDp() int
	GetWindowHeightInPx() int
	IsStageRunning() bool
	Theme() *material.Theme
	Service() service2.Service
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

type DatabaseListener interface {
	OnDatabaseChange(event service2.Event)
}

type ViewWithDBListener interface {
	View
	DatabaseListener
}

type URL string

const (
	SettingsPageURL      URL = "/settings"
	BookingsPageURL          = SettingsPageURL + "/bookings"
	CustomersPageURL         = SettingsPageURL + "/customers"
	ThemePageURL             = SettingsPageURL + "/theme"
	NotificationsPageURL     = SettingsPageURL + "/notifications"
	HelpPageURL              = SettingsPageURL + "/help"
	AboutPageURL             = SettingsPageURL + "/about"
)

func AddEditBookingPageURL(bookingID uint) URL {
	return URL(fmt.Sprintf("%s/%d", BookingsPageURL, bookingID))
}

func AddEditCustomerPageURL(customerID uint) URL {
	return URL(fmt.Sprintf("%s/%d", CustomersPageURL, customerID))
}

type (
	Gtx       = layout.Context
	Dim       = layout.Dimensions
	Animation = component.VisibilityAnimation
)
