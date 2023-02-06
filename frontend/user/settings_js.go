//go:build js

package user

import (
	"errors"
	"github.com/mearaj/bhagad-house-booking/frontend/assets/fonts"
	log "github.com/sirupsen/logrus"
	"syscall/js"
)

const LocalStorage = "localStorage"

func LoadSettings() {
	localStorage := js.Global().Get(LocalStorage)
	var err error
	if localStorage.IsNull() || localStorage.IsUndefined() {
		err = errors.New(LocalStorage + " is not defined")
		log.Errorln(err)
		fallbackToDefault()
		return
	}
	data := localStorage.Call("getItem", SettingsFileName)
	if data.IsNull() || data.IsUndefined() {
		err = errors.New("settings were not saved to local storage")
		log.Errorln(err)
		fallbackToDefault()
		return
	}
	settingsJSON, err := unmarshalJSON([]byte(data.String()))
	if err != nil {
		log.Errorln(err)
		fallbackToDefault()
		return
	}
	th := fonts.NewTheme()
	bg := settingsJSON.FlatTheme.Background
	fg := settingsJSON.FlatTheme.Foreground
	ctBg := settingsJSON.FlatTheme.ContrastBackground
	ctFg := settingsJSON.FlatTheme.ContrastForeground
	th.Fg, th.Bg, th.ContrastBg, th.ContrastFg = fg, bg, ctBg, ctFg
	user := settingsJSON.User
	rate := settingsJSON.BookingRate
	settingsMutex.Lock()
	settings = Settings{
		languageCode: settingsJSON.LanguageCode,
		theme:        *th,
		user:         user,
		bookingRate:  rate,
	}
	settingsMutex.Unlock()
}
func SaveSettings() {
	var err error
	data, err := marshalJSON()
	if err != nil {
		fallbackToDefault()
		log.Errorln(err)
		return
	}
	localStorage := js.Global().Get(LocalStorage)
	if localStorage.IsNull() || localStorage.IsUndefined() {
		err = errors.New(LocalStorage + " is not defined")
		log.Errorln(err)
		fallbackToDefault()
		return
	}
	localStorage.Call("setItem", SettingsFileName, string(data))
}
