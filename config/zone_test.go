package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZoneIsOk(t *testing.T) {
	testCases := []struct {
		name     string
		zone     Zone
		expected bool
	}{
		{
			name: "cam1: CONNECTED, carOnAnyCam: true, ignore_car_state: true, Input: true, Person: True ",
			zone: Zone{
				Cameras: []Cam{
					{
						Id:       "1",
						ConState: "CONNECTED",
						Inputs: map[string]*Input{
							"123": {
								EntityId: "123",
								State:    StateTrue,
							},
						},
						Car:            "",
						CarEventId:     "",
						Person:         StateTrue,
						PersonEventId:  "",
						InputsDisabled: false,
					},
					{
						Id:       "2",
						ConState: "CONNECTED",
						Inputs: map[string]*Input{
							"123": {
								EntityId: "124",
								State:    StateTrue,
							},
						},
						Car:            StateTrue,
						CarEventId:     "",
						Person:         StateTrue,
						PersonEventId:  "",
						InputsDisabled: false,
					},
				},
				CarOnAnyCamera: true,
				IgnoreCarState: true,
			},
			expected: true,
		},
		{
			name: "cam1: CONNECTED, carOnAnyCam: false, ignore_car_state: true, Input: true, Person: True ",
			zone: Zone{
				Cameras: []Cam{
					{
						Id:       "1",
						ConState: "CONNECTED",
						Inputs: map[string]*Input{
							"123": {
								EntityId: "123",
								State:    StateTrue,
							},
						},
						Car:            "",
						CarEventId:     "",
						Person:         StateTrue,
						PersonEventId:  "",
						InputsDisabled: false,
					},
					{
						Id:       "2",
						ConState: "CONNECTED",
						Inputs: map[string]*Input{
							"123": {
								EntityId: "124",
								State:    StateTrue,
							},
						},
						Car:            StateTrue,
						CarEventId:     "",
						Person:         StateTrue,
						PersonEventId:  "",
						InputsDisabled: false,
					},
				},
				CarOnAnyCamera: false,
				IgnoreCarState: true,
			},
			expected: true,
		},
		{
			name: "cam1: CONNECTED, carOnAnyCam: true, ignore_car_state: false, Input: true, Person: True ",
			zone: Zone{
				Cameras: []Cam{
					{
						Id:       "1",
						ConState: "CONNECTED",
						Inputs: map[string]*Input{
							"123": {
								EntityId: "123",
								State:    StateTrue,
							},
						},
						Car:            "",
						CarEventId:     "",
						Person:         StateTrue,
						PersonEventId:  "",
						InputsDisabled: false,
					},
					{
						Id:       "2",
						ConState: "CONNECTED",
						Inputs: map[string]*Input{
							"123": {
								EntityId: "124",
								State:    StateTrue,
							},
						},
						Car:            StateTrue,
						CarEventId:     "",
						Person:         StateTrue,
						PersonEventId:  "",
						InputsDisabled: false,
					},
				},
				CarOnAnyCamera: true,
				IgnoreCarState: false,
			},
			expected: true,
		},
		{
			name: "cam1: CONNECTED, carOnAnyCam: false, ignore_car_state: fasle, Input: true, Person: True ",
			zone: Zone{
				Cameras: []Cam{
					{
						Id:       "1",
						ConState: "CONNECTED",
						Inputs: map[string]*Input{
							"123": {
								EntityId: "123",
								State:    StateTrue,
							},
						},
						Car:            "",
						CarEventId:     "",
						Person:         StateTrue,
						PersonEventId:  "",
						InputsDisabled: false,
					},
					{
						Id:       "2",
						ConState: "CONNECTED",
						Inputs: map[string]*Input{
							"123": {
								EntityId: "124",
								State:    StateTrue,
							},
						},
						Car:            StateTrue,
						CarEventId:     "",
						Person:         StateTrue,
						PersonEventId:  "",
						InputsDisabled: false,
					},
				},
				CarOnAnyCamera: false,
				IgnoreCarState: false,
			},
			expected: false,
		},
		{
			name: "cam1: CONNECTED, carOnAnyCam: true, ignore_car_state: true, Input: false, Person: True ",
			zone: Zone{
				Cameras: []Cam{
					{
						Id:       "1",
						ConState: "CONNECTED",
						Inputs: map[string]*Input{
							"123": {
								EntityId: "123",
								State:    StateFalse,
							},
						},
						Car:            "",
						CarEventId:     "",
						Person:         StateTrue,
						PersonEventId:  "",
						InputsDisabled: false,
					},
					{
						Id:       "2",
						ConState: "CONNECTED",
						Inputs: map[string]*Input{
							"123": {
								EntityId: "124",
								State:    StateTrue,
							},
						},
						Car:            StateTrue,
						CarEventId:     "",
						Person:         StateTrue,
						PersonEventId:  "",
						InputsDisabled: false,
					},
				},
				CarOnAnyCamera: true,
				IgnoreCarState: true,
			},
			expected: false,
		},
		{
			name: "cam1: CONNECTED, carOnAnyCam: true, ignore_car_state: true, Input: unknown, Person: True ",
			zone: Zone{
				Cameras: []Cam{
					{
						Id:       "1",
						ConState: "CONNECTED",
						Inputs: map[string]*Input{
							"123": {
								EntityId: "123",
								State:    StateUnknown,
							},
						},
						Car:            "",
						CarEventId:     "",
						Person:         StateTrue,
						PersonEventId:  "",
						InputsDisabled: false,
					},
					{
						Id:       "2",
						ConState: "CONNECTED",
						Inputs: map[string]*Input{
							"123": {
								EntityId: "124",
								State:    StateTrue,
							},
						},
						Car:            StateTrue,
						CarEventId:     "",
						Person:         StateTrue,
						PersonEventId:  "",
						InputsDisabled: false,
					},
				},
				CarOnAnyCamera: true,
				IgnoreCarState: true,
			},
			expected: false,
		},
		{
			name: "cam1: CONNECTED, carOnAnyCam: true, ignore_car_state: true, Input: true, Person: False ",
			zone: Zone{
				Cameras: []Cam{
					{
						Id:       "1",
						ConState: "CONNECTED",
						Inputs: map[string]*Input{
							"123": {
								EntityId: "123",
								State:    StateTrue,
							},
						},
						Car:            "",
						CarEventId:     "",
						Person:         StateFalse,
						PersonEventId:  "",
						InputsDisabled: false,
					},
					{
						Id:       "2",
						ConState: "CONNECTED",
						Inputs: map[string]*Input{
							"123": {
								EntityId: "124",
								State:    StateTrue,
							},
						},
						Car:            StateTrue,
						CarEventId:     "",
						Person:         StateTrue,
						PersonEventId:  "",
						InputsDisabled: false,
					},
				},
				CarOnAnyCamera: true,
				IgnoreCarState: true,
			},
			expected: false,
		},
		{
			name: "cam1: CONNECTED, carOnAnyCam: true, ignore_car_state: true, Input: true, Person: unknown ",
			zone: Zone{
				Cameras: []Cam{
					{
						Id:       "1",
						ConState: "CONNECTED",
						Inputs: map[string]*Input{
							"123": {
								EntityId: "123",
								State:    StateTrue,
							},
						},
						Car:            "",
						CarEventId:     "",
						Person:         StateUnknown,
						PersonEventId:  "",
						InputsDisabled: false,
					},
					{
						Id:       "2",
						ConState: "CONNECTED",
						Inputs: map[string]*Input{
							"123": {
								EntityId: "124",
								State:    StateTrue,
							},
						},
						Car:            StateTrue,
						CarEventId:     "",
						Person:         StateTrue,
						PersonEventId:  "",
						InputsDisabled: false,
					},
				},
				CarOnAnyCamera: true,
				IgnoreCarState: true,
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ZoneIsOk(&tc.zone)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
