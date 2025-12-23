package handler

import (
    "embed"
    "encoding/json"
    "io/fs"
    "net/http"
    "os"
    "path"
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
    fileHandler := ServeFS(adminAssets)

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if path.Clean(r.URL.Path) == "/env.runtime.json" {
            payload, _ := json.Marshal(runtimeEnvFromEnv())
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            w.Write(payload)
            return
        }

        fileHandler.ServeHTTP(w, r)
    })
}

// runtimeEnvFromEnv builds the runtime env payload from process environment.
// Empty values are encoded as null to allow client-side fallbacks.
func runtimeEnvFromEnv() map[string]*string {
    envOrNil := func(key string) *string {
        if val, ok := os.LookupEnv(key); ok && val != "" {
            return &val
        }
        return nil
    }

    return map[string]*string{
        "NEXT_PUBLIC_BRAND_PRODUCT_NAME":    envOrNil("NEXT_PUBLIC_BRAND_PRODUCT_NAME"),
        "NEXT_PUBLIC_BRAND_COLOR":           envOrNil("NEXT_PUBLIC_BRAND_COLOR"),
        "NEXT_PUBLIC_GRAPHQL_API_URL":       envOrNil("NEXT_PUBLIC_GRAPHQL_API_URL"),
        "NEXT_PUBLIC_GRAPHQL_QUERY_API_URL": envOrNil("NEXT_PUBLIC_GRAPHQL_QUERY_API_URL"),
        "NEXT_PUBLIC_GRAPHQL_QUERY_PLAYGROUND_URL": envOrNil("NEXT_PUBLIC_GRAPHQL_QUERY_PLAYGROUND_URL"),
        "NEXT_PUBLIC_GRAPHQL_ADMIN_API_URL":       envOrNil("NEXT_PUBLIC_GRAPHQL_ADMIN_API_URL"),
        "NEXT_PUBLIC_GRAPHQL_ADMIN_PLAYGROUND_URL": envOrNil("NEXT_PUBLIC_GRAPHQL_ADMIN_PLAYGROUND_URL"),
        "NEXT_PUBLIC_CENTRIFUGO_URL":              envOrNil("NEXT_PUBLIC_CENTRIFUGO_URL"),
    }
}
