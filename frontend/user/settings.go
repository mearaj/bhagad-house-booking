package user

import (
	"encoding/json"
	"gioui.org/widget/material"
	"github.com/mearaj/bhagad-house-booking/common/alog"
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
	LanguageCode code.Code `json:"languageCode"`
	FlatTheme    flatTheme `json:"flatTheme"`
}

type Settings struct {
	languageCode code.Code
	theme        material.Theme
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
	st := jsonSettings{
		LanguageCode: *lang,
		FlatTheme: flatTheme{
			Foreground:         th.Fg,
			ContrastForeground: th.ContrastFg,
			Background:         th.Bg,
			ContrastBackground: th.ContrastBg,
		},
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

func fallbackToDefault(err error) {
	settingsMutex.Lock()
	settings = Settings{
		languageCode: code.English,
		theme:        *fonts.NewTheme(),
	}
	settingsMutex.Unlock()
	if err != nil {
		alog.Logger().Println(err)
	}
}
