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

	m := NewMessage()
	err = m.Parse(b)
	if err != nil {
		log.Printf("Error [WebHooksHandler]: %s\n", err)
	}
	//fmt.Println(string(b))
	fmt.Printf("Notification parsed [%s], has [%d] events:\n", m.Type, len(m.Notifications))
	for k, e := range m.Notifications {
		fmt.Printf("%02d.Event [%s]=>[%s]:\n\n", k, e.Type, e.Event.Type)
	}

	w.WriteHeader(http.StatusOK)
}
