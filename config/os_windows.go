//go:build windows

package config

import (
	"korrectkm/utility"
	"path/filepath"
)

// если каталог ../spaserver/templates в режиме разработки существует, то прописываем его в переменную
// для поиска шаблонов динамической обработки для отладки
var rootPathTemplates = "../spaserver/templates"

func RootPathTemplates() string {
	if Mode == "development" {
		absPath, err := filepath.Abs(rootPathTemplates)
		if err != nil {
			return ""
		}
		if utility.PathOrFileExists(absPath) {
			return absPath
		} else {
			return ""
		}
	} else {
		return ""
	}
}

var (
	ConfigPath       = ".korrectkm"
	DbPath           = ".korrectkm"
	LogPath          = ".korrectkm"
	Supported        = true
	Windows          = true
	Linux            = false
	PosixUserUIDGUID = 1002
	PosixChownPath   = 0755
	PosixChownFile   = 0644
)
