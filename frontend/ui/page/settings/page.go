package settings

import (
	"fmt"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/colorpicker"
	"gioui.org/x/component"
	"github.com/mearaj/bhagad-house-booking/frontend/assets/fonts"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/code"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	. "github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/view"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"time"
)

type page struct {
	layout.List
	Manager
	Theme           *material.Theme
	OrigTheme       *material.Theme
	SavedTheme      *material.Theme
	title           string
	btnNavIcon      widget.Clickable
	btnMenuIcon     widget.Clickable
	btnMenuContent  widget.Clickable
	btnSaveSettings widget.Clickable
	btnResetTheme   widget.Clickable
	btnDefaultTheme widget.Clickable
	navIcon         *widget.Icon
	menuIcon        *widget.Icon
	colorpicker.MuxState
	colorpicker.State
	menuVisibilityAnim component.VisibilityAnimation
	langEnum           widget.Enum
	languageCode       code.Code
}

func New(manager Manager) Page {
	navIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	menuIcon, _ := widget.NewIcon(icons.NavigationMoreVert)
	OrigTheme := *user.Theme()
	Theme := *user.Theme()
	SavedTheme := *user.Theme()
	pg := page{
		Manager:    manager,
		Theme:      &Theme,
		OrigTheme:  &OrigTheme,
		SavedTheme: &SavedTheme,
		navIcon:    navIcon,
		List:       layout.List{Axis: layout.Vertical},
		menuIcon:   menuIcon,
		menuVisibilityAnim: component.VisibilityAnimation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		},
	}
	pg.MuxState = colorpicker.NewMuxState([]colorpicker.MuxOption{
		{
			Label: "Contrast Background",
			Value: &pg.Theme.ContrastBg,
		},
		{
			Label: "Contrast Foreground",
			Value: &pg.Theme.ContrastFg,
		},
		{
			Label: "Background",
			Value: &pg.Theme.Bg,
		},
		{
			Label: "Foreground",
			Value: &pg.Theme.Fg,
		},
	}...)
	pg.State.SetColor(*pg.MuxState.Color())
	pg.langEnum.Value = string(*user.LanguageCode())
	pg.languageCode = *user.LanguageCode()
	pg.title = i18n.GetFromCode(key.Theme, pg.languageCode)
	return &pg
}

func (p *page) Layout(gtx Gtx) Dim {
	flex := layout.Flex{Axis: layout.Vertical,
		Spacing:   layout.SpaceEnd,
		Alignment: layout.Start,
	}
	p.title = i18n.GetFromCode(key.Settings, p.languageCode)

	d := flex.Layout(gtx,
		layout.Rigid(p.DrawAppBar),
		layout.Rigid(p.drawContentLayout),
	)
	p.drawMenuLayout(gtx)
	for _, e := range gtx.Queue.Events(p) {
		if e, ok := e.(pointer.Event); ok {
			if e.Type == pointer.Press {
				if !p.btnMenuContent.Pressed() {
					p.menuVisibilityAnim.Disappear(gtx.Now)
				}
			}
		}
	}
	return d
}

func (p *page) DrawAppBar(gtx Gtx) Dim {
	if p.btnNavIcon.Clicked() {
		p.PopUp()
	}
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	gtx.Constraints.Max.Y = gtx.Dp(56)
	th := p.OrigTheme
	return view.DrawAppBarLayout(gtx, th, func(gtx Gtx) Dim {
		return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx Gtx) Dim {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx Gtx) Dim {
						navigationIcon := p.navIcon
						button := material.IconButton(th, &p.btnNavIcon, navigationIcon, "Nav Icon Button")
						button.Size = unit.Dp(40)
						button.Background = th.Palette.ContrastBg
						button.Color = th.Palette.ContrastFg
						button.Inset = layout.UniformInset(unit.Dp(8))
						return button.Layout(gtx)
					}),
					layout.Rigid(func(gtx Gtx) Dim {
						return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx Gtx) Dim {
							titleText := p.title
							title := material.Body1(th, titleText)
							title.Color = th.Palette.ContrastFg
							title.TextSize = unit.Sp(18)
							return title.Layout(gtx)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if p.btnMenuIcon.Clicked() {
					p.menuVisibilityAnim.Appear(gtx.Now)
				}
				button := material.IconButton(th, &p.btnMenuIcon, p.menuIcon, "Context Menu")
				button.Size = unit.Dp(40)
				button.Background = th.Palette.ContrastBg
				button.Color = th.Palette.ContrastFg
				button.Inset = layout.UniformInset(unit.Dp(8))
				d := button.Layout(gtx)
				return d
			}),
		)
	})

}

