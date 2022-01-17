package home

import (
	"embed"
	_ "embed"
	"log"
	"net/http"
)

//go:embed index.html
var homePage []byte

// Handler serves the all-in-one home page
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	if _, err := w.Write(homePage); err != nil {
		log.Printf("Error writing response to %s: %s", r.RemoteAddr, err)
	}
}

//go:embed static
var static embed.FS

var Static = http.FileServer(http.FS(static))
