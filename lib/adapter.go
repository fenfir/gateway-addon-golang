package lib

import (
	"errors"
	"log"
)

type Adapter struct {
	Manager     AddonManager
	Id          string            `json:"id"`
	PackageName string            `json:"packageName"`
	Name        string            `json:"name"`
	Devices     map[string]Device `json:"devices"`
	Actions     map[string]Action `json:"actions"`
	Ready       bool              `json:"ready"`
}

func NewAdapter(addonManager *AddonManager, id string, packageName string) (*Adapter, error) {
	a := new(Adapter)
	a.Manager = *addonManager
	a.Id = id
	a.PackageName = packageName
	a.Name = "Adapter"
	a.Devices = make(map[string]Device)
	a.Actions = make(map[string]Action)
	a.Ready = true

	return a, nil
}

func (a *Adapter) Dump() {
	log.Printf("Adapter: %s Dump() not implemented")
}

func (a *Adapter) GetId() string {
	return a.Id
}

func (a *Adapter) GetPackageName() string {
	return a.PackageName
}

func (a *Adapter) GetDevice(id string) Device {
	return a.Devices[id]
}

func (a *Adapter) GetDevices() map[string]Device {
	return a.Devices
}

func (a *Adapter) GetName() string {
	return a.Name
}

func (a *Adapter) IsReady() bool {
	return a.Ready
}

func (a *Adapter) HandleDeviceAdded(device *Device) {
	a.Devices[device.Id] = *device
	a.Manager.HandleDeviceAdded(device)
}

func (a *Adapter) HandleDeviceRemoved(device *Device) {
	delete(a.Devices, device.Id)
	a.Manager.HandleDeviceRemoved(device)
}

func (a *Adapter) StartPairing(timeoutSeconds int) {
	log.Printf("Adapter: %s id %s pairing started", a.Name, a.Id)
}

func (a *Adapter) CancelPairing() {
	log.Printf("Adapter: %s id %s pairing cancelled", a.Name, a.Id)
}

func (a *Adapter) RemoveThing(device *Device) {
	log.Printf("Adapter: %s id %s RemoveThing(%s) started", a.Name, a.Id, device.Id)
}

func (a *Adapter) CancelRemoveThing(device *Device) {
	log.Printf("Adapter: %s id %s CancelRemoveThing(%s)", a.Name, a.Id, device.Id)
}

func (a *Adapter) Unload() {
	log.Printf("Adapter: %s unloaded", a.Name)
}

func (a *Adapter) SetPin(deviceId string, pin string) error {
	d := a.GetDevice(deviceId)
	if len(d.Id) > 0 {
		log.Printf("Adapter: %s id %s SetPin(%s, %s)", a.Name, a.Id, d.Id, pin)
		return nil
	}

	return errors.New("Device not found")
}
