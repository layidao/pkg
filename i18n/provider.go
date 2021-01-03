// Copyright (c) 2020. pkg Inc. All rights reserved.
// Author bozz@stc.plus
// Create Time 2020/12/10

package i18n

import (
	"github.com/unknwon/i18n"
)

func Init(locales map[string]string) error {
	for lang, localeFile := range locales {
		if err := i18n.SetMessage(lang, localeFile); err != nil {
			return err
		}
	}
	return nil
}

func Tr(lang, key string) string {
	return i18n.Tr(lang, key)
}
