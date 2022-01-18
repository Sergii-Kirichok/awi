package main

import (
	"awi/service"
	"awi/version"
	"fmt"
	"golang.org/x/sys/windows/svc"
	"log"
	"os"
	"strings"
)

func main() {
	info := version.NewInfo()
	m := service.Myservice{Name: info.Name, Version: info.Version}

	isIntSess, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Fatalf("Не удалось определить работаем-ли мы в интерактивном сеансе: %v", err)
	}

	if !isIntSess {
		service.RunService(info.SvcName, false, &m)
		return
	}

	if len(os.Args) < 2 {
		usage("Команда не определена.")
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "debug":
		service.RunService(info.SvcName, true, &m)
		return
	case "install":
		err = service.InstallService(info.SvcName, fmt.Sprintf("%s", info))
	case "remove":
		err = service.RemoveService(info.SvcName)
	case "start":
		err = service.StartService(info.SvcName)
	case "stop":
		err = service.ControlService(info.SvcName, svc.Stop, svc.Stopped)
	case "pause":
		err = service.ControlService(info.SvcName, svc.Pause, svc.Paused)
	case "continue":
		err = service.ControlService(info.SvcName, svc.Continue, svc.Running)
	default:
		usage(fmt.Sprintf("Комманда %s не поддерживается или написана не верно.", cmd))
	}
	if err != nil {
		log.Fatalf("Не удалось выполнить %s %s: %v", cmd, info.SvcName, err)
	}
	return
}

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"Используем так: %s <комманда>\n"+
			"       где <комманда> должна быть одна из ниже перечиленных\n"+
			"       install, remove, start, stop.\n",
		errmsg, os.Args[0])
	os.Exit(2)
}
