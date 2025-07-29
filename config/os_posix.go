//go:build linux || darwin || freebsd

package config

// import "github.com/mechiko/telebot/entity"
// if !entity.supported
var (
	DbPath               = "/var/local/kminfo"
	LogPath              = "/var/log/kminfo"
	ConfigPath           = "/etc/kminfo"
	Supported            = true
	Linux                = true
	Windows              = false
	PosixUserUIDGUID int = 1002
	PosixChownPath   int = 0755
	PosixChownFile   int = 0644
)
