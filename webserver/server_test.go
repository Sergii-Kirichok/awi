package webserver

import (
	"awi/controller"
	"testing"
)

func TestIsStartCountdown(t *testing.T) {
	testCases := []struct {
		name   string
		zone   controller.Zone
		output bool
	}{
		{
			zone: controller.Zone{
				Cameras: map[string]*controller.Camera{
					"first": {
						Car:   true,
						Human: true,
						Inputs: map[string]*controller.Input{
							"first": {
								State: true,
							},
						},
					},
					"second": {
						Car:   true,
						Human: true,
						Inputs: map[string]*controller.Input{
							"first": {
								State: true,
							},
							"second": {
								State: true,
							},
						},
					},
				},
			},
			output: true,
		},
		{
			zone: controller.Zone{
				Cameras: map[string]*controller.Camera{
					"first": {
						Car:   true,
						Human: true,
						Inputs: map[string]*controller.Input{
							"first": {
								State: true,
							},
						}},
					"second": {
						Car:    true,
						Human:  true,
						Inputs: map[string]*controller.Input{},
					},
				},
			},
			output: true,
		},
		{
			zone: controller.Zone{
				Cameras: map[string]*controller.Camera{
					"first": {
						Car:    true,
						Human:  false,
						Inputs: map[string]*controller.Input{},
					},
					"second": {
						Car:    true,
						Human:  true,
						Inputs: map[string]*controller.Input{},
					},
				},
			},
			output: false,
		},
		{
			zone: controller.Zone{
				Cameras: map[string]*controller.Camera{
					"first": {
						Car:   false,
						Human: true,
						Inputs: map[string]*controller.Input{
							"first": {
								State: false,
							},
						},
					},
					"second": {
						Car:   true,
						Human: true,
						Inputs: map[string]*controller.Input{
							"second": {
								State: false,
							},
						},
					},
				},
			},
			output: false,
		},
		{
			zone: controller.Zone{
				Cameras: map[string]*controller.Camera{
					"first": {
						Car:   false,
						Human: false,
						Inputs: map[string]*controller.Input{
							"first": {
								State: false,
							},
						},
					},
					"second": {
						Car:   false,
						Human: false,
						Inputs: map[string]*controller.Input{
							"second": {
								State: false,
							},
						},
					},
				},
			},
			output: false,
		},
	}

	mutex.RLock()
	defer mutex.RUnlock()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := isStartCountdown(tc.zone); tc.output != actual {
				t.Fail()
			}
		})
	}
}
