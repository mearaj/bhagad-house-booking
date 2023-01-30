package transactions

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/common/model"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/helper"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"strconv"
	"strings"
	"time"
)

type page struct {
	initialized            bool
	isFetchingTransactions bool
	title                  string
	ViewList               layout.List
	Theme                  *material.Theme
	btnNav                 widget.Clickable
	btnAddTransaction      widget.Clickable
	closeSnapBar           widget.Clickable
	navigationIcon         *widget.Icon
	iconAddTransaction     *widget.Icon
	bookingDetails         view.BookingDetails
	loginUserResponse      service.UserResponse
	subscription           service.Subscriber
	fwk.Manager
	Booking             service.Booking
	selectedTransaction view.TransactionForm
	transactions        []*transactionItem
	ModalContent        *view.ModalContent
}

func New(manager fwk.Manager, booking service.Booking) fwk.Page {
	navIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	addTransactionIcon, _ := widget.NewIcon(icons.ContentAdd)
	theme := *user.Theme()
	p := page{
		Manager:            manager,
		Theme:              &theme,
		navigationIcon:     navIcon,
		ViewList:           layout.List{Axis: layout.Vertical},
		Booking:            booking,
		iconAddTransaction: addTransactionIcon,
		bookingDetails: view.BookingDetails{
			Booking: booking,
			Theme:   user.Theme(),
		},
	}
	p.subscription = manager.Service().Subscribe()
	p.ModalContent = view.NewModalContent(p.onModalCloseClick)
	p.subscription.SubscribeWithCallback(p.OnServiceStateChange)
	return &p
}

func (p *page) Layout(gtx fwk.Gtx) fwk.Dim {
	if !p.initialized {
		if p.Theme == nil {
			p.Theme = user.Theme()
		}
		p.fetchTransactions()
		p.initialized = true
	}
	p.title = i18n.Get(key.ManageTransactions)
	if !p.loginUserResponse.IsAuthorized() {
		return fwk.Dim{}
	}
	if p.btnAddTransaction.Clicked() {
		p.selectedTransaction = view.NewTransactionForm(model.Transaction{BookingID: p.Booking.ID}, user.Theme())
		p.Modal().Show(func(gtx layout.Context) layout.Dimensions {
			return p.drawTransactionForm(gtx)
		}, nil, view.Animation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		})
	}

	if p.selectedTransaction.Submit.Clicked() {
		p.handleSubmitTransaction()
	}
	if p.selectedTransaction.Cancel.Clicked() {
		p.Modal().Dismiss(nil)
		p.selectedTransaction = view.NewTransactionForm(model.Transaction{BookingID: p.Booking.ID}, user.Theme())
	}

	flex := layout.Flex{Axis: layout.Vertical}
	d := flex.Layout(gtx,
		layout.Rigid(p.DrawAppBar),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			inset := layout.UniformInset(view.BookingDetailsFormInset)
			return inset.Layout(gtx, func(gtx fwk.Gtx) fwk.Dim {
				flex := layout.Flex{Axis: layout.Vertical}
				return flex.Layout(gtx,
					layout.Rigid(p.drawBookingDetails),
					layout.Rigid(layout.Spacer{Height: 16}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min = gtx.Constraints.Max
						return p.ViewList.Layout(gtx, len(p.transactions), func(gtx layout.Context, index int) layout.Dimensions {
							return p.drawTransactionItem(gtx, index)
						})
					}),
				)
			})
		}),
	)
	p.drawFloatingButton(gtx)
	return d
}

func (p *page) DrawAppBar(gtx fwk.Gtx) fwk.Dim {
	gtx.Constraints.Max.Y = gtx.Dp(56)
	th := p.Theme
	if p.btnNav.Clicked() {
		p.PopUp()
	}

	return view.DrawAppBarLayout(gtx, th, func(gtx fwk.Gtx) fwk.Dim {
		return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
						navigationIcon := p.navigationIcon
						button := material.IconButton(th, &p.btnNav, navigationIcon, "Nav Icon Button")
						button.Size = unit.Dp(40)
						button.Background = th.Palette.ContrastBg
						button.Color = th.Palette.ContrastFg
						button.Inset = layout.UniformInset(unit.Dp(8))
						return button.Layout(gtx)
					}),
					layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
						return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx fwk.Gtx) fwk.Dim {
							titleText := p.title
							title := material.Body1(th, titleText)
							title.Color = th.Palette.ContrastFg
							title.TextSize = unit.Sp(18)
							return title.Layout(gtx)
						})
					}),
				)
			}),
		)
	})
}

