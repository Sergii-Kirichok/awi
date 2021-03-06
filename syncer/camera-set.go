package syncer

import (
	"awi/awp"
	"awi/config"
)

// TODO: Переделать на методы конфига. Синхронизатор не должен сам лезть в изменение данных, должен только дёрнуть метод и передать данные для заполнения
// Но, возможно мы столкнёмся с циклическим вызовом awi и config, поэтому надо будет передавать структуру камеры из конфига
func (s *syncer) cameraSet(camera awp.Camera) error {
	s.auth.Config.Lock()
	defer s.auth.Config.Unlock()

	for zIndex, zData := range s.auth.Config.Zones {
		for camIndex, camData := range zData.Cameras {
			// Нашли камеру и обновляем данные по ней
			if camera.Serial == camData.Serial {
				s.auth.Config.Zones[zIndex].Cameras[camIndex].Id = camera.Id
				s.auth.Config.Zones[zIndex].Cameras[camIndex].Name = camera.Name

				camState := config.CamState(camera.ConnectionState)
				s.auth.Config.Zones[zIndex].Cameras[camIndex].ConState = camState

				// Если камера отключена, то сбрасываем её состояния
				if camState != config.CamConnected {
					s.auth.Config.ClearCamStates(camera.Id)
				}
				// Устанавливаем/Обновляем данные о входах
				s.inputsSetUpdate(zIndex, camIndex, &camera)
				return nil
			}
		}
	}
	return nil
}
