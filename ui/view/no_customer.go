package view

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mearaj/bhagad-house-booking/service"
	. "github.com/mearaj/bhagad-house-booking/ui/fwk"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
)

type NoCustomerView struct {
	Manager
	buttonAddCustomer *IconButton
	*material.Theme
	*widget.Icon
	CustomerFormView View
	*ModalContent
}

func NewNoCustomer(manager Manager, onSuccess func(contactAddr string), btnText string) *NoCustomerView {
	btnIcon, _ := widget.NewIcon(icons.CommunicationContacts)
	if btnText == "" {
		btnText = "Add Customer"
	}
	nc := NoCustomerView{
		Manager: manager,
		Theme:   manager.Theme(),
		buttonAddCustomer: &IconButton{
			Theme: manager.Theme(),
			Icon:  btnIcon,
			Text:  btnText,
		},
	}
	return &nc
}

func (nc *NoCustomerView) Layout(gtx Gtx) Dim {
	flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceSides, Alignment: layout.Middle}
	gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	if nc.buttonAddCustomer.Button.Clicked() {
		nc.Manager.NavigateToUrl(AddEditCustomerPageURL(service.Customer{}.ID), func() {})
	}
	d := flex.Layout(gtx,
		layout.Rigid(func(gtx Gtx) Dim {
			return DrawAppIconImageCenter(gtx, nc.Theme)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx Gtx) Dim {
			return layout.Center.Layout(gtx, func(gtx Gtx) Dim {
				bdy := material.Body1(nc.Theme, "No Customer(s) Found")
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
				return nc.buttonAddCustomer.Layout(gtx)
			}))
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
	)
	return d
}
