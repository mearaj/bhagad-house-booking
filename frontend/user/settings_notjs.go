//go:build !js

package user

import (
	"gioui.org/app"
	"github.com/mearaj/bhagad-house-booking/frontend/assets/fonts"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func LoadSettings() {
	dirPath, err := app.DataDir()
	if err != nil {
		log.Errorln(err)
		fallbackToDefault()
		return
	}
	dirPath = filepath.Join(dirPath, AppDirName)
	if _, err = os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0700)
		if err != nil {
			log.Errorln(err)
			fallbackToDefault()
			return
		}
	}
	filePath := filepath.Join(dirPath, SettingsFileName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Errorln(err)
		fallbackToDefault()
		return
	}
	settingsJSON, err := unmarshalJSON(data)
	if err != nil {
		fallbackToDefault()
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
	dirPath, err := app.DataDir()
	if err != nil {
		log.Errorln(err)
		fallbackToDefault()
		return
	}
	dirPath = filepath.Join(dirPath, AppDirName)
	if _, err = os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0700)
		if err != nil {
			log.Errorln(err)
			fallbackToDefault()
			return
		}
	}
	filePath := filepath.Join(dirPath, SettingsFileName)
	data, err := marshalJSON()
	if err != nil {
		log.Errorln(err)
		fallbackToDefault()
		return
	}
	err = os.WriteFile(filePath, data, 0600)
	if err != nil {
		log.Errorln(err)
		fallbackToDefault()
	}
}
