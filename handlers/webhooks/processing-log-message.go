package webhooks

import "log"

func (h HandlerData) ProcessingLogMessage(e *Event) {
	cam, err := h.cfg.GetCameraById(e.CameraId)
	if err != nil {
		cam.Name = e.CameraId
	}
	log.Printf("[%s] Камера: \"%s\", Cобытие: \"%s\"\n", e.Type, cam.Name, e.AnalyticEventName)
}

func (h HandlerData) ProcessingInputLogMessage(e *Event) {
	for _, cId := range e.CameraIds {
		cam, err := h.cfg.GetCameraById(cId)
		if err != nil {
			cam.Name = e.CameraId
		}
		log.Printf("[%s] Камера: \"%s\"\n", e.Type, cam.Name)
	}
}
