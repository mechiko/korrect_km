package templates

import (
	"embed"
)

//go:embed header footer home index kmstate
var root embed.FS
