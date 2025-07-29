package templates

import (
	"embed"
)

//go:embed header footer home index setup
var root embed.FS
