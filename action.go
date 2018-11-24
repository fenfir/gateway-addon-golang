package main

import "time"

type Input struct{}

type Action struct {
	Id            string            `json:"id"`
	Device        *Device           `json:"device"`
	Name          string            `json:"name"`
	Input         map[string]string `json:"input"`
	Status        string            `json:"status"`
	TimeRequested time.Time         `json:"timeRequested"`
	TimeCompleted time.Time         `json:"timeCompleted"`
}

func NewAction(id string, device *Device, name string, input map[string]string) (*Action, error) {
	a := new(Action)
	a.Id = id
	a.Device = device
	a.Name = name
	a.Input = input
	a.Status = "created"
	a.TimeRequested = time.Now()

	return a, nil
}

func (a *Action) Start() error {
	a.Status = "pending"
	a.Device.ActionNotify(a)

	return nil
}

func (a *Action) Finish() error {
	a.Status = "completed"
	a.TimeCompleted = time.Now()
	a.Device.ActionNotify(a)

	return nil
}