func (p *page) drawTransactionItem(gtx fwk.Gtx, index int) fwk.Dim {
	if len(p.transactions) == 0 {
		return view.Dim{}
	}
	if p.transactions[index].BtnEdit.Clicked() {
		transaction := p.transactions[index].transaction
		p.selectedTransaction.Transaction = transaction
		p.selectedTransaction.AmountField.SetText(fmt.Sprintf("%.2f", transaction.Amount))
		p.selectedTransaction.DetailsField.SetText(transaction.Details)
		p.Modal().Show(func(gtx layout.Context) layout.Dimensions {
			return p.drawTransactionForm(gtx)
		}, nil, view.Animation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		})
	}
	if p.transactions[index].BtnDelete.Clicked() {
		transaction := p.transactions[index].transaction
		p.selectedTransaction.Transaction = transaction
		p.selectedTransaction.AmountField.SetText(fmt.Sprintf("%.2f", transaction.Amount))
		p.selectedTransaction.DetailsField.SetText(transaction.Details)
		p.Modal().Show(func(gtx layout.Context) layout.Dimensions {
			return p.drawDeleteTransactionModel(gtx, index)
		}, nil, view.Animation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		})
	}
	return p.transactions[index].Layout(gtx, index)
}

func (p *page) drawBookingDetails(gtx fwk.Gtx) fwk.Dim {
	flex := layout.Flex{Axis: layout.Vertical}
	return flex.Layout(gtx,
		layout.Rigid(p.bookingDetails.Layout),
		layout.Rigid(p.drawPaymentFields),
	)
}
func (p *page) drawPaymentFields(gtx fwk.Gtx) fwk.Dim {
	flex := layout.Flex{Axis: layout.Vertical}
	d := flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			labelTxt := fmt.Sprintf("%s %s", i18n.Get(key.Total), i18n.Get(key.Received))
			var value float64
			var valueTxt string
			for _, tr := range p.transactions {
				value += tr.transaction.Amount
			}
			valueTxt = fmt.Sprintf("%.2f", value)
			return p.drawBookingField(gtx, labelTxt, valueTxt)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			labelTxt := i18n.Get(key.Balance)
			var totalReceived float64
			var valueTxt string
			for _, tr := range p.transactions {
				totalReceived += tr.transaction.Amount
			}
			totalPrice := helper.BookingTotalPrice(p.Booking.RatePerDay, p.Booking.StartDate, p.Booking.EndDate)
			balanceLeft := totalPrice - totalReceived
			valueTxt = fmt.Sprintf("%.2f", balanceLeft)
			return p.drawBookingField(gtx, labelTxt, valueTxt)
		}),
	)
	return d
}
func (p *page) drawBookingField(gtx fwk.Gtx, labelText string, valueText string) fwk.Dim {
	flex := layout.Flex{Alignment: layout.Start}
	labelField := p.drawLabelField(gtx, labelText)
	return flex.Layout(gtx,
		labelField[0],
		labelField[1],
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			ins := layout.Inset{}
			return ins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				b := material.Body1(p.Theme, valueText)
				b.Font.Weight = view.BookingDetailsBodyFontWeight
				b.TextSize = view.BookingDetailsBodyFontSize
				return b.Layout(gtx)
			})
		}),
	)
}
func (p *page) drawLabelField(gtx fwk.Gtx, labelText string) [2]layout.FlexChild {
	return [2]layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = gtx.Dp(view.BookingDetailsLabelWidth)
			gtx.Constraints.Min.X = gtx.Dp(view.BookingDetailsLabelWidth)
			label := material.Label(p.Theme, view.BookingDetailsHeadFontSize, labelText)
			label.Font.Weight = view.BookingDetailsHeadFontWeight
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Label(p.Theme, view.BookingDetailsHeadFontSize, ": ").Layout(gtx)
		}),
	}
}

func (p *page) drawFloatingButton(gtx fwk.Gtx) fwk.Dim {
	st := layout.Stack{Alignment: layout.NE}
	op.Offset(image.Pt(0, gtx.Dp(view.AppBarHeight))).Add(gtx.Ops)
	return st.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			inset := layout.UniformInset(16)
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.IconButton(p.Theme, &p.btnAddTransaction, p.iconAddTransaction, "Add Transaction")
				return btn.Layout(gtx)
			})
		}),
	)
}

