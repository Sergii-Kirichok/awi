package webhooks

import (
	"time"
)

// Всё сообщение
type WebhookMessage struct {
	Core
	Notifications       []Notification `json:"notifications"`
	AuthenticationToken string         `json:"authenticationToken"` // 3131746f6b656e3131537472696e67252164284d495353494e4729
}

// Основная часть. HELLO и HEARTBEAT состоят только из неё
type Core struct {
	SiteId string    `json:"siteId"` // IN2ir_lQRli_PuW2Un48ZQ
	Type   string    `json:"type"`   // HELLO, NOTIFICATION, HEARTBEAT
	Time   time.Time `json:"time"`   // 2022-01-12T19:07:38.725Z
}

type Notification struct {
	Core
	Id    string `json:"id"` // H8Ay08LvRM2Ie5hiWrSYRA
	Event *Event `json:"event"`
}

type Event struct {
	CameraId              string              `json:"cameraId"`      // 4xIx1DMwMLSwMDW2tDBKNNBLTsw1MBASCDilIfJR0W3apqrIovO_tncAAA
	ThisId                string              `json:"thisId"`        // 18036
	LinkedEventId         string              `json:"linkedEventId"` // -1
	Timestamp             time.Time           `json:"timestamp"`     // 2022-01-12T19:11:51.183Z
	RecordTriggerParams   recordTriggerParams `json:"recordTriggerParams"`
	OriginatingEventId    string              `json:"originatingEventId"`    // 18036
	OriginatingServerName string              `json:"originatingServerName"` // ARCHER
	Location              string              `json:"location"`              // Unknown
	Type                  string              `json:"type"`                  // DEVICE_MOTION_ST[ART|OP]/DEVICE_ANALYTICS_ST[ART|OP],DEVICE_DIGITAL_INPUT_O[N|FF] ....
	OriginatingServerId   string              `json:"originatingServerId"`   // UMooFPEhRpayellyz_q3iA
	CameraIds             []string            `json:"cameraIds"`             //
	EntityIds             []string            `json:"entityIds"`             //
	EntityId              string              `json:"entityId"`              // Для событий входов так точно нужен
	TargetIdsDeprecated   []interface{}       `json:"targetIdsDeprecated"`
	AnalyticEventName     string              `json:"analyticEventName"` // [Automation Analytic Event]
	Area                  string              `json:"area"`              // ""
	Activity              string              `json:"activity"`          // UNKNOWN, OBJECT_COUNTING_ENTER
	EventTriggerTime      time.Time           `json:"eventTriggerTime"`  //2020-05-20T19:55:55.550Z
	classifiedObjects     []classifiedObject  `json:"classifiedObjects"`
}
type classifiedObject struct {
	Subclass string `json:"subclass"` //VEHICLE_CAR, PERSON_FACE
	ObjectId int    `json:"objectId"` //242288522
}
type recordTriggerParams struct {
	PrePostRecordTime string `json:"prePostRecordTime"` // 5000000000
	IsPrePost         bool   `json:"isPrePost"`         // true
	TriggerJiffy      int    `json:"triggerJiffy"`      // 29079210
	IsTriggerJiffy    bool   `json:"isTriggerJiffy"`    // true

}

type Activity string

