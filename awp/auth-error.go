package awp

// Используем для отображения ошибке в вебе
func (a *Auth) LoginSetError(err error) {
	a.Lock()
	a.err = err
	a.Unlock()
}

func (a *Auth) GetError() error {
	//a.Lock()
	err := a.err
	//a.Unlock()
	return err
}
