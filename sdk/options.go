package sdk

import "fmt"

type Language string

const (
	English Language = "en"
	Welsh   Language = "cy"
)

type Options struct {
	Offset int
	Limit  int
	Lang   Language
}

// ErrUnrecognisedLanguage builds error message when the language is not recognisable
func ErrUnrecognisedLanguage(lang Language) error {
	return fmt.Errorf("unrecognised language: %s", lang)
}

func (lang Language) String() (string, error) {
	switch lang {
	case English:
		return string(lang), nil
	case Welsh:
		return string(lang), nil
	}

	return "", ErrUnrecognisedLanguage(lang)
}
