package view

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/common/model"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/helper"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
)

type TransactionForm struct {
	Transaction    model.Transaction
	DetailsField   component.TextField
	AmountField    component.TextField
	PreviousAmount string
	Theme          *material.Theme
	Submit         widget.Clickable
	Cancel         widget.Clickable
}

func NewTransactionForm(transaction model.Transaction, theme *material.Theme) TransactionForm {
	tr := TransactionForm{Transaction: transaction, Theme: theme}
	return tr
}

func (tr *TransactionForm) Layout(gtx fwk.Gtx) fwk.Dim {
	if tr.Theme == nil {
		tr.Theme = user.Theme()
	}
	if tr.PreviousAmount != tr.AmountField.Text() {
		tr.AmountField.ClearError()
		tr.PreviousAmount = tr.AmountField.Text()
	}
	// If booking id is new, then Transaction cannot be made as it depends upon existing booking
	if helper.IsNilObjectID(tr.Transaction.BookingID) {
		return fwk.Dim{}
	}
	flex := layout.Flex{Axis: layout.Vertical}
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				labelText := i18n.Get(key.EditTransaction)
				if helper.IsNilObjectID(tr.Transaction.ID) {
					labelText = i18n.Get(key.AddTransaction)
				}
				label := material.H5(tr.Theme, labelText)
				return label.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			labelText := i18n.Get(key.Details)
			return DrawFormFieldRowWithLabel(gtx, tr.Theme, labelText, labelText, &tr.DetailsField, nil)
		}),
		layout.Rigid(layout.Spacer{Height: 16}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			labelText := i18n.Get(key.Amount)
			return DrawFormFieldRowWithLabel(gtx, tr.Theme, labelText, labelText, &tr.AmountField, nil)
		}),
		layout.Rigid(layout.Spacer{Height: 16}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			flex := layout.Flex{Spacing: layout.SpaceSides}
			return flex.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btnText := i18n.Get(key.AddTransaction)
					if !helper.IsNilObjectID(tr.Transaction.ID) {
						btnText = i18n.Get(key.UpdateTransaction)
					}
					btn := material.Button(tr.Theme, &tr.Submit, btnText)
					return btn.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: 16}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btnText := i18n.Get(key.Cancel)
					btn := material.Button(tr.Theme, &tr.Cancel, btnText)
					return btn.Layout(gtx)
				}),
			)
		}),
	)
}
