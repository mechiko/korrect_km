package templates

import (
	"embed"
)

//go:embed header footer home root kmstate
var root embed.FS
