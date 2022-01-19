package webhooks

import "fmt"

//{
//	"siteId":"IN2ir_lQRli_PuW2Un48ZQ",
//	"type":"HELLO",
//	"time":"2022-01-12T13:46:19.981Z",
//	"authenticationToken":"3333746f6b656e3333537472696e67252164284d495353494e4729"
//}
func (h *HandlerData) processingHello(t MainType) error {
	fmt.Printf("Processing: [%s] \n", t)
	h.controller.UpdateHeartBeat()
	return nil
}
