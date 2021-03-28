package bundle

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	chineseI18n = i18n.Message{
		ID:    "Emails",
		One:   "{{.Name}} 您好! 歡迎加入 AmazingTalker",
		Other: "{{.Name}} 您好! 歡迎加入 AmazingTalker",
	}

	englishI18n = i18n.Message{
		ID:    "Emails",
		One:   "Hi {{.Name}} ! Welcome to AmazingTalker.",
		Other: "Hi {{.Name}} ! Welcome to AmazingTalker.",
	}
)

// NewBundle ...
func NewBundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.Chinese)
	bundle.AddMessages(language.Chinese, &chineseI18n)
	bundle.AddMessages(language.English, &englishI18n)

	return bundle
}
