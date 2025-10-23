package stats

import (
	_ "embed"
)

//go:embed templates\index.html
var indexTmpl string

//go:embed templates\navbar.html
var navbarTmpl string

//go:embed templates\page.html
var pageTmpl string

//go:embed templates\footer.html
var footerTmpl string

//go:embed templates\progress.html
var progressTmpl string

//go:embed templates\chunk.html
var chunkTmpl string
