package i18n

import (
	"encoding/json"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func InitBundle() *goi18n.Bundle {
	bundle := goi18n.NewBundle(language.Russian) // язык по умолчанию
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	return bundle
}
