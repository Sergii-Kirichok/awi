package webhooks

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
)

//go:embed styles.css
var stylesCSS []byte

var exampleUpload []byte

func WebHooksHandler(w http.ResponseWriter, r *http.Request) {

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(b))
	w.WriteHeader(http.StatusOK)
}
