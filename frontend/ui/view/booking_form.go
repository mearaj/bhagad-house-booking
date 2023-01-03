package view

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/common/assets/fonts"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/giowidgets/calendar"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
	"time"
)

type BookingForm struct {
	Manager
	Theme              *material.Theme
	btnStartDate       widget.Clickable
	btnEndDate         widget.Clickable
	btnClearStartDate  IconButton
	btnClearEndDate    IconButton
	startFieldCalendar calendar.Calendar
	endFieldCalendar   calendar.Calendar
	initialized        bool
	StartDate          time.Time
	EndDate            time.Time
	ButtonSubmit       widget.Clickable
	Details            component.TextField
	ShowDetails        bool
	ButtonText         string
}

// NewBookingForm Always call this function to create BookingForm
func NewBookingForm(manager Manager, booking service.Booking, showDetails bool) BookingForm {
	clearIcon, _ := widget.NewIcon(icons.ContentClear)
	inActiveTheme := fonts.NewTheme()
	inActiveTheme.ContrastBg = color.NRGBA(colornames.Grey500)
	contForm := BookingForm{
		Manager:   manager,
		Theme:     manager.Theme(),
		StartDate: booking.StartDate,
		EndDate:   booking.EndDate,
		btnClearStartDate: IconButton{
			Theme: manager.Theme(),
			Icon:  clearIcon,
			Text:  "",
		},
		btnClearEndDate: IconButton{
			Theme: manager.Theme(),
			Icon:  clearIcon,
			Text:  "",
		},
		ShowDetails: showDetails,
	}
	contForm.Details.SetText(booking.Details)
	return contForm
}

func (bf *BookingForm) Layout(gtx Gtx) Dim {
	if !bf.initialized {
		if bf.Theme == nil {
			bf.Theme = fonts.NewTheme()
		}
		ModalContentInstance.OnCloseClick = func() {
			bf.Modal().Dismiss(nil)
		}
	}

	flex := layout.Flex{Axis: layout.Vertical}
	return flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if bf.btnStartDate.Clicked() {
				bf.Modal().Show(bf.showStartFieldCalendar, nil, Animation{
					Duration: time.Millisecond * 250,
					State:    component.Invisible,
					Started:  time.Time{},
				})
			}
			if bf.btnClearStartDate.Button.Clicked() {
				bf.StartDate = time.Time{}
			}
			return bf.drawDateField(gtx, "Start Date", &bf.btnStartDate, &bf.btnClearStartDate, bf.StartDate)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if bf.btnEndDate.Clicked() {
				bf.Modal().Show(bf.showEndFieldCalendar, nil, Animation{
					Duration: time.Millisecond * 250,
					State:    component.Invisible,
					Started:  time.Time{},
				})
			}
			if bf.btnClearEndDate.Button.Clicked() {
				bf.EndDate = time.Time{}
			}
			return bf.drawDateField(gtx, "End Date", &bf.btnEndDate, &bf.btnClearEndDate, bf.EndDate)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !bf.ShowDetails {
				return Dim{}
			}
			inset := layout.UniformInset(16)
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return DrawFormFieldRowWithLabel(gtx, bf.Theme, "Details", "Details", &bf.Details, nil)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			inset := layout.UniformInset(16)
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				if bf.ButtonText == "" {
					bf.ButtonText = "Apply"
				}
				btn := material.Button(bf.Theme, &bf.ButtonSubmit, bf.ButtonText)
				return btn.Layout(gtx)
			})
		}),
	)
}

func (bf *BookingForm) drawDateField(gtx Gtx, label string, btnDate *widget.Clickable, btnClearDate *IconButton, t time.Time) Dim {
	// fieldVal := "Enter dd/mm/yyyy"
	fieldVal := "Tap to enter date"
	labelWidth := gtx.Dp(100)
	flx := layout.Flex{Axis: layout.Vertical}
	inset := layout.UniformInset(unit.Dp(16))
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return flx.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X, gtx.Constraints.Max.X = labelWidth, labelWidth
				inset := layout.Inset{Bottom: unit.Dp(4)}
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return material.Label(bf.Theme, bf.Theme.TextSize, label).Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				flx := layout.Flex{Alignment: layout.Middle}
				return flx.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return btnDate.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return widget.Border{
								Color:        bf.Theme.ContrastBg,
								CornerRadius: 0,
								Width:        unit.Dp(1),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								inset := layout.UniformInset(unit.Dp(12))
								return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									lbl := material.Label(bf.Theme, bf.Theme.TextSize, fieldVal)
									lbl.Color = color.NRGBA(colornames.Grey500)
									lbl.Font.Weight = text.Normal
									if !t.IsZero() {
										lbl.Text = fmt.Sprintf("%d %s %d", t.Day(), t.Month().String(), t.Year())
										lbl.Color = bf.Theme.Fg
										lbl.Font.Weight = text.Bold
									}
									return lbl.Layout(gtx)
								})
							})
						})
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btnClearDate.Inset = layout.UniformInset(unit.Dp(8))
						return btnClearDate.Layout(gtx)
					}),
				)
			}),
		)
	})
}

func (bf *BookingForm) showStartFieldCalendar(gtx Gtx) Dim {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.85)
	bf.startFieldCalendar.OnCalendarDateClick = bf.onCalendarStartDateFieldClick
	bf.startFieldCalendar.Inset = layout.UniformInset(unit.Dp(16))
	return ModalContentInstance.DrawContent(gtx, bf.Theme, bf.startFieldCalendar.Layout)
}

func (bf *BookingForm) showEndFieldCalendar(gtx Gtx) Dim {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.85)
	bf.endFieldCalendar.OnCalendarDateClick = bf.onCalendarEndDateFieldClick
	bf.endFieldCalendar.Inset = layout.UniformInset(unit.Dp(16))
	return ModalContentInstance.DrawContent(gtx, bf.Theme, bf.endFieldCalendar.Layout)
}
func (bf *BookingForm) onCalendarStartDateFieldClick(t time.Time) {
	bf.Modal().Dismiss(nil)
	bf.StartDate = t
}

func (bf *BookingForm) onCalendarEndDateFieldClick(t time.Time) {
	bf.Modal().Dismiss(nil)
	bf.EndDate = t
}
