package webserver

import (
	"testing"
)

func TestIsStartCountdown(t *testing.T) {
	testCases := []struct {
		name          string
		camerasStates map[string]CameraStates
		output        bool
	}{
		{
			camerasStates: map[string]CameraStates{
				"first": {
					Cars:   true,
					Humans: true,
					Inputs: true,
				},
				"second": {
					Cars:   true,
					Humans: true,
					Inputs: true,
				},
			},
			output: true,
		},
		{
			camerasStates: map[string]CameraStates{
				"first": {
					Cars:   true,
					Humans: true,
					Inputs: true,
				},
				"second": {
					Cars:   true,
					Humans: true,
					Inputs: false,
				},
			},
			output: false,
		},
		{
			camerasStates: map[string]CameraStates{
				"first": {
					Cars:   true,
					Humans: false,
					Inputs: true,
				},
				"second": {
					Cars:   true,
					Humans: true,
					Inputs: false,
				},
			},
			output: false,
		},
		{
			camerasStates: map[string]CameraStates{
				"first": {
					Cars:   false,
					Humans: true,
					Inputs: false,
				},
				"second": {
					Cars:   true,
					Humans: true,
					Inputs: false,
				},
			},
			output: false,
		},
		{
			camerasStates: map[string]CameraStates{
				"first": {
					Cars:   false,
					Humans: false,
					Inputs: false,
				},
				"second": {
					Cars:   false,
					Humans: false,
					Inputs: false,
				},
			},
			output: false,
		},
	}

	mutex.RLock()
	defer mutex.RUnlock()
	camerasStatesCopy := camerasStates
	defer func() { camerasStates = camerasStatesCopy }()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			camerasStates = tc.camerasStates
			if actual := isStartCountdown(); tc.output != actual {
				t.Fail()
			}
		})
	}
}
