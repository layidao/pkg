// Copyright (c) 2020. pkg Inc. All rights reserved.
// Author bozz
// Create Time 2020/12/7

package i18n

import (
	"github.com/unknwon/i18n"
)

func New(locales map[string]string) error {
	//i18n.SetMessage("en-US", "conf/locale/locale_en-US.ini")
	//i18n.SetMessage("zh-CN", "conf/locale/locale_zh-CN.ini")
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
