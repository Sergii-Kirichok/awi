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
	"sync"
	"sync/atomic"
	"time"
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

type Input struct {
	Id    string `json:"id"`
	State bool   `json:"state"`
}

type CameraStates struct {
	Cars   bool              `json:"cars"`
	Humans bool              `json:"humans"`
	Inputs map[string]*Input `json:"inputs"`
}

//const stabilizationTime = 5 * time.Minute / time.Second
const stabilizationTime = (15 * time.Second) / time.Second

var (
	timeLeft = int32(stabilizationTime)

	mutex sync.RWMutex
)

func (s *Server) getCountdown(w http.ResponseWriter, r *http.Request) {
	zone := s.controller.GetZoneData(mux.Vars(r)["zone-id"])

	if !isStartCountdown(zone) {
		atomic.StoreInt32(&timeLeft, int32(stabilizationTime))
	} else if atomic.LoadInt32(&timeLeft) > 0 {
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

func isStartCountdown(zone controller.Zone) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	for _, camera := range zone.Cameras {
		for _, input := range camera.Inputs {
			if !input.State {
				return false
			}
		}

		if !camera.Car || !camera.Human {
			return false
		}
	}

	return true
}

func (s *Server) getCamerasIDs(w http.ResponseWriter, r *http.Request) {
	zone := s.controller.GetZoneData(mux.Vars(r)["zone-id"])

	mutex.RLock()
	cameraIDs := make([]string, 0, len(zone.Cameras))
	for cameraID, c := range zone.Cameras {
		fmt.Printf("cameraID: %s, CAM: %v\n", cameraID, c)
		cameraIDs = append(cameraIDs, cameraID)
	}
	mutex.RUnlock()

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(&cameraIDs); err != nil {
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

func (s *Server) getCameraInfo(w http.ResponseWriter, r *http.Request) {
	zone := s.controller.GetZoneData(mux.Vars(r)["zone-id"])

	mutex.RLock()
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(zone.Cameras[mux.Vars(r)["camera-id"]]); err != nil {
		log.Printf("encoding error with data <%s> : %s\n", b.String(), err)
		w.WriteHeader(http.StatusInternalServerError)
		mutex.RUnlock()
		return
	}
	mutex.RUnlock()

	if _, err := w.Write(b.Bytes()); err != nil {
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

	s.router.HandleFunc("/index.js", home.Scripts()).Methods(http.MethodGet)
	s.router.HandleFunc("/style.css", home.Styles()).Methods(http.MethodGet)
	s.router.HandleFunc("/favicon.ico", home.Favicon()).Methods(http.MethodGet)
	s.router.HandleFunc("/dseg7.woff2", home.DSEG7()).Methods(http.MethodGet)

	wh := webhooks.NewHandler(s.config)
	s.router.HandleFunc("/webhooks", wh.WebHooksHandler).Methods(http.MethodPost)

	zone := s.router.PathPrefix("/zones/{zone-id}").Subrouter()
	zone.HandleFunc("", home.Handler()).Methods(http.MethodGet)
	zone.HandleFunc("/countdown", s.getCountdown).Methods(http.MethodGet)
	zone.HandleFunc("/cameras-ids", s.getCamerasIDs).Methods(http.MethodGet)
	zone.HandleFunc("/cameras-info/{camera-id}", s.getCameraInfo).Methods(http.MethodGet)

	//go updateCamerasStates()
	// Запуск веб-сервера
	rootDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	certFile := fmt.Sprintf("%s/%s", rootDir, s.config.WWWCertificate)
	keyFile := fmt.Sprintf("%s/%s", rootDir, s.config.WWWCertificateKey)
	if err := http.ListenAndServeTLS(bind, certFile, keyFile, s.router); err != nil {
		//if err := http.ListenAndServeTLS(s.Bind, s.Certificate, s.CertificateKey, nil); err != nil {
		log.Fatal("HTTPS-Err: ", err)
	}
}