func (p *page) drawContentLayout(gtx Gtx) Dim {
	th := fonts.NewTheme()
	paint.FillShape(gtx.Ops, th.Bg,
		clip.UniformRRect(image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y), 0).Op(gtx.Ops))
	inset := layout.UniformInset(unit.Dp(16))
	gtx.Constraints.Min = gtx.Constraints.Max
	if p.langEnum.Changed() {
		p.languageCode = code.Code(p.langEnum.Value)
		*user.LanguageCode() = p.languageCode
		user.SaveSettings()
		op.InvalidateOp{}.Add(gtx.Ops)
	}
	if p.MuxState.Changed() {
		p.State.SetColor(*p.MuxState.Color())
	}
	if p.State.Changed() {
		k := p.MuxState.Value
		clr := p.MuxState.Options[k]
		clr.R = p.State.Color().R
		clr.G = p.State.Color().G
		clr.B = p.State.Color().B
		clr.A = p.State.Color().A
		p.State.Editor.SetText(fmt.Sprintf("%02x%02x%02x%02x", clr.R, clr.G, clr.B, clr.A))
	}
	return inset.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			p.List.Axis = layout.Vertical
			return p.List.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
				flex := layout.Flex{Axis: layout.Vertical}
				return flex.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						h := material.H4(th, i18n.GetFromCode(key.Language, p.languageCode))
						h.Alignment = text.Middle
						return h.Layout(gtx)
					}),
					layout.Rigid(func(gtx Gtx) Dim {
						return material.RadioButton(
							th,
							&p.langEnum,
							string(code.English),
							i18n.GetFromCode(key.English, p.languageCode),
						).Layout(gtx)
					}),
					layout.Rigid(func(gtx Gtx) Dim {
						return material.RadioButton(
							th,
							&p.langEnum,
							string(code.Gujarati),
							i18n.GetFromCode(key.Gujarati, p.languageCode),
						).Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
					layout.Rigid(func(gtx Gtx) Dim {
						size := image.Pt(gtx.Constraints.Max.X, gtx.Dp(1))
						bounds := image.Rectangle{Max: size}
						bgColor := color.NRGBA(colornames.Grey500)
						bgColor.A = 75
						paint.FillShape(gtx.Ops, bgColor, clip.UniformRRect(bounds, 0).Op(gtx.Ops))
						return Dim{Size: image.Pt(size.X, size.Y)}
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						h := material.H4(th, i18n.GetFromCode(key.Theme, p.languageCode))
						h.Alignment = text.Middle
						return h.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Min.X = gtx.Dp(50)
								return material.Body1(th, i18n.GetFromCode(key.Red, p.languageCode)).Layout(gtx)
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return material.Slider(th, &p.State.R, 0, 1).Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
							p.drawColorIntBox(p.State.Color().R),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Min.X = gtx.Dp(50)
								return material.Body1(th, i18n.GetFromCode(key.Green, p.languageCode)).Layout(gtx)
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return material.Slider(th, &p.State.G, 0, 1).Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
							p.drawColorIntBox(p.State.Color().G),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Min.X = gtx.Dp(50)
								return material.Body1(th, i18n.GetFromCode(key.Blue, p.languageCode)).Layout(gtx)
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return material.Slider(th, &p.State.B, 0, 1).Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
							p.drawColorIntBox(p.State.Color().B),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Min.X = gtx.Dp(50)
								return material.Body1(th, i18n.GetFromCode(key.Alpha, p.languageCode)).Layout(gtx)
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return material.Slider(th, &p.State.A, 0, 1).Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
							p.drawColorIntBox(p.State.Color().A),
						)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Min.X = gtx.Dp(50)
								return material.Body1(th, i18n.GetFromCode(key.Hex, p.languageCode)).Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx Gtx) Dim {
									return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
										layout.Rigid(func(gtx Gtx) Dim {
											return material.Body1(p.OrigTheme, "#"+p.Editor.Text()).Layout(gtx)
										}),
									)
								})
							}),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.Flex{Alignment: layout.Start, Spacing: layout.SpaceEnd}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								mux := colorpicker.Mux(th, &p.MuxState, "")
								return mux.Layout(gtx)
							}),
						)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(32)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.Y = gtx.Dp(64)
						gtx.Constraints.Min = gtx.Constraints.Max
						background := p.MuxState.Options["Contrast Background"]
						foreground := p.MuxState.Options["Contrast Foreground"]
						paint.FillShape(gtx.Ops, *background, clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Op())
						layout.Center.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								body := material.Body1(p.OrigTheme, "Contrast Foreground")
								body.Color = *foreground
								return body.Layout(gtx)
							},
						)
						return Dim{Size: gtx.Constraints.Max}
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.Y = gtx.Dp(64)
						gtx.Constraints.Min = gtx.Constraints.Max
						background := p.MuxState.Options["Background"]
						foreground := p.MuxState.Options["Foreground"]
						paint.FillShape(gtx.Ops, *background, clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Op())
						layout.Center.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								body := material.Body1(p.OrigTheme, "Foreground")
								body.Color = *foreground
								return body.Layout(gtx)
							},
						)
						return Dim{Size: gtx.Constraints.Max}
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(32)}.Layout),
				)
			})
		},
	)
}

