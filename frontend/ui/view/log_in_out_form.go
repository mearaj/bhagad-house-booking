package view

import (
	key2 "gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/frontend/assets/fonts"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"image"
	"image/color"
	"strings"
)

type UserForm struct {
	Manager
	Theme             *material.Theme
	email             FormField
	password          FormField
	btnLogInOut       widget.Clickable
	isLoggingInOut    bool
	loginUserResponse service.UserResponse
	subscription      service.Subscriber
}

func NewUserForm(manager Manager) *UserForm {
	inActiveTheme := fonts.NewTheme()
	inActiveTheme.ContrastBg = color.NRGBA(colornames.Grey500)
	contForm := UserForm{
		Manager: manager,
		Theme:   user.Theme(),
		email: FormField{FieldName: i18n.Get(key.Email),
			TextField: component.TextField{Editor: widget.Editor{SingleLine: true, Submit: true}},
		},
		password: FormField{FieldName: i18n.Get(key.Password),
			TextField: component.TextField{Editor: widget.Editor{SingleLine: true, Submit: true}},
		},
		subscription: manager.Service().Subscribe(service.TopicLoggedInOut),
	}
	contForm.email.InputHint = key2.HintEmail
	contForm.subscription.SubscribeWithCallback(contForm.OnServiceStateChange)
	return &contForm
}

func (p *UserForm) Layout(gtx Gtx) Dim {
	if p.Theme == nil {
		p.Theme = fonts.NewTheme()
	}

	inset := layout.UniformInset(unit.Dp(16))
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		flex := layout.Flex{Axis: layout.Vertical}
		return flex.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					title := i18n.Get(key.AdminLogin)
					if p.loginUserResponse.IsAuthorized() {
						title = p.loginUserResponse.User.Name
						if title == "" {
							title = p.loginUserResponse.User.Email
						}
					}
					return material.H4(p.Theme, title).Layout(gtx)
				})
			}),
			layout.Rigid(layout.Spacer{Height: 8}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return p.inputField(gtx, &p.email)
			}),
			layout.Rigid(layout.Spacer{Height: 16}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return p.inputField(gtx, &p.password)
			}),
			layout.Rigid(layout.Spacer{Height: 16}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return p.logInOut(gtx)
			}),
		)
	})
}

func (p *UserForm) inputField(gtx Gtx, field *FormField) Dim {
	if p.loginUserResponse.IsAuthorized() {
		return Dim{}
	}
	return DrawFormFieldRowWithLabel(gtx, p.Theme, "", field.FieldName, &field.TextField, nil)
}

func (p *UserForm) submitted() (submitted bool) {
	for _, e := range p.email.Events() {
		if _, ok := e.(widget.SubmitEvent); ok {
			submitted = true
			return
		}
	}
	for _, e := range p.password.Events() {
		if _, ok := e.(widget.SubmitEvent); ok {
			submitted = true
			return
		}
	}
	return
}

func (p *UserForm) logInOut(gtx Gtx) Dim {
	title := i18n.Get(key.LogIn)
	if p.loginUserResponse.IsAuthorized() {
		title = i18n.Get(key.LogOut)
	}

	submitted := p.submitted()
	if !p.isLoggingInOut {
		if p.btnLogInOut.Clicked() || submitted {
			p.isLoggingInOut = true
			if !p.loginUserResponse.IsAuthorized() {
				email := strings.TrimSpace(p.email.Text())
				password := strings.TrimSpace(p.password.Text())
				if email != "" && password != "" {
					p.Service().LogInUser(p.email.Text(), p.password.Text())
				} else {
					p.isLoggingInOut = false
				}
			}
			if p.loginUserResponse.IsAuthorized() {
				p.Service().LogOutUser()
			}
		}
		btn := material.Button(p.Theme, &p.btnLogInOut, title)
		return btn.Layout(gtx)
	}
	loader := Loader{
		Size: image.Point{X: 24, Y: 24},
	}
	gtx.Constraints.Max.X, gtx.Constraints.Max.Y = 24, 24
	return loader.Layout(gtx)
}
func (p *UserForm) OnServiceStateChange(event service.Event) {
	if userResponse, ok := event.Data.(service.UserResponse); ok {
		p.isLoggingInOut = false
		p.loginUserResponse = userResponse
	}

}
