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
	"math/rand"
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

func (s *Server) getZoneName(w http.ResponseWriter, r *http.Request) {
	zone := s.controller.GetZoneData(mux.Vars(r)["zone-id"])
	if err := sendJSON(w, zone.Name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) getHeartbeat(w http.ResponseWriter, r *http.Request) {
	//zone := s.controller.GetZoneData(mux.Vars(r)["zone-id"])
	//if err := sendJSON(w, zone.Heartbeat); err != nil {
	if err := sendJSON(w, 0 != rand.Intn(15)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) getCountdown(w http.ResponseWriter, r *http.Request) {
	zone := s.controller.GetZoneData(mux.Vars(r)["zone-id"])
	if err := sendJSON(w, zone.TimeLeftSec); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) getCamerasID(w http.ResponseWriter, r *http.Request) {
	zone := s.controller.GetZoneData(mux.Vars(r)["zone-id"])
	cameraIDs := make([]string, 0, len(zone.Cameras))
	for cameraID := range zone.Cameras {
		cameraIDs = append(cameraIDs, cameraID)
	}

	if err := sendJSON(w, cameraIDs); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) getCamera(w http.ResponseWriter, r *http.Request) {
	zone := s.controller.GetZoneData(mux.Vars(r)["zone-id"])
	if err := sendJSON(w, zone.Cameras[mux.Vars(r)["camera-id"]]); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) buttonPress(w http.ResponseWriter, r *http.Request) {
	zone := mux.Vars(r)["zone-id"]
	if err := s.controller.MakeAction(zone); err != nil {
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

	s.router.HandleFunc("/index.js", home.Scripts()).Methods(http.MethodGet)
	s.router.HandleFunc("/style.css", home.Styles()).Methods(http.MethodGet)
	s.router.HandleFunc("/favicon.ico", home.Favicon()).Methods(http.MethodGet)
	s.router.HandleFunc("/dseg7.woff2", home.DSEG7()).Methods(http.MethodGet)

	wh := webhooks.NewHandler(s.config)
	s.router.HandleFunc("/webhooks", wh.WebHooksHandler).Methods(http.MethodPost)

	zone := s.router.PathPrefix("/zones/{zone-id}").Subrouter()
	zone.HandleFunc("", home.Handler()).Methods(http.MethodGet)
	zone.HandleFunc("/zone-name", s.getZoneName).Methods(http.MethodGet)
	zone.HandleFunc("/heartbeat", s.getHeartbeat).Methods(http.MethodGet)
	zone.HandleFunc("/countdown", s.getCountdown).Methods(http.MethodGet)
	zone.HandleFunc("/cameras-id", s.getCamerasID).Methods(http.MethodGet)
	zone.HandleFunc("/cameras/{camera-id}", s.getCamera).Methods(http.MethodGet)
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
