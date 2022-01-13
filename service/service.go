package service

import (
	"awi/config"
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"os"
	"strings"
	"suv02-server/webserver"
	"time"
)

var (
	cfg  *config.Config
	elog debug.Log
)

type Myservice struct {
	Name    string
	Version string
}

func startServer(s Myservice) {
	cfg = config.New().Load()
	msg := fmt.Sprintf("Запуск сервера %s [%s]", s.Name, s.Version)

	//db.LogCreate(&cfg, db.System, "", "", msg)
	elog.Info(1, msg)

	//Тестовый пользователь. Создавался и удалялся во время написания и отладки программы
	//if cfg.Debug {
	//	db.ClearTestCardTmp(&cfg, "1A000000000000000242630B", "23")
	//	_, sd := rest.VisitorDelete(&cfg, "U23")
	//	log.Printf("пользователь U23 удален, если код  == 200, код [%s]\n", sd)
	//}

	//Запускаем сервисы синхронизации
	//go sync.Syncer(&cfg)

	//Пдоготавливаем конфиг для веба

	certificatePath := fmt.Sprintf("%v/%s", rootDir, cfg.WWWCertificate)
	certificateKeyPath := fmt.Sprintf("%v/%s", rootDir, cfg.WWWCertificateKey)
	if _, err := os.Stat(certificatePath); os.IsNotExist(err) {
		elog.Warning(1, fmt.Sprintf("Не верный путь к сертификату %v", certificatePath))
	}
	if _, err := os.Stat(certificateKeyPath); os.IsNotExist(err) {
		elog.Warning(1, fmt.Sprintf("Не верный путь к ключу сертификата %v", certificateKeyPath))
	}

	srv := webserver.NewServer()
	srv.Name = s.Name
	srv.Version = s.Version
	srv.Bind = cfg.WWWAddr
	srv.Certificate = certificatePath
	srv.CertificateKey = certificateKeyPath
	//srv.Conf = &cfg

	//Запуск веб сервера http
	//if cfg.WWWHTTPRedirect {
	//	go srv.ListenAndServeHTTP()
	//}
	//Запуск веб сервера https
	//srv.ListenAndServeHTTPS()
}

func (m *Myservice) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	//Какие комманды мы будем отрабатывать
	//const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	//Пишем ответ, что мы приняли комманду и начинаем её обрабатывать.
	changes <- svc.Status{State: svc.StartPending}

	//Это пишет данные в канал с выбранным интервалом
	//fasttick := time.Tick(5 * time.Second)
	//slowtick := time.Tick(20 * time.Second)
	//tick := fasttick

	//Отдаём текущий статус, и какие комманды мы можем отрабатывать
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	go startServer(*m)

loop:
	//Я так понимаю, вот тут мы крутимся пока сервис разботает или нет
	for {
		//выбираем любой ГОТОВЫЙ канал  и работаетм с ним
		select {
		//Тут мы получаем блокируемое чтение. т.е. когда по-истечении времени таймер напишет данные в канал
		//будет выполенно условие
		//case <-tick:
		//	beep()
		//	elog.Info(1, "beep123")

		//А вот тут мы читаем данные из канала которые нам посылает система
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				// golang.org/x/sys/windows/svc.TestExample is verifying this output.
				testOutput := strings.Join(args, "-")
				//testOutput += fmt.Sprintf("-%d", c.Context)
				elog.Info(1, testOutput)
				break loop
			//case svc.Pause:
			//	changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			//	//tick = slowtick
			//case svc.Continue:
			//	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			//	//tick = fasttick
			default:
				elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func RunService(name string, isDebug bool, m *Myservice) {
	var err error
	if isDebug {
		elog = debug.New(name)
	} else {
		elog, err = eventlog.Open(name)
		if err != nil {
			return
		}
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", name))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	err = run(name, m)
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", name))
}
