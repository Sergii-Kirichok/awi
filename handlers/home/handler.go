package home

import (
	_ "embed"
	"log"
	"net/http"
)

//go:embed index.html
var homePage []byte

// Handler serves the all-in-one home page
func Handler() http.HandlerFunc {
	return getStatic(homePage, "text/html")
}

//go:embed index.js
var scripts []byte

func Scripts() http.HandlerFunc {
	return getStatic(scripts, "application/javascript")
}

//go:embed style.css
var styles []byte

func Styles() http.HandlerFunc {
	return getStatic(styles, "text/css")
}

//go:embed favicon.ico
var favicon []byte

func Favicon() http.HandlerFunc {
	return getStatic(favicon, "image/x-icon")
}

//go:embed dseg7.woff2
var dseg7 []byte

func DSEG7() http.HandlerFunc {
	return getStatic(dseg7, "font/woff2")
}

func getStatic(data []byte, mimeType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("From %s %s %s", r.RemoteAddr, r.Method, r.URL.Path)

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", mimeType)

		if _, err := w.Write(data); err != nil {
			log.Printf("Error writing response to %s: %s", r.RemoteAddr, err)
		}
	}
}
