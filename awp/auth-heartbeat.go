package awp

import "time"

func (a *Auth) UpdateHeartBeat() error {
	a.Lock()
	a.LastHeartbeat = time.Now()
	a.Unlock()
	return nil
}

func (a *Auth) GetHeartBeat() bool {
	//a.Lock()
	var hbState bool
	if time.Since(a.LastHeartbeat).Milliseconds() <= (HeartBeatDelayMs * 2) {
		hbState = true
	}
	//a.Unlock()
	return hbState
}
