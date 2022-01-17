package main

import (
	"awi/handlers/webhooks"
	"fmt"
)

func main() {
	hookData := `{"siteId":"IN2ir_lQRli_PuW2Un48ZQ","type":"NOTIFICATION","time":"2022-01-17T09:35:55.447Z","notifications":[{"siteId":"IN2ir_lQRli_PuW2Un48ZQ","type":"EVENT","time":"2022-01-17T09:35:55.447Z","id":"nKYs4HUDREWNT02FGWPMtA","event":{"targetIdsDeprecated":[],"classifiedObjects":[{"subclass":"PERSON_BODY","objectId":32775346}],"analyticEventName":"человек на дороге","area":"","activity":"OBJECT_PRESENT","eventTriggerTime":"2022-01-17T09:35:55.429Z","cameraId":"4xIx1DMwMLSwMDW2tDBKNNBLTsw1MBASCDilIfJR0W3apqrIovO_tncAAA","thisId":"32732","linkedEventId":"-1","timestamp":"2022-01-17T09:35:55.048Z","originatingEventId":"32732","originatingServerName":"ARCHER","location":"Unknown","type":"DEVICE_ANALYTICS_START","originatingServerId":"UMooFPEhRpayellyz_q3iA","cameraIds":[],"entityIds":[]}}],"authenticationToken":"323032322d30312d31372031313a33343a31382e32373439313335202b3032303020454554206d3d2b362e343934333134313031"}`
	//h := webhooks.NewHandler(nil)
	msg := &webhooks.WebhookMessage{}
	msg.Parse([]byte(hookData))
	fmt.Printf("Parsed data: %#v\n", msg)
	for _, notification := range msg.Notifications {
		fmt.Printf("Event Data: %#v\n", notification.Event)
	}
}
