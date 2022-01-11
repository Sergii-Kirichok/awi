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
	info := version.GetInfo()
	//rootDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	//cfg := config.New().Load(rootDir)

	//on("./")
	//u := rest.UserGetFromRESTByAddress(&cfg, "U3")
	//fmt.Println(u)

	//val ,_ :=time.Parse("2006-01-02T15:04:05",u.Users[0].CreationDate)
	//val.Format("2006-01-02 15:04:05")
	//fmt.Printf("DateTime %s, dateTimeFormated %s",val, val.Format("2006-01-02 15:04:05") )
	//fmt.Println(coder.Coder("1"))
	//return

	m := service.Myservice{info.Name, info.Version}
	isIntSess, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("Не удалось определить работаем-ли мы в интерактивном сеансе: %v", err)
	}
	if !isIntSess {
		service.RunService(info.SvcName, false, &m)
		return
	}

	if len(os.Args) < 2 {
		usage("Комманда не определена.")
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "debug":
		service.RunService(info.SvcName, true, &m)
		return
	case "install":
		err = service.InstallService(info.SvcName, fmt.Sprintf("%s [%s]", info.Name, info.Version))
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
