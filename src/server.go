package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func RootHandler(fileSystem http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Clear-Size-Data", "*")
		log.Printf("Page %v Visited\n", request.URL.Path)
		fileSystem.ServeHTTP(writer, request)
	}
}

func EmbededServer(dir embed.FS, location string) {
	EmbedDir, _ := fs.Sub(dir, "win")
	http.Handle("/", RootHandler(http.FileServer(http.FS(EmbedDir))))
	log.Printf("Serving on HTTP location \"%v\"\n", location)
	log.Fatal(http.ListenAndServe(location, nil))
}

func DevServer(dir, location string) {
	AbsoluteDir, _ := filepath.Abs(dir)
	os.Chdir(AbsoluteDir)
	http.Handle("/", RootHandler(http.FileServer(http.Dir(AbsoluteDir))))
	log.Printf("Serving on HTTP location \"%v\"\n", location)
	log.Fatal(http.ListenAndServe(location, nil))
}