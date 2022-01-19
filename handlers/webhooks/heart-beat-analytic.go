package webhooks

import (
	"log"
)

//{
//	"siteId":"IN2ir_lQRli_PuW2Un48ZQ",
//	"type":"HEARTBEAT",
//	"time":"2022-01-12T16:42:31.349Z",
//	"authenticationToken":"3733746f6b656e3733537472696e67252164284d495353494e4729"
//}
func (h *HandlerData) processingHeartbeat(t MainType) error {
	log.Printf("Processing: [%s] \n", t)
	return h.controller.UpdateHeartBeat()
}
