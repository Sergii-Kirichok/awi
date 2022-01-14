package syncer

import (
	"awi/awp"
)

func (s *syncer) cameraSet(camera awp.Camera) error {
	s.auth.Config.Lock()
	defer s.auth.Config.Unlock()
	for zIndex, zData := range s.auth.Config.Zones {
		for camIndex, camData := range zData.Cameras {
			// Нашли камеру и обновляем данные по ней
			if camera.Serial == camData.Serial && camData.Id == "" {
				s.auth.Config.Zones[zIndex].Cameras[camIndex].Id = camera.Id
				s.auth.Config.Zones[zIndex].Cameras[camIndex].Name = camera.Name
				// Устанавливаем/Обновляем данные о входах
				s.inputsSetUpdate(zIndex, camIndex, &camera)
				return nil
			}
		}
	}
	return nil
}
