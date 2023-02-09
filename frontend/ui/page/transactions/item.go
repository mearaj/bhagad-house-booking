package transactions

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/common/model"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	color2 "image/color"
	"math"
	"net/mail"
	"strings"
	"time"
)

type transactionItem struct {
	fwk.Animation
	transaction    model.Transaction
	BtnHeader      widget.Clickable
	BtnEdit        widget.Clickable
	BtnDelete      widget.Clickable
	BtnYes         widget.Clickable
	BtnNo          widget.Clickable
	BtnSendSMS     widget.Clickable
	BtnSendEmail   widget.Clickable
	IsSendingEmail bool
	IsSendingSMS   bool
	parent         *page
	*material.Theme
}

var editIcon, _ = widget.NewIcon(icons.EditorModeEdit)
var deleteIcon, _ = widget.NewIcon(icons.ActionDelete)
var iconDone, _ = widget.NewIcon(icons.ActionDone)

func (tr *transactionItem) Layout(gtx fwk.Gtx, index int) view.Dim {
	if tr.Animation == (fwk.Animation{}) {
		tr.Animation.Duration = time.Millisecond * 100
		tr.Animation.State = component.Invisible
	}
	if tr.BtnHeader.Clicked() {
		tr.Animation.ToggleVisibility(gtx.Now)
	}
	tr.handleSendSMSSubmit()
	tr.handleSendEmailSubmit()
	inset := layout.Inset{Top: 8, Bottom: 8}
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		flex := layout.Flex{Axis: layout.Vertical}
		return flex.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btnStyle := material.ButtonLayoutStyle{Button: &tr.BtnHeader}
				btnStyle.Background = tr.Theme.ContrastBg
				d := btnStyle.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return tr.layoutHeader(gtx, index)
				})
				return d
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				progress := tr.Animation.Revealed(gtx)
				macro := op.Record(gtx.Ops)
				d := layout.Flex{}.Layout(gtx, layout.Flexed(1.0, func(gtx view.Gtx) view.Dim {
					return layout.Inset{
						Top:    unit.Dp(0),
						Bottom: unit.Dp(6),
					}.Layout(gtx, func(gtx view.Gtx) view.Dim {
						return tr.layoutChild(gtx)
					})
				}))
				call := macro.Stop()
				height := int(math.Round(float64(float32(d.Size.Y) * progress)))
				d.Size.Y = height
				defer clip.Rect(image.Rectangle{
					Max: d.Size,
				}).Push(gtx.Ops).Pop()
				call.Add(gtx.Ops)
				return d
			}),
		)
	})
}
func (tr *transactionItem) layoutHeader(gtx fwk.Gtx, index int) fwk.Dim {
	th := tr.Theme
	inset := layout.Inset{Top: 6, Right: 12, Bottom: 6, Left: 12}
	d := inset.Layout(gtx, func(gtx fwk.Gtx) fwk.Dim {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
				gtx.Constraints.Min.X = gtx.Dp(30)
				txt := fmt.Sprintf("%d.", index+1)
				label := material.Label(tr.Theme, 16.0, txt)
				label.Color = tr.Theme.ContrastFg
				return label.Layout(gtx)
			}),
			layout.Flexed(1, func(gtx fwk.Gtx) fwk.Dim {
				flex := layout.Flex{Spacing: layout.SpaceBetween}
				return flex.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						txt := utils.GetFormattedDate(tr.transaction.CreatedAt)
						label := material.Label(th, unit.Sp(14), txt)
						label.Color = th.ContrastFg
						label.Font.Weight = text.Bold
						return label.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						txt := fmt.Sprintf("%.2f", tr.transaction.Amount)
						label := material.Label(th, unit.Sp(14), txt)
						label.Color = th.ContrastFg
						label.Font.Weight = text.Bold
						return label.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(func(gtx fwk.Gtx) (d fwk.Dim) {
				affine := f32.Affine2D{}
				ic, _ := widget.NewIcon(icons.NavigationChevronRight)
				cl := th.ContrastFg
				origin := f32.Pt(12, 12)
				rotation := float32(0)
				if tr.Animation.Visible() {
					rotation = float32(math.Pi * 0.5)
				}
				if tr.Animation.Animating() {
					rotation *= tr.Animation.Revealed(gtx)
				}
				affine = affine.Rotate(origin, rotation)
				defer op.Affine(affine).Push(gtx.Ops).Pop()
				return ic.Layout(gtx, cl)
			}),
		)
	})
	return d
}

func (tr *transactionItem) layoutChild(gtx fwk.Gtx) fwk.Dim {
	inset := layout.Inset{Top: 6, Right: 12, Bottom: 6, Left: 12}
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		flex := layout.Flex{Axis: layout.Vertical}
		return flex.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				flex := layout.Flex{Alignment: layout.Middle}
				return flex.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						txt := tr.transaction.Details
						label := material.Label(tr.Theme, 16.0, txt)
						return label.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(layout.Spacer{Height: 16}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return tr.drawButtons(gtx)
			}),
		)
	})
}

