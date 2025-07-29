package embedded

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed root
var Root embed.FS

func GetFileSystem() http.FileSystem {
	fsys, err := fs.Sub(Root, "root")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
