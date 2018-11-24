package main

type Property struct {
	Device        *Device           `json:"device"`
	Name          string            `json:"name"`
	Value         string            `json:"value"`
	Description   map[string]string `json:"description"`
	Visible       bool              `json:"visbile"`
	FireAndForget bool              `json:"fire_and_forget"`
}
