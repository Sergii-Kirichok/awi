package webserver

//"www_certificate": "certificates/src_certs_server.pem",
//"www_certificate_key": "certificates/src_certs_server.key",
//"www_certificate": "certificates/public.crt",
//"www_certificate_key": "certificates/private.key",
import (
	"awi/config"
	"awi/controller"
	"awi/handlers/home"
	"awi/handlers/webhooks"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Server struct {
	name    string
	version string
	router  *mux.Router

	Index      int
	Delay      int
	config     *config.Config
	controller *controller.Controller
}

// NewServer returns new Server.
func New(name, version string, config *config.Config, control *controller.Controller) *Server {
	return &Server{
		name:       name,
		version:    version,
		config:     config,
		router:     mux.NewRouter(),
		controller: control,
	}
}

func (s *Server) getZoneData(w http.ResponseWriter, r *http.Request) {
	zone, err := s.controller.GetZoneData(mux.Vars(r)["zone-id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := sendJSON(w, zone); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) buttonPress(w http.ResponseWriter, r *http.Request) {
	if err := s.controller.MakeAction(mux.Vars(r)["zone-id"]); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func sendJSON(w http.ResponseWriter, data interface{}) error {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(&data); err != nil {
		return err
	}

	if _, err := w.Write(b.Bytes()); err != nil {
		return err
	}

	return nil
}

// ListenAndServeHTTPS listens on the TCP address and serves requests.
func (s *Server) ListenAndServeHTTPS() {
	bind := fmt.Sprintf("%s:%s", s.config.WWWAddr, s.config.WWWPort)
	fmt.Printf("Веб-сервер %s [%s] - 'httpS' запущен %s\n", s.name, s.version, bind)

	s.router.PathPrefix("/static").Handler(home.Static)

	wh := webhooks.NewHandler(s.config, s.controller)
	s.router.HandleFunc("/webhooks", wh.WebHooksHandler).Methods(http.MethodPost)

	zone := s.router.PathPrefix("/zones/{zone-id}").Subrouter()
	zone.HandleFunc("", home.Handler).Methods(http.MethodGet)
	zone.HandleFunc("/data", s.getZoneData).Methods(http.MethodGet)
	zone.HandleFunc("/button-press", s.buttonPress).Methods(http.MethodGet)

	// Запуск веб-сервера
	rootDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	certFile := fmt.Sprintf("%s/%s", rootDir, s.config.WWWCertificate)
	keyFile := fmt.Sprintf("%s/%s", rootDir, s.config.WWWCertificateKey)
	if err := http.ListenAndServeTLS(bind, certFile, keyFile, s.router); err != nil {
		//if err := http.ListenAndServeTLS(s.Bind, s.Certificate, s.CertificateKey, nil); err != nil {
		log.Fatal("HTTPS-Err: ", err)
	}
}
