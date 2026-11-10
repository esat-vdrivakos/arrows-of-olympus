package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dir := os.Getenv("PUBLIC_DIR")
	if dir == "" {
		dir = "public"
	}

	http.Handle("/", gzipStatic(dir, http.FileServer(http.Dir(dir))))
	log.Printf("Arrows of Olympus WebAssembly server listening on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// gzipStatic serves a precompressed "<file>.gz" sibling when the client accepts gzip.
// game.wasm is ~24 MB raw and ~6 MB gzipped, and Cloud Run does not compress responses.
func gzipStatic(dir string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		if strings.HasSuffix(r.URL.Path, ".wasm") && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gz := filepath.Join(dir, filepath.Clean("/"+r.URL.Path)+".gz")
			if _, err := os.Stat(gz); err == nil {
				w.Header().Set("Content-Type", "application/wasm")
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Add("Vary", "Accept-Encoding")
				http.ServeFile(w, r, gz)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
