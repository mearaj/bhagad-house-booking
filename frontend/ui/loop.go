package ui

import (
	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/mearaj/bhagad-house-booking/common/alog"
	. "github.com/mearaj/bhagad-house-booking/frontend/ui/fwk"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"image"
	"os/exec"
	"strings"
	"time"
)

// FixTimezone https://github.com/golang/go/issues/20455
func FixTimezone() {
	out, err := exec.Command("/system/bin/getprop", "persist.sys.timezone").Output()
	if err != nil {
		return
	}
	z, err := time.LoadLocation(strings.TrimSpace(string(out)))
	if err != nil {
		return
	}
	time.Local = z
}

func init() {
	go FixTimezone()
}

func Loop(w *app.Window) error {
	var ops op.Ops
	appManager := NewAppManager()
	appManager.window = w

	// backClickTag is meant for tracking user's backClick action, specially on mobile
	var backClickTag struct{}

	//subscription := appManager.Service().Subscribe()

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				alog.Logger().Errorln("system.DestroyEvent called", e.Err)
				return e.Err
			case system.FrameEvent:
				appManager.Insets = e.Insets
				e.Insets = system.Insets{}
				gtx := layout.NewContext(&ops, e)
				for _, event := range gtx.Events(&backClickTag) {
					if e, ok := event.(key.Event); ok {
						if e.Name == key.NameBack {
							if len(appManager.pagesStack) > 1 {
								appManager.PopUp()
							}
						}
					}
				}
				// Listen to back command only when appManager.pagesStack is greater than 1,
				//  so we can pop up page else we want the android's default behavior
				if len(appManager.pagesStack) > 1 {
					key.InputOp{Tag: &backClickTag, Keys: key.NameBack}.Add(gtx.Ops)
				}
				appManager.Constraints = gtx.Constraints
				appManager.Metric = gtx.Metric
				// Create a clip area the size of the window.
				areaStack := clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Push(gtx.Ops)
				// In desktop layout, sidebar exists and needs to listen to entire window's pointer event
				// hence added here. It avoids conflict with page that contains sidebar
				for _, elem := range []interface{}{appManager.CurrentPage(), appManager.settingsSideBar} {
					pointer.InputOp{
						Types: pointer.Enter | pointer.Leave | pointer.Drag | pointer.Press | pointer.Release | pointer.Scroll | pointer.Move,
						Tag:   elem,
					}.Add(gtx.Ops)
				}
				layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Flexed(1, func(gtx Gtx) Dim {
						size := image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y}
						bounds := image.Rectangle{Max: size}
						paint.FillShape(gtx.Ops, user.Theme().Bg, clip.UniformRRect(bounds, 0).Op(gtx.Ops))
						return appManager.Layout(gtx)
					}),
				)
				areaStack.Pop()
				e.Frame(gtx.Ops)
			case system.StageEvent:
				if e.Stage == system.StagePaused {
					alog.Logger().Infoln("window is running in background")
					appManager.isStageRunning = false
				} else if e.Stage == system.StageRunning {
					alog.Logger().Infoln("window is running in foreground")
					appManager.isStageRunning = true
				}
			}
			//case event := <-subscription.Events():
			//	var settingsBarFound bool
			//	for _, eachPage := range appManager.pagesStack {
			//		if l, ok := eachPage.(ServiceListener); ok {
			//			l.OnServiceStateChange(event)
			//		}
			//		settingsBarFound = eachPage == appManager.settingsSideBar
			//	}
			//	if l, ok := appManager.settingsSideBar.(ServiceListener); ok && !settingsBarFound {
			//		l.OnServiceStateChange(event)
			//	}
		}
	}
}
