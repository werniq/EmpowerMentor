package utils

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// TranslateToUkrainian function is used to translate text to Ukrainian
func TranslateToUkrainian(text string) string {
	ukrainian := language.Ukrainian
	translator := message.NewPrinter(ukrainian)
	return translator.Sprintf(text)
}
