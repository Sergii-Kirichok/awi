package syncer

import (
	"awi/awp"
	"awi/config"
)

// Обновляем данные о входах у выбранной камеры
// Мьютекс установлен уровнем выше, и мы можем спокойно писать и читать из/в структуру *Config
func (s *syncer) inputsSetUpdate(zIndex int, camIndex int, camSrc *awp.Camera) {
	camDest := &s.auth.Config.Zones[zIndex].Cameras[camIndex]
	// Для обратной проверки (удаления лишних входов если они там есть)
	srsInputs := make(map[string]awp.Link, len(camSrc.Links))

	// Такой мапы вообще ещё нет, программа только стартонула
	if camDest.Inputs == nil {
		camDest.Inputs = make(map[string]*config.Input)
	}

	//Создаём массив цифровых входов, остальное нас не интересует
	for _, v := range camSrc.Links {
		if v.Type == awp.DIGITAL_INPUT {
			srsInputs[v.Source] = v
			// Если такого входа нет - добавляем его
			if _, ok := camDest.Inputs[v.Source]; !ok {
				camDest.Inputs[v.Source] = &config.Input{EntityId: v.Source}
			}
		}
	}

	// Обратная проверка
	for k, v := range camDest.Inputs {
		// такого входа не существует, удаляем его
		if _, ok := srsInputs[v.EntityId]; !ok {
			delete(camDest.Inputs, k)
		}
	}
}
