package service

import (
	"awi/awp"
	"awi/config"
	"awi/version"
	"fmt"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"log"
)

var (
	cfg *config.Config
)

type Myservice struct {
	Name    string
	Version string
	Auth    *awp.Auth // Надо дернуть logout при выходе.
	debug.Log
}

func New(info *version.Info, isDebug string) *Myservice {
	var err error
	s := &Myservice{
		Name:    info.Name,
		Version: info.Version,
	}

	if isDebug == "debug" {
		s.Log = debug.New(info.SvcName)
	} else {
		s.Log, err = eventlog.Open(info.SvcName)
		if err != nil {
			log.Fatalf("Can't open eventlog: %s", err)
		}
	}

	return s
}

func (m *Myservice) String() string {
	return fmt.Sprintf("%s [%s]", m.Name, m.Version)
}
