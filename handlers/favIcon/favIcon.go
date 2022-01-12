package favIcon

import (
	_ "embed"
	"net/http"
)

//go:embed favIcon.go
var icon []byte

func Icon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	w.Write(icon)
}
