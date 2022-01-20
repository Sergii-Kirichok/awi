package service

import (
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

func RunService(name string, isDebug bool, m *Myservice) {
	m.Info(1, fmt.Sprintf("starting %s service", name))

	run := svc.Run
	if isDebug {
		run = debug.Run
	}

	if err := run(name, m); err != nil {
		m.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}

	// Подчищаем за собой.
	if err := m.garbageCleaner(); err != nil {
		m.Error(123, fmt.Sprintf("garmabeCleaner: %s", err))
	}

	m.Info(1, fmt.Sprintf("%s service stopped", name))
	m.Log.Close()
}

func (m *Myservice) garbageCleaner() error {
	// Выходим из веб-поинта если сессия была создана
	err := m.authLogout()
	return err
}

func (m *Myservice) authLogout() error {
	m.Auth.Lock()
	err := m.Auth.Logout()
	m.Auth.Unlock()
	return err
}