func (p *page) fetchTransactions() {
	if !p.isFetchingTransactions && !helper.IsNilObjectID(p.Booking.ID) {
		p.isFetchingTransactions = true
		p.Service().Transactions(service.TransactionsRequest{BookingID: p.Booking.ID}, p)
	}
}
func (p *page) drawTransactionForm(gtx fwk.Gtx) fwk.Dim {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(p.GetWindowHeightInPx()) * 0.85)
	inset := layout.UniformInset(view.BookingDetailsFormInset)
	return p.ModalContent.DrawContent(gtx, p.Theme, func(gtx view.Gtx) view.Dim {
		gtx.Constraints.Max.Y = int(float32(p.GetWindowHeightInPx()) * 0.85)
		return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return p.selectedTransaction.Layout(gtx)
		})
	})
}

func (p *page) handleSubmitTransaction() {
	amountStr := strings.TrimSpace(p.selectedTransaction.AmountField.Text())
	if amountStr == "" {
		amountStr = "0"
	}
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		p.selectedTransaction.AmountField.SetError("invalid amount")
	}
	if err == nil {
		p.Modal().Dismiss(nil)
		p.Service().AddUpdateTransaction(service.AddUpdateTransactionRequest{
			ID:        p.selectedTransaction.Transaction.ID,
			BookingID: p.selectedTransaction.Transaction.BookingID,
			Amount:    amount,
			Details:   p.selectedTransaction.DetailsField.Text(),
		}, p)
		p.selectedTransaction = view.NewTransactionForm(model.Transaction{BookingID: p.Booking.ID}, user.Theme())
	}
}
func (p *page) drawDeleteTransactionModel(gtx fwk.Gtx, index int) fwk.Dim {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.85)
	if p.transactions[index].BtnYes.Clicked() {
		p.Modal().Dismiss(func() {
			p.Service().DeleteTransaction(service.DeleteTransactionRequest{
				ID:        p.transactions[index].transaction.ID,
				BookingID: p.transactions[index].transaction.BookingID,
			}, p)
		})
	}
	if p.transactions[index].BtnNo.Clicked() {
		p.Modal().Dismiss(func() {})
	}

	delPrompt := i18n.Get(key.TransactionDeletePrompt)
	transactionID := fmt.Sprintf("%s %s %s", i18n.Get(key.Transaction), i18n.Get(key.ID),
		p.transactions[index].transaction.ID.Hex(),
	)
	transactionDetails := p.transactions[index].transaction.Details
	transactionAmount := p.transactions[index].transaction.Amount
	promptContent := view.NewPromptContent(p.Theme,
		i18n.Get(key.TransactionDeletion),
		fmt.Sprintf(
			"%s\n%s\n%s :- %s\n%s :- %.2f\n",
			delPrompt,
			transactionID,
			i18n.Get(key.Details),
			transactionDetails,
			i18n.Get(key.Amount),
			transactionAmount,
		),
		&p.transactions[index].BtnYes,
		&p.transactions[index].BtnNo)
	return p.ModalContent.DrawContent(gtx, p.Theme, promptContent.Layout)
}

func (p *page) OnServiceStateChange(event service.Event) {
	var errTxt string
	switch eventData := event.Data.(type) {
	case service.UserResponse:
		p.loginUserResponse = eventData
	case service.TransactionsResponse:
		p.isFetchingTransactions = false
		if event.Cached || event.ID != p {
			return
		}
		errTxt = eventData.Error
		var transactions []*transactionItem
		for _, tr := range eventData.Transactions {
			transactions = append(transactions, &transactionItem{
				transaction: tr,
				Theme:       p.Theme,
			})
		}
		p.transactions = transactions
	case service.AddUpdateTransactionResponse:
		p.isFetchingTransactions = false
		if event.Cached || event.ID != p {
			return
		}
		p.fetchTransactions()
	case service.DeleteTransactionResponse:
		p.isFetchingTransactions = false
		if event.Cached || event.ID != p {
			return
		}
		errTxt = eventData.Error
		p.fetchTransactions()
	}
	if errTxt != "" {
		if strings.Contains(errTxt, "connection refused") {
			errTxt = "connection refused"
		}
		p.Snackbar().Show(errTxt, &p.closeSnapBar, color.NRGBA{R: 255, A: 255}, i18n.Get(key.Close))
	}
}

func (p *page) onModalCloseClick() {
	p.Modal().Dismiss(nil)
}

func (p *page) URL() fwk.URL {
	return fwk.AddEditTransactionsPageURL(p.Booking.ID.Hex())
}
