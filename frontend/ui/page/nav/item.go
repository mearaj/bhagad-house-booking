package nav

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
)

type pageItem struct {
	fwk.Manager
	*material.Theme
	widget.Clickable
	Title key.Key
	*widget.Icon
	url        fwk.URL
	parentPage *items
}

func (c *pageItem) Layout(gtx fwk.Gtx) fwk.Dim {
	if c.Theme == nil {
		c.Theme = user.Theme()
	}
	shouldReturn := c.URL() == fwk.SearchBookingsPageURL && !(c.parentPage.loginUserResponse.IsAuthorized())
	if shouldReturn {
		return fwk.Dim{}
	}

	return c.layoutContent(gtx)
}

func (c *pageItem) layoutContent(gtx fwk.Gtx) fwk.Dim {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	btnStyle := material.ButtonLayoutStyle{Background: c.Theme.ContrastBg, Button: &c.Clickable}
	if c.Clicked() {
		c.NavigateToURL(c.URL())
	}
	if c.Hovered() || c.URL() == c.CurrentPage().URL() {
		btnStyle.Background.A = 50
	} else {
		btnStyle.Background.A = 10
	}
	d := btnStyle.Layout(gtx, func(gtx fwk.Gtx) fwk.Dim {
		inset := layout.UniformInset(unit.Dp(16))
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		d := inset.Layout(gtx, func(gtx fwk.Gtx) fwk.Dim {
			flex := layout.Flex{Spacing: layout.SpaceEnd, Alignment: layout.Middle}
			d := flex.Layout(gtx,
				layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
					flex := layout.Flex{Spacing: layout.SpaceSides, Alignment: layout.Middle, Axis: layout.Vertical}
					d := flex.Layout(gtx, layout.Rigid(c.drawIcon))
					return d
				}),
				layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
					flex := layout.Flex{Spacing: layout.SpaceSides, Alignment: layout.Start, Axis: layout.Vertical}
					inset := layout.UniformInset(unit.Dp(16))
					d := inset.Layout(gtx, func(gtx fwk.Gtx) fwk.Dim {
						d := flex.Layout(gtx,
							layout.Rigid(func(gtx fwk.Gtx) fwk.Dim {
								title := i18n.Get(c.Title)
								bd := material.Body1(c.Theme, title)
								bd.Font.Weight = text.Bold
								return bd.Layout(gtx)
							}))
						return d
					})
					return d
				}),
			)
			return d
		})
		return d
	})
	return d
}

func (c *pageItem) drawIcon(gtx fwk.Gtx) fwk.Dim {
	gtx.Constraints.Max.X = gtx.Dp(40)
	gtx.Constraints.Max.Y = gtx.Dp(40)
	gtx.Constraints.Min = gtx.Constraints.Max
	if c.Icon == nil {
		return fwk.Dim{Size: gtx.Constraints.Max}
	}

	iconButton := material.IconButton(c.Theme, &widget.Clickable{}, c.Icon, "Booking")
	iconButton.Size = unit.Dp(24)
	iconButton.Background = c.Theme.ContrastBg
	iconButton.Color = c.Theme.Bg
	iconButton.Inset = layout.UniformInset(8)
	return iconButton.Layout(gtx)
}

func (c *pageItem) URL() fwk.URL {
	return c.url
}
