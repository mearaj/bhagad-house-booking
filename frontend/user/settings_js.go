//go:build js

package user

import (
	"errors"
	"github.com/mearaj/bhagad-house-booking/frontend/assets/fonts"
	"syscall/js"
)

const LocalStorage = "localStorage"

func LoadSettings() {
	localStorage := js.Global().Get(LocalStorage)
	var err error
	if localStorage.IsNull() || localStorage.IsUndefined() {
		err = errors.New(LocalStorage + " is not defined")
		fallbackToDefault(err)
		return
	}
	data := localStorage.Call("getItem", SettingsFileName)
	if data.IsNull() || data.IsUndefined() {
		err = errors.New("settings not saved to local storage")
		fallbackToDefault(err)
		return
	}
	settingsJSON, err := unmarshalJSON([]byte(data.String()))
	if err != nil {
		fallbackToDefault(err)
		return
	}
	th := fonts.NewTheme()
	bg := settingsJSON.FlatTheme.Background
	fg := settingsJSON.FlatTheme.Foreground
	ctBg := settingsJSON.FlatTheme.ContrastBackground
	ctFg := settingsJSON.FlatTheme.ContrastForeground
	th.Fg, th.Bg, th.ContrastBg, th.ContrastFg = fg, bg, ctBg, ctFg
	settingsMutex.Lock()
	settings = Settings{
		languageCode: settingsJSON.LanguageCode,
		theme:        *th,
	}
	settingsMutex.Unlock()
}
func SaveSettings() {
	var err error
	data, err := marshalJSON()
	if err != nil {
		fallbackToDefault(err)
		return
	}
	localStorage := js.Global().Get(LocalStorage)
	if localStorage.IsNull() || localStorage.IsUndefined() {
		err = errors.New(LocalStorage + " is not defined")
		fallbackToDefault(err)
		return
	}
	localStorage.Call("setItem", SettingsFileName, string(data))
}
