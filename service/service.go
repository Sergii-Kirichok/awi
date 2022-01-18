package service

import (
	"awi/awp"
	"awi/config"
	"awi/controller"
	"awi/syncer"
	"awi/webserver"
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"strings"
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
	msg := fmt.Sprintf("Запуск сервера %s", s)
	elog.Info(1, msg)

start:
	cfg, err := config.New().Load()
	if err != nil {
		elog.Error(1, fmt.Sprintf("Load config: %s", err))
		time.Sleep(10 * time.Second)
		goto start
	}

	//Если не смогли авторизоваться при старте, то явно какая-то ошибка в конфиге или с самим сервером. Ругаемся, но идём дальше (нужен веб для отображения информации)
	auth, err := awp.NewAuth(cfg).Login()
	if err != nil {
		elog.Error(2, fmt.Sprintf("%s", err))
	}

	// Синхронизатор. Проверяет конфиг, находит Ид камер, создаёт и удаляет вебхуки, ...
	go syncer.New(auth).Sync()

	// Cтруктура хранящая данные используемые для генерации данных по зонам, камерам и ошибок (связи,синхронизации) передаваемых при get-запросе из веба
	control := controller.New(auth)

	// Сервис отвечающий за: таймеры обратного отсчёта по зонам и их состояния. Веб работает с ней и  методами controllera\
	go control.Service()

	// WebServer, принимает и обрабатываем webhook-и от WebPointa, так-же отдаёт страничку с Кнопкой, таймером обратного отсчёта, значками состояния, ...
	webserver.New(s.Name, s.Version, cfg, control).ListenAndServeHTTPS()
}

func (m *Myservice) String() string {
	return fmt.Sprintf("%s [%s]", m.Name, m.Version)
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
