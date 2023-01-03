package view

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/common/assets/fonts"
	"github.com/mearaj/bhagad-house-booking/frontend/service"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"image"
	"image/color"
)

const UserFieldEmail = "Email"
const UserFieldPassword = "Password"
const loginTitle = "Log In"
const logoutTitle = "Log Out"

type FormField struct {
	FieldName string
	component.TextField
}

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
		Manager:      manager,
		Theme:        manager.Theme(),
		email:        FormField{FieldName: UserFieldEmail},
		password:     FormField{FieldName: UserFieldPassword},
		subscription: manager.Service().Subscribe(service.TopicUserLoggedInOut),
	}
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
					title := "Admin Login"
					isAuthorized := p.loginUserResponse.IsLoggedIn() && p.loginUserResponse.IsAdmin()
					if isAuthorized {
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
	isAuthorized := p.loginUserResponse.IsLoggedIn() && p.loginUserResponse.IsAdmin()
	if isAuthorized {
		return Dim{}
	}
	return DrawFormFieldRowWithLabel(gtx, p.Theme, "", field.FieldName, &field.TextField, nil)
}

func (p *UserForm) logInOut(gtx Gtx) Dim {
	title := loginTitle
	isAuthorized := p.loginUserResponse.IsLoggedIn() && p.loginUserResponse.IsAdmin()
	if isAuthorized {
		title = logoutTitle
	}
	if !p.isLoggingInOut {
		if p.btnLogInOut.Clicked() {
			p.isLoggingInOut = true
			if !isAuthorized {
				p.Service().LogInUser(p.email.Text(), p.password.Text())
			}
			if isAuthorized {
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
	switch userResponse := event.Data.(type) {
	case service.UserResponse:
		p.isLoggingInOut = false
		p.loginUserResponse = userResponse
	}

}
