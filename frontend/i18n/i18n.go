package i18n

import (
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/code"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/internal/keys"
	"github.com/mearaj/bhagad-house-booking/frontend/i18n/key"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
)

type TranslationMap = utils.Map[key.Key, string]

var translator = utils.NewFromMap(map[code.Code]TranslationMap{
	code.English:  utils.NewFromMap(keys.EnMapNative),
	code.Gujarati: utils.NewFromMap(keys.GuMapNative),
})

func init() {

}

func Get(langKey key.Key) string {
	return GetFromCode(langKey, *user.LanguageCode())
}

func GetFromCode(langKey key.Key, langCode code.Code) string {
	lang, ok := translator.Get(langCode)
	if !ok {
		return string(langKey)
	}
	val, ok := lang.Get(langKey)
	if !ok {
		return string(langKey)
	}
	return val
}
