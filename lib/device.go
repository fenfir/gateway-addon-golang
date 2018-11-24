package lib

import "errors"

type Device struct {
	Adapter     *Adapter            `json:"adapter"`
	Id          string              `json:"id"`
	Type        string              `json:"type"`
	AtContext   string              `json:"@context"`
	AtType      string              `json:"@type"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Properties  map[string]Property `json:"properties"`
	Actions     map[string]string   `json:"actions"`
	Events      map[string]string   `json:"events"`
	UiHref      string              `json:"uiHref"`
	PinRequired bool                `json:"pinRequired"`
	PinPattern  string              `json:"pinPattern"`
}

func NewDevice(adapter *Adapter, id string) (*Device, error) {
	d := new(Device)
	d.Adapter = adapter
	d.Id = id
	d.Type = "thing"
	d.AtContext = "https://iot.mozilla.org/schemas"
	d.AtType = "[]"
	d.Name = ""
	d.Description = ""
	d.Properties = make(map[string]Property)
	d.Actions = make(map[string]string)
	d.Events = make(map[string]string)
	d.UiHref = ""
	d.PinRequired = false
	d.PinPattern = ""

	return d, nil
}

func (d *Device) NotifyPropertyChanged(property *Property) {
	d.Adapter.Manager.SendPropertyChangedNotification(property)
}

func (d *Device) ActionNotify(action *Action) {
	d.Adapter.Manager.SendActionStatusNotification(action)
}

func (d *Device) EventNotify(event *Event) {
	d.Adapter.Manager.SendEventNotification(event)
}

func (d *Device) RequestAction(actionId string, actionName string, input map[string]string) error {
	action := d.Actions[actionName]
	if action == "" {
		return errors.New("Action not found")
	}

	a, err := NewAction(actionId, d, actionName, input)
	if err != nil {
		return err
	}

	d.PerformAction(a)

	return nil
}

func (d *Device) RemoveAction(actionId string, actionName string) error {
	action := d.Actions[actionName]
	if action == "" {
		return errors.New("Action not found")
	}

	d.CancelAction(actionId, actionName)

	return nil
}

func (d *Device) PerformAction(action *Action) {

}

func (d *Device) CancelAction(actionId string, actionName string) {

}

func (d *Device) AddAction(actionName string, metadata string) {
	d.Actions[actionName] = metadata
}

func (d *Device) AddEvent(eventName string, metadata string) {
	d.Events[eventName] = metadata
}
