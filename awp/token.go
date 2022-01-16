package awp

import (
	"crypto/sha256"
	"fmt"
	"time"
)

//FO#09121901:1641806505:a455e137865a4a85866f1de723814131546095170f9b2c22d540a2e23e49c3f1
func (a *Auth) genToken() string {
	a.Config.Lock()
	defer a.Config.Unlock()

	tNow := time.Now()
	timeStamp := tNow.Unix()
	//timeStamp := 1641806505
	data := fmt.Sprintf("%d%s", timeStamp, a.Config.DevKey)
	hexEncoded := sha256.Sum256([]byte(data))
	//fmt.Printf("Nonce: %s, Key: %s, Time: %d, data: %s,Hash: %x\n", c.DevNonce, c.DevKey, timeStamp, data, hexEncoded)

	token := fmt.Sprintf("%s:%d:%x", a.Config.DevNonce, timeStamp, hexEncoded)
	return token
}
