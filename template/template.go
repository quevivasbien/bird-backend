package template

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed all:build
var content embed.FS

// Return contents of build directory
func GetBuild() http.FileSystem {
	dist, err := fs.Sub(content, "build")
	if err != nil {
		log.Fatal(err)
	}
	return http.FS(dist)
}
