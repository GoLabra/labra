package handler

import (
    "embed"
    "net/http"
    "io/fs"
	"github.com/GoLabra/labra/handler"
)

// ...existing code...
//go:embed labradmin/index.html
var indexHTML string

//go:embed labradmin/**
var adminFiles embed.FS

var adminAssets fs.FS

func init() {
    var err error
    adminAssets, err = fs.Sub(adminFiles, "labradmin")
    if err != nil {
        panic(err)
    }
}

func ServeAdmin() http.Handler {
    return handler.ServeFS(adminAssets)
}
