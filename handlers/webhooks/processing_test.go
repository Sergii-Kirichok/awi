package webhooks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var b = []byte(`
"event":{
"targetIdsDeprecated":[],
"classifiedObjects":[{
"subclass":"PERSON_BODY",
"objectId":32774759
}],
"analyticEventName":"человек на дороге",
"area":"","activity":
"OBJECT_PRESENT",
"eventTriggerTime":"2022-01-17T08:20:47.483Z",
"cameraId":"4xIx1DMwMLSwMDW2tDBKNNBLTsw1MBASCDilIfJR0W3apqrIovO_tncAAA",
"thisId":"32544",
"linkedEventId":"-1",
"timestamp":"2022-01-17T08:20:47.181Z",
"originatingEventId":"32544",
"originatingServerName":"ARCHER",
"location":"Unknown",
"type":"DEVICE_ANALYTICS_START",
"originatingServerId":"UMooFPEhRpayellyz_q3iA",
"cameraIds":[],
"entityIds":[]
}`)

func TestHandlerData_personAndCarAnalyticStart(t *testing.T) {
	testCases := []struct {
		name     string
		input    *Event
		hasError bool
	}{
		{
			name: "PERSON Start Event - OK",
			input: &Event{
				CameraId: "4xIx1DMwMLSwMDW2tDBKNNBLTsw1MBASCDilIfJR0W3apqrIovO_tncAAA",
				Type:     DEVICE_ANALYTICS_START,
				Activity: OBJECT_PRESENT,
				ClassifiedObjects: []classifiedObject{
					{
						Subclass: PERSON_BODY,
						ObjectId: 32774759,
					},
				},
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHandler(nil)
			err := h.personAndCarAnalyticStart(tc.input)
			if tc.hasError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
