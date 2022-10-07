package view

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mearaj/bhagad-house-booking/common/assets/fonts"
	. "github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	. "github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
)

type NoBookingView struct {
	Manager
	buttonNewBooking *IconButton
	*material.Theme
	*widget.Icon
	inActiveTh      *material.Theme
	iconCreateNewID *widget.Icon
	*ModalContent
}

func NewNoBooking(manager Manager) *NoBookingView {
	acc := NoBookingView{Manager: manager, Theme: manager.Theme()}
	acc.ModalContent = NewModalContent(func() {
		acc.Modal().Dismiss(nil)
	})
	return &acc
}

func (na *NoBookingView) Layout(gtx Gtx) Dim {
	if na.Theme == nil {
		na.Theme = fonts.NewTheme()
	}
	if na.Icon == nil {
		na.Icon, _ = widget.NewIcon(icons.ActionAccountCircle)
	}
	if na.inActiveTh == nil {
		inActiveTh := *fonts.NewTheme()
		inActiveTh.ContrastBg = color.NRGBA(colornames.Grey500)
		na.inActiveTh = &inActiveTh
	}
	if na.iconCreateNewID == nil {
		na.iconCreateNewID, _ = widget.NewIcon(icons.ContentCreate)
	}
	if na.buttonNewBooking == nil {
		na.buttonNewBooking = &IconButton{
			Theme: na.Theme,
			Icon:  na.Icon,
			Text:  "Add/Edit Booking",
		}
	}

	flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceSides, Alignment: layout.Middle}
	gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	if na.buttonNewBooking.Button.Clicked() {
		na.Manager.NavigateToUrl(AddEditBookingPageURL(Booking{}.ID), func() {})
	}
	d := flex.Layout(gtx,
		layout.Rigid(func(gtx Gtx) Dim {
			return DrawAppIconImageCenter(gtx, na.Theme)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx Gtx) Dim {
			return layout.Center.Layout(gtx, func(gtx Gtx) Dim {
				bdy := material.Body1(na.Theme, "No Booking(s) Created")
				bdy.Alignment = text.Middle
				bdy.Font.Weight = text.Black
				bdy.Color = color.NRGBA{R: 102, G: 117, B: 127, A: 255}
				return bdy.Layout(gtx)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx Gtx) Dim {
			return layout.Flex{Spacing: layout.SpaceSides}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X = gtx.Dp(250)
				return na.buttonNewBooking.Layout(gtx)
			}))
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
	)
	return d
}

func (na *NoBookingView) onSuccess() {
	na.Modal().Dismiss(func() {
		na.Manager.NavigateToUrl(SettingsPageURL, nil)
	})
}