func (p *page) drawColorIntBox(num uint8) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X = gtx.Dp(40)
		gtx.Constraints.Max.Y = gtx.Dp(32)
		gtx.Constraints.Min = gtx.Constraints.Max
		paint.FillShape(gtx.Ops, color.NRGBA(colornames.White), clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Op())
		layout.Center.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				body := material.Body1(p.OrigTheme, fmt.Sprintf("%d", num))
				body.Color = color.NRGBA(colornames.Black)
				return body.Layout(gtx)
			},
		)
		bounds := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
		rect := clip.UniformRRect(bounds, gtx.Dp(4))
		paint.FillShape(gtx.Ops,
			color.NRGBA(colornames.Black),
			clip.Stroke{Path: rect.Path(gtx.Ops), Width: float32(gtx.Dp(1))}.Op(),
		)
		return Dim{Size: gtx.Constraints.Max}
	})
}

func (p *page) drawMenuLayout(gtx Gtx) Dim {
	return layout.Stack{Alignment: layout.NE}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			progress := p.menuVisibilityAnim.Revealed(gtx)
			gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * progress)
			gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * progress)
			return component.Rect{Size: gtx.Constraints.Max, Color: color.NRGBA{}}.Layout(gtx)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			progress := p.menuVisibilityAnim.Revealed(gtx)
			macro := op.Record(gtx.Ops)
			d := p.btnMenuContent.Layout(gtx, p.drawMenuItems)
			call := macro.Stop()
			d.Size.X = int(float32(d.Size.X) * progress)
			d.Size.Y = int(float32(d.Size.Y) * progress)
			component.Rect{Size: d.Size, Color: color.NRGBA(colornames.White)}.Layout(gtx)
			clipOp := clip.Rect{Max: d.Size}.Push(gtx.Ops)
			call.Add(gtx.Ops)
			clipOp.Pop()
			return d
		}),
	)
}

func (p *page) drawMenuItems(gtx Gtx) Dim {
	inset := layout.UniformInset(unit.Dp(12))
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) / 1.5)
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	if p.btnSaveSettings.Clicked() {
		*user.Theme() = *p.Theme
		*p.OrigTheme = *p.Theme
		*user.LanguageCode() = p.languageCode
		user.SaveSettings()
		p.menuVisibilityAnim.Disappear(gtx.Now)
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	if p.btnResetTheme.Clicked() {
		*p.Theme = *p.SavedTheme
		*user.Theme() = *p.Theme
		*p.OrigTheme = *p.Theme
		p.State.SetColor(*p.MuxState.Color())
		user.SaveSettings()
		p.menuVisibilityAnim.Disappear(gtx.Now)
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	if p.btnDefaultTheme.Clicked() {
		*user.Theme() = *fonts.NewTheme()
		*p.Theme = *user.Theme()
		*p.OrigTheme = *user.Theme()
		p.languageCode = code.English
		*user.LanguageCode() = code.English
		p.langEnum.Value = string(code.English)
		user.SaveSettings()
		p.State.SetColor(*p.MuxState.Color())
		p.menuVisibilityAnim.Disappear(gtx.Now)
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btnStyle := material.ButtonLayoutStyle{Button: &p.btnSaveSettings}
			btnStyle.Background = color.NRGBA(colornames.White)
			return btnStyle.Layout(gtx,
				func(gtx Gtx) Dim {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					inset := inset
					return inset.Layout(gtx, func(gtx Gtx) Dim {
						return layout.Flex{Spacing: layout.SpaceEnd}.Layout(gtx,
							layout.Rigid(func(gtx Gtx) Dim {
								bd := material.Body1(p.Theme, "Save")
								bd.Color = color.NRGBA(colornames.Black)
								bd.Alignment = text.Start
								return bd.Layout(gtx)
							}),
						)
					})
				},
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btnStyle := material.ButtonLayoutStyle{Button: &p.btnResetTheme}
			btnStyle.Background = color.NRGBA(colornames.White)
			return btnStyle.Layout(gtx,
				func(gtx Gtx) Dim {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					inset := inset
					return inset.Layout(gtx, func(gtx Gtx) Dim {
						return layout.Flex{Spacing: layout.SpaceEnd}.Layout(gtx,
							layout.Rigid(func(gtx Gtx) Dim {
								bd := material.Body1(p.Theme, "Reset")
								bd.Color = color.NRGBA(colornames.Black)
								bd.Alignment = text.Start
								return bd.Layout(gtx)
							}),
						)
					})
				},
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btnStyle := material.ButtonLayoutStyle{Button: &p.btnDefaultTheme}
			btnStyle.Background = color.NRGBA(colornames.White)
			return btnStyle.Layout(gtx,
				func(gtx Gtx) Dim {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					inset := inset
					return inset.Layout(gtx, func(gtx Gtx) Dim {
						return layout.Flex{Spacing: layout.SpaceEnd}.Layout(gtx,
							layout.Rigid(func(gtx Gtx) Dim {
								bd := material.Body1(p.Theme, "Default")
								bd.Color = color.NRGBA(colornames.Black)
								bd.Alignment = text.Start
								return bd.Layout(gtx)
							}),
						)
					})
				},
			)
		}),
	)
}

func (p *page) URL() URL {
	return SettingsPageURL
}
