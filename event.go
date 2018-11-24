package main

import "time"

type Event struct {
	Device    *Device   `json:"device"`
	Name      string    `json:"name"`
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}