func (tr *transactionItem) drawButtons(gtx fwk.Gtx) fwk.Dim {
	flex := layout.Flex{Alignment: layout.Middle}
	isWidthSmall := gtx.Constraints.Max.X < gtx.Dp(400)
	_, err := mail.ParseAddress(tr.parent.Booking.CustomerEmail)
	isEmailValid := err == nil
	isPhoneValid := utils.ValidateIndianPhoneNumber(tr.parent.Booking.CustomerPhone)

	return flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			background := color2.NRGBA(colornames.Green500)
			if isWidthSmall {
				iconBtn := material.IconButton(tr.Theme, &tr.BtnEdit, editIcon, "Edit")
				iconBtn.Background = background
				return iconBtn.Layout(gtx)
			}
			btn := material.Button(tr.Theme, &tr.BtnEdit, i18n.Get(key.Edit))
			btn.Background = background
			return btn.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Width: 16}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			background := color2.NRGBA(colornames.Red500)
			if isWidthSmall {
				iconBtn := material.IconButton(tr.Theme, &tr.BtnDelete, deleteIcon, "Delete")
				iconBtn.Background = background
				return iconBtn.Layout(gtx)
			}
			btn := material.Button(tr.Theme, &tr.BtnDelete, i18n.Get(key.Delete))
			btn.Background = background
			return btn.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Width: 16}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !isPhoneValid {
				return layout.Dimensions{}
			}
			txt := i18n.Get(key.SendSMS)
			if isWidthSmall {
				txt = "SMS"
			}
			btn := material.Button(tr.Theme, &tr.BtnSendSMS, txt)
			return btn.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !isPhoneValid {
				return layout.Dimensions{}
			}
			inset := layout.Inset{Left: 8}
			if tr.IsSendingSMS {
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					loader := view.Loader{
						Theme: tr.Theme,
						Size:  image.Point{X: gtx.Dp(32), Y: gtx.Dp(32)},
					}
					return loader.Layout(gtx)
				})
			}
			if tr.transaction.ConfirmSMSSent {
				constraints := image.Point{X: gtx.Dp(48), Y: gtx.Dp(48)}
				gtx.Constraints.Min, gtx.Constraints.Max = constraints, constraints
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return iconDone.Layout(gtx, color2.NRGBA(colornames.Green500))
				})
			}
			return layout.Dimensions{}
		}),
		layout.Rigid(layout.Spacer{Width: 16}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !isEmailValid {
				return layout.Dimensions{}
			}
			txt := i18n.Get(key.SendEmail)
			if isWidthSmall {
				txt = i18n.Get(key.Email)
			}
			btn := material.Button(tr.Theme, &tr.BtnSendEmail, txt)
			return btn.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !isEmailValid {
				return layout.Dimensions{}
			}
			inset := layout.Inset{Left: 8}
			if tr.IsSendingEmail {
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					loader := view.Loader{
						Theme: tr.Theme,
						Size:  image.Point{X: gtx.Dp(32), Y: gtx.Dp(32)},
					}
					return loader.Layout(gtx)
				})
			}
			if tr.transaction.ConfirmEmailSent {
				constraints := image.Point{X: gtx.Dp(48), Y: gtx.Dp(48)}
				gtx.Constraints.Min, gtx.Constraints.Max = constraints, constraints
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return iconDone.Layout(gtx, color2.NRGBA(colornames.Green500))
				})
			}
			return layout.Dimensions{}
		}),
	)
}
func (tr *transactionItem) handleSendSMSSubmit() {
	if tr.BtnSendSMS.Clicked() {
		isPhoneValid := utils.ValidateIndianPhoneNumber(tr.parent.Booking.CustomerPhone)
		if isPhoneValid {
			tr.IsSendingSMS = true
			tr.parent.Service().SendNewTransactionSMS(tr.transaction.Number, tr)
		}
	}
}
func (tr *transactionItem) handleSendEmailSubmit() {
	if tr.BtnSendEmail.Clicked() {
		_, err := mail.ParseAddress(tr.parent.Booking.CustomerEmail)
		isEmailValid := err == nil
		if isEmailValid {
			tr.IsSendingEmail = true
			tr.parent.Service().SendNewTransactionEmail(tr.transaction.Number, tr)
		}
	}
}
func (tr *transactionItem) OnServiceStateChange(event service.Event) {
	var errTxt string
	switch eventData := event.Data.(type) {
	case service.NewTransactionSMSResponse:
		if event.Cached || event.ID != tr {
			return
		}
		tr.IsSendingSMS = false
		errTxt = eventData.Error
		tr.parent.Window().Invalidate()
	case service.NewTransactionEmailResponse:
		if event.Cached || event.ID != tr {
			return
		}
		tr.IsSendingEmail = false
		errTxt = eventData.Error
		tr.parent.Window().Invalidate()
	}

	if errTxt != "" {
		if strings.Contains(errTxt, "connection refused") {
			errTxt = "connection refused"
		}
		tr.parent.Snackbar().Show(errTxt, &tr.parent.closeSnapBar, color2.NRGBA{R: 255, A: 255}, i18n.Get(key.Close))
	}
	tr.parent.Window().Invalidate()
}
