package awp

func (a *Auth) Lock() {
	a.mu.Lock()
}

func (a *Auth) Unlock() {
	a.mu.Unlock()
}