const (
	UNKNOWN                           Activity = "UNKNOWN"
	OBJECT_PRESENT                    Activity = "OBJECT_PRESENT" // Objects in area
	OBJECT_NOT_PRESENT                Activity = "OBJECT_NOT_PRESENT"
	PROHIBITED_DIRECTION              Activity = "PROHIBITED_DIRECTION"              // Direction violated. The event is triggered for each object that moves within 22 degrees of the prohibited direction for longer than the threshold time. One event is activated for each classified object that moves in the prohibited direction.
	OBJECT_LOITERING                  Activity = "OBJECT_LOITERING"                  // The event is triggered for each object that stays within the region of interest longer than the threshold time. Each object triggers a separate event.  The event resets when the object leaves the region of interest or the event times out.
	OBJECT_STOPPED                    Activity = "OBJECT_STOPPED"                    // The event is triggered if a classified object is detected moving within the region of interest then stops moving for longer than the threshold time. One event is activated for each object that stops. An object can only be tracked for up to 15 minutes.
	OBJECT_ENTERS                     Activity = "OBJECT_ENTERS"                     // Objects enter area
	OBJECT_LEAVES                     Activity = "OBJECT_LEAVES"                     // Objects leave area
	OBJECTS_CROSSED_BEAM              Activity = "OBJECTS_CROSSED_BEAM"              // The event is triggered when the specified number of objects have crossed the beam in the specified direction within the threshold time. If the number of objects is 1, the event is triggered after the threshold time elapses.
	OBJECT_COUNTING_ENTER             Activity = "OBJECT_COUNTING_ENTER"             // The event is triggered for each object that enters an occupancy area. One event is activated for each classified object that enters an area in the specified direction. o define an occupancy area, create Enter and Exit occupancy area events. To ensure accurate counts, be sure to create events for each camera with an entrance to the occupancy area.
	OBJECT_COUNTING_EXIT              Activity = "OBJECT_COUNTING_EXIT"              // Exit occupancy area
	OBJECTS_TOO_CLOSE                 Activity = "OBJECTS_TOO_CLOSE"                 // The event is triggered when two detected people are closer than the specified distance for longer than the threshold time. If there is a group of people, an event will be triggered for each pair that is too close.
	PERSON_ABNORMALLY_LOW_TEMPERATURE Activity = "PERSON_ABNORMALLY_LOW_TEMPERATURE" // This event is triggered when a lower temperature at or below the threshold is detected by a camera. Default is the calibrated value from the thermal camera. For example, if the default is 35.0 °C and 34.9 °C is detected, an event triggers the camera to record.
	PERSON_ELEVATED_TEMPERATURE       Activity = "PERSON_ELEVATED_TEMPERATURE"       // This event is triggered when an elevated temperature at or above the threshold is detected by a camera. Default is the calibrated value from the thermal camera. For example, if the default is 37.5 °C and 37.5 °C is detected, an event triggers the camera to record.
	PERSON_NORMAL_TEMPERATURE         Activity = "PERSON_NORMAL_TEMPERATURE"         // This event is triggered when a temperature within the acceptable range is detected by a camera. Using the same examples above, if 37.5 °C is the Object with elevated temperature and 35.0 °C is the Object with lower temperature, and a temperature is detected within the 35 °C to 37.5° C range, an event triggers the camera to record.
)

//type Object string
//
//const (
//	UNKNOWN Object = "UNKNOWN"
//	PERSON  Object = "PERSON"
//	VEHICLE Object = "VEHICLE"
//)

//The list of objects that triggered the event. Each classified object specifies a supertype (e.g. vehicle) and subtype (e.g. truck) and has a unique ID associated with it.
type ClassifiedObject string

const (
	PERSON            ClassifiedObject = "PERSON"
	PERSON_BODY       ClassifiedObject = "PERSON_BODY"
	PERSON_FACE       ClassifiedObject = "PERSON_FACE"
	VEHICLE           ClassifiedObject = "VEHICLE"
	VEHICLE_BICYCLE   ClassifiedObject = "VEHICLE_BICYCLE"
	VEHICLE_MOTORCYLE ClassifiedObject = "VEHICLE_MOTORCYLE"
	VEHICLE_CAR       ClassifiedObject = "VEHICLE_CAR"
	VEHICLE_TRUCK     ClassifiedObject = "VEHICLE_TRUCK"
	VEHICLE_BUS       ClassifiedObject = "VEHICLE_BUS"
)

type EventTypes string

const (
	DEVICE_DIGITAL_INPUT_ON  EventTypes = "DEVICE_DIGITAL_INPUT_ON"
	DEVICE_DIGITAL_INPUT_OFF EventTypes = "DEVICE_DIGITAL_INPUT_OFF"
	DEVICE_ANALYTICS_START   EventTypes = "DEVICE_ANALYTICS_START"
	DEVICE_ANALYTICS_STOP    EventTypes = "DEVICE_ANALYTICS_STOP"
)
