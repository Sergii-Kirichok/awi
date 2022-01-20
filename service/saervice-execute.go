package service

import (
	"fmt"
	"golang.org/x/sys/windows/svc"
	"strings"
	"time"
)

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

	go startServer(m)

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
				m.Info(1, testOutput)
				break loop
			//case svc.Pause:
			//	changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			//	//tick = slowtick
			//case svc.Continue:
			//	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			//	//tick = fasttick
			default:
				m.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}
