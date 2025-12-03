package handler

import (
    "io"
    "io/fs"
    "net/http"
    "path"
    //"strings"
)

// ServeFS serves files from a provided fs.FS under the given urlPrefix.
// Behavior: serves files, serves directory/index.html if present, falls back to root index.html for SPA routes.
func ServeFS(fsys fs.FS) http.Handler {

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rel := r.URL.Path
        rel = path.Clean("/" + rel)[1:]

        if rel == "" || rel == "." {
            // root -> serve root/index.html
            serveFileIfExists(w, r, http.FS(fsys), "index.html")
            return
        }

        if f, err := fsys.Open(rel); err == nil {
            defer f.Close()
            if fi, err := f.Stat(); err == nil {
                if fi.IsDir() {
                    // if directory, try directory/index.html
                    idx := path.Join(rel, "index.html")
                    if serveFileIfExists(w, r, http.FS(fsys), idx) {
                        return
                    }
					
                    http.NotFound(w, r)
                    return
                }
                // it's a file -> serve it (use FileServer for range/headers)
				serveFileIfExists(w, r, http.FS(fsys), rel) 
                return
            }
        }

        http.NotFound(w, r)
    })
}


// helper: open file from fs and serve with http.ServeContent. returns true if served.
func serveFileIfExists(w http.ResponseWriter, r *http.Request, fsys http.FileSystem, name string) bool {
    f, err := fsys.Open(name)
    if err != nil {
        return false
    }
    defer f.Close()

    fi, err := f.Stat()
    if err != nil || fi.IsDir() {
        return false
    }

    // use http.ServeContent to get correct headers (modtime, range support)
    http.ServeContent(w, r, name, fi.ModTime().UTC(), f.(io.ReadSeeker))
    return true
}