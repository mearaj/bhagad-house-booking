package user

import (
	"encoding/json"
	"gioui.org/widget/material"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"github.com/mearaj/bhagad-house-booking/frontend/assets/fonts"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/code"
	"image/color"
	"sync"
)

const SettingsFileName = "settings.json"
const AppDirName = "bhagad_house_booking"

var settingsMutex sync.RWMutex
var settings Settings

type flatTheme struct {
	Foreground         color.NRGBA `json:"foreground"`
	ContrastForeground color.NRGBA `json:"contrastForeground"`
	Background         color.NRGBA `json:"background"`
	ContrastBackground color.NRGBA `json:"contrastBackground"`
}

type jsonSettings struct {
	LanguageCode code.Code          `json:"languageCode"`
	FlatTheme    flatTheme          `json:"flatTheme"`
	User         response.LoginUser `json:"user"`
	BookingRate  float64            `json:"booking_rate"`
}

type Settings struct {
	languageCode code.Code
	theme        material.Theme
	user         response.LoginUser
	bookingRate  float64
}

func init() {
	LoadSettings()
}

func marshalJSON() ([]byte, error) {
	th := Theme()
	lang := LanguageCode()
	if th == nil {
		th = fonts.NewTheme()
	}
	user := User()
	rate := BookingRate()
	st := jsonSettings{
		LanguageCode: *lang,
		FlatTheme: flatTheme{
			Foreground:         th.Fg,
			ContrastForeground: th.ContrastFg,
			Background:         th.Bg,
			ContrastBackground: th.ContrastBg,
		},
		User:        *user,
		BookingRate: rate,
	}
	return json.MarshalIndent(&st, "", "  ")
}

func unmarshalJSON(data []byte) (jsonSettings, error) {
	var res jsonSettings
	return res, json.Unmarshal(data, &res)
}

func LanguageCode() *code.Code {
	settingsMutex.RLock()
	defer settingsMutex.RUnlock()
	langCode := &settings.languageCode
	return langCode
}

func Theme() *material.Theme {
	settingsMutex.RLock()
	defer settingsMutex.RUnlock()
	theme := &settings.theme
	return theme
}

func User() *response.LoginUser {
	settingsMutex.RLock()
	defer settingsMutex.RUnlock()
	return &settings.user
}

func BookingRate() float64 {
	settingsMutex.RLock()
	defer settingsMutex.RUnlock()
	return settings.bookingRate
}

func SetBookingRate(rate float64) {
	settingsMutex.Lock()
	settings.bookingRate = rate
	settingsMutex.Unlock()
	SaveSettings()
}

func SetUser(user response.LoginUser) {
	settingsMutex.Lock()
	settings.user = user
	settingsMutex.Unlock()
	SaveSettings()
}

func fallbackToDefault() {
	settingsMutex.Lock()
	settings = Settings{
		languageCode: code.English,
		theme:        *fonts.NewTheme(),
	}
	settingsMutex.Unlock()
}
