//go:build !js

package user

import (
	"gioui.org/app"
	"github.com/mearaj/bhagad-house-booking/frontend/assets/fonts"
	"os"
	"path/filepath"
)

func LoadSettings() {
	dirPath, err := app.DataDir()
	if err != nil {
		fallbackToDefault(err)
		return
	}
	dirPath = filepath.Join(dirPath, AppDirName)
	if _, err = os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0700)
		if err != nil {
			fallbackToDefault(err)
			return
		}
	}
	filePath := filepath.Join(dirPath, SettingsFileName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		fallbackToDefault(err)
		return
	}
	settingsJSON, err := unmarshalJSON(data)
	if err != nil {
		fallbackToDefault(err)
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
	dirPath, err := app.DataDir()
	if err != nil {
		fallbackToDefault(err)
		return
	}
	dirPath = filepath.Join(dirPath, AppDirName)
	if _, err = os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0700)
		if err != nil {
			fallbackToDefault(err)
			return
		}
	}
	filePath := filepath.Join(dirPath, SettingsFileName)
	data, err := marshalJSON()
	if err != nil {
		fallbackToDefault(err)
		return
	}
	err = os.WriteFile(filePath, data, 0600)
	if err != nil {
		fallbackToDefault(err)
	}
}
