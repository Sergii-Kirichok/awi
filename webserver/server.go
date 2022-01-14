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
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

type Server struct {
	name    string
	version string
	router  *mux.Router

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
		router:  mux.NewRouter(),
	}
}

//const stabilizationTime = 5 * time.Minute / time.Second
const stabilizationTime = (15 * time.Second) / time.Second

var timeLeft = int32(stabilizationTime)

var camerasStates = map[string]*CameraStates{
	"4xIx1DMwMLSwMDW2tDBKNNBLTsw1MBASCDilIfJR0W3apqrIovO_tncAAA": {
		Cars:   false,
		Humans: false,
		Inputs: false,
	},
}

func getCountdown(w http.ResponseWriter, r *http.Request) {
	if !isStartCountdown() {
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

func isStartCountdown() bool {
	for _, state := range camerasStates {
		if !state.Cars || !state.Humans || !state.Inputs {
			return false
		}
	}

	return true
}

func getCamerasIDs(w http.ResponseWriter, r *http.Request) {
	cameraIDs := []string{"4xIx1DMwMLSwMDW2tDBKNNBLTsw1MBASCDilIfJR0W3apqrIovO_tncAAA"}
	//cameraIDs := []string{"4xIx1DMw", "DW2tDBK", "CDilIfJ"}

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

type CameraStates struct {
	Cars   bool `json:"cars"`
	Humans bool `json:"humans"`
	Inputs bool `json:"inputs"`
}

func getCameraInfo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["camera-id"]

	states := map[string]*CameraStates{
		"4xIx1DMwMLSwMDW2tDBKNNBLTsw1MBASCDilIfJR0W3apqrIovO_tncAAA": {
			Cars:   0 == rand.Intn(2),
			Humans: 0 == rand.Intn(2),
			Inputs: 0 == rand.Intn(2),
		},
		//"4xIx1DMw": {
		//	Cars:   0 == rand.Intn(2),
		//	Humans: 0 == rand.Intn(2),
		//	Inputs: 0 == rand.Intn(2),
		//},
		//"DW2tDBK": {
		//	Cars:   0 == rand.Intn(2),
		//	Humans: 0 == rand.Intn(2),
		//	Inputs: 0 == rand.Intn(2),
		//},
		//"CDilIfJ": {
		//	Cars:   0 == rand.Intn(2),
		//	Humans: 0 == rand.Intn(2),
		//	Inputs: 0 == rand.Intn(2),
		//},
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(states[id]); err != nil {
		log.Printf("encoding error with data <%s> : %s\n", b.String(), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b.Bytes()); err != nil {
		log.Printf("response error with empty data: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func resetTimer(w http.ResponseWriter, r *http.Request) {
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

	s.router.HandleFunc("/", home.Handler()).Methods(http.MethodGet)
	s.router.HandleFunc("/index.js", home.Scripts()).Methods(http.MethodGet)
	s.router.HandleFunc("/style.css", home.Styles()).Methods(http.MethodGet)
	s.router.HandleFunc("/favicon.ico", home.Favicon()).Methods(http.MethodGet)
	s.router.HandleFunc("/dseg7.woff2", home.DSEG7()).Methods(http.MethodGet)

	wh := webhooks.NewHandler(s.config)
	s.router.HandleFunc("/webhooks", wh.WebHooksHandler).Methods(http.MethodPost)

	s.router.HandleFunc("/countdown", getCountdown).Methods(http.MethodGet)
	s.router.HandleFunc("/cameras-ids", getCamerasIDs).Methods(http.MethodGet)
	s.router.HandleFunc("/cameras-info/{camera-id}", getCameraInfo).Methods(http.MethodGet)
	s.router.HandleFunc("/reset-timer", resetTimer).Methods(http.MethodPost)

	// Запуск веб-сервера
	rootDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	certFile := fmt.Sprintf("%s/%s", rootDir, s.config.WWWCertificate)
	keyFile := fmt.Sprintf("%s/%s", rootDir, s.config.WWWCertificateKey)
	if err := http.ListenAndServeTLS(bind, certFile, keyFile, s.router); err != nil {
		//if err := http.ListenAndServeTLS(s.Bind, s.Certificate, s.CertificateKey, nil); err != nil {
		log.Fatal("HTTPS-Err: ", err)
	}
}
