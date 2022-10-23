package view

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/common/assets/fonts"
	. "github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"image/color"
)

const customerFieldName = "Name"
const customerFieldEmail = "Email"
const customerFieldContact = "Contact"
const customerFieldAddress = "Address"

type customerField struct {
	FieldName string
	component.TextField
}

var customerFields = []*customerField{
	{FieldName: customerFieldName},
	{FieldName: customerFieldEmail},
	{FieldName: customerFieldContact},
	{FieldName: customerFieldAddress},
}

// CustomerForm Always call NewCustomerForm function to create CustomerForm
type CustomerForm struct {
	Manager
	Theme              *material.Theme
	Customer           Customer
	customerFieldsList layout.List
	OnSuccess          func(addr string)
}

// NewCustomerForm Always call this function to create CustomerForm
func NewCustomerForm(manager Manager, customer Customer, OnSuccess func(addr string)) CustomerForm {
	inActiveTheme := fonts.NewTheme()
	inActiveTheme.ContrastBg = color.NRGBA(colornames.Grey500)
	contForm := CustomerForm{
		Manager:   manager,
		Theme:     manager.Theme(),
		Customer:  customer,
		OnSuccess: OnSuccess,
	}
	return contForm
}

func (p *CustomerForm) Layout(gtx Gtx) Dim {
	if p.Theme == nil {
		p.Theme = fonts.NewTheme()
	}
	p.customerFieldsList.Axis = layout.Vertical

	return p.customerFieldsList.Layout(gtx, len(customerFields), func(gtx layout.Context, index int) layout.Dimensions {
		inset := layout.UniformInset(unit.Dp(16))
		return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return p.inputField(gtx, customerFields[index])
		})
	})
}

func (p *CustomerForm) inputField(gtx Gtx, field *customerField) Dim {
	return DrawFormFieldRowWithLabel(gtx, p.Theme, field.FieldName, field.FieldName, &field.TextField, nil)
}
