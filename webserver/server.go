package webserver

//"www_certificate": "certificates/src_certs_server.pem",
//"www_certificate_key": "certificates/src_certs_server.key",
//"www_certificate": "certificates/public.crt",
//"www_certificate_key": "certificates/private.key",
import (
	"awi/config"
	"awi/handlers/home"
	"awi/handlers/webhooks"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

type Server struct {
	name    string
	version string

	Index  int
	Delay  int
	config *config.Config
}

// NewServer returns new Server.
func New(name, version string, config *config.Config) *Server {
	return &Server{
		name:    name,
		version: version,
		config:  config,
	}
}

//const stabilizationTime = 5 * time.Minute / time.Second
const stabilizationTime = (15 * time.Second) / time.Second

var timeLeft = int32(stabilizationTime)

func getCountdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if atomic.LoadInt32(&timeLeft) > 0 {
		atomic.AddInt32(&timeLeft, -1)
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(&timeLeft); err != nil {
		log.Printf("encoding error with data <%s> : %s\n", b.String(), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b.Bytes()); err != nil {
		log.Printf("response error with data %s : %s\n", b.String(), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func resetTimer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if atomic.LoadInt32(&timeLeft) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	atomic.StoreInt32(&timeLeft, int32(stabilizationTime))
	if _, err := w.Write([]byte("{}")); err != nil {
		log.Printf("response error with empty data: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// ListenAndServe listens on the TCP address and serves requests.
//func (s *Server) ListenAndServe() error {
func (s *Server) ListenAndServeHTTPS() {
	bind := fmt.Sprintf("%s:%s", s.config.WWWAddr, s.config.WWWPort)
	fmt.Printf("Веб-сервер %s [%s] - 'httpS' запущен %s\n", s.name, s.version, bind)

	http.HandleFunc("/", home.Handler())
	http.HandleFunc("/index.js", home.Scripts())
	http.HandleFunc("/style.css", home.Styles())
	http.HandleFunc("/favicon.ico", home.Favicon())
	http.HandleFunc("/dseg7.woff2", home.DSEG7())
	http.HandleFunc("/webhooks", webhooks.WebHooksHandler)

	http.HandleFunc("/get/countdown", getCountdown)
	http.HandleFunc("/post/reset-timer", resetTimer)

	// Запуск веб-сервера
	rootDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	certFile := fmt.Sprintf("%s/%s", rootDir, s.config.WWWCertificate)
	keyFile := fmt.Sprintf("%s/%s", rootDir, s.config.WWWCertificateKey)
	if err := http.ListenAndServeTLS(bind, certFile, keyFile, nil); err != nil {
		//if err := http.ListenAndServeTLS(s.Bind, s.Certificate, s.CertificateKey, nil); err != nil {
		log.Fatal("HTTPS-Err: ", err)
	}
}
