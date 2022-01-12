package webserver

//"www_certificate": "certificates/src_certs_server.pem",
//"www_certificate_key": "certificates/src_certs_server.key",
//"www_certificate": "certificates/public.crt",
//"www_certificate_key": "certificates/private.key",
import (
	"awi/config"
	"awi/handlers/home"
	"awi/handlers/webhooks"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	// Запуск веб-сервера
	rootDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	certFile := fmt.Sprintf("%s/%s", rootDir, s.config.WWWCertificate)
	keyFile := fmt.Sprintf("%s/%s", rootDir, s.config.WWWCertificateKey)
	if err := http.ListenAndServeTLS(bind, certFile, keyFile, nil); err != nil {
		//if err := http.ListenAndServeTLS(s.Bind, s.Certificate, s.CertificateKey, nil); err != nil {
		log.Fatal("HTTPS-Err: ", err)
	}
}
