package webserver

//"www_certificate": "certificates/src_certs_server.pem",
//"www_certificate_key": "certificates/src_certs_server.key",
//"www_certificate": "certificates/public.crt",
//"www_certificate_key": "certificates/private.key",
import (
	"awi/config"
	"awi/handlers/favIcon"
	"awi/handlers/webhooks"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Server struct {
	Name    string
	Version string

	Index int
	Delay int
	Conf  *config.Config
}

// NewServer returns new Server.
func NewServer() *Server {
	return &Server{}
}

// ListenAndServe listens on the TCP address and serves requests.
//func (s *Server) ListenAndServe() error {
func (s *Server) ListenAndServeHTTPS() {
	bind := fmt.Sprintf("%s:%s", s.Conf.WWWAddr, s.Conf.WWWPort)
	fmt.Printf("Веб-сервер %s [%s] - 'httpS' запущен %s\n", s.Name, s.Version, bind)
	http.HandleFunc("/favicon.ico", favIcon.Icon)
	http.HandleFunc("/webhooks", webhooks.WebHooksHandler)

	// Запуск веб-сервера
	rootDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	if err := http.ListenAndServeTLS(bind, fmt.Sprintf("%s/%s", rootDir, s.Conf.WWWCertificate), fmt.Sprintf("%s/%s", rootDir, s.Conf.WWWCertificateKey), nil); err != nil {
		//if err := http.ListenAndServeTLS(s.Bind, s.Certificate, s.CertificateKey, nil); err != nil {
		log.Fatal("HTTPS-Err: ", err)
	}
}
