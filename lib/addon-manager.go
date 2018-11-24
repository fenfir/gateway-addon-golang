package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"nanomsg.org/go/mangos/v2"

	_ "nanomsg.org/go/mangos/v2/transport/all"
)

var lock sync.Mutex

type AddonManager struct {
	Adapters  []Adapter
	IpcClient *IpcClient
	PluginId  string
	Verbose   bool
	Running   bool
}

func NewAddonManager(pluginId string, verbose bool) (*AddonManager, error) {
	a := new(AddonManager)
	a.Adapters = []Adapter{}
	ipcClient, err := NewIpcClient(pluginId)
	if err != nil {
		return nil, err
	}
	a.IpcClient = ipcClient
	a.PluginId = pluginId
	a.Verbose = verbose
	a.Running = true

	return a, nil
}

func (a *AddonManager) Close() {
	a.IpcClient.ManagerSocket.Close()
	a.IpcClient.PluginSocket.Close()
}

type BaseMessage struct {
	PluginId string `json:"pluginId"`
}

type ErrorMessage struct {
	BaseMessage
	Message string `json:"message"`
}

func (a *AddonManager) SendError(message string) error {
	e := ErrorMessage{
		Message: message,
	}

	e.PluginId = a.PluginId

	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	a.PluginSend(PLUGIN_ERROR, data)
	return nil
}

type AddAdapterMessage struct {
	BaseMessage
	AdapterId   string `json:"adapterId"`
	Name        string `json:"name"`
	PackageName string `json:"packageName"`
}

func (a *AddonManager) AddAdapter(adapter *Adapter) error {
	e := AddAdapterMessage{
		AdapterId:   adapter.Id,
		Name:        adapter.Name,
		PackageName: adapter.PackageName,
	}

	e.PluginId = a.PluginId

	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	a.PluginSend(ADD_ADAPTER, data)
	return nil
}

type HandleDeviceAddedMessage struct {
	BaseMessage
	AdapterId string `json:"adapterId"`
}

func (a *AddonManager) HandleDeviceAdded(device *Device) error {
	e := AddAdapterMessage{
		AdapterId: device.Adapter.Id,
	}

	e.PluginId = a.PluginId

	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	a.PluginSend(HANDLE_DEVICE_ADDED, data)
	return nil
}

type HandleDeviceRemovedMessage struct {
	BaseMessage
	AdapterId string `json:"adapterId"`
	Id        string `json:"id"`
}

func (a *AddonManager) HandleDeviceRemoved(device *Device) error {
	e := HandleDeviceRemovedMessage{
		AdapterId: device.Adapter.Id,
		Id:        device.Id,
	}

	e.PluginId = a.PluginId

	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	a.PluginSend(HANDLE_DEVICE_REMOVED, data)
	return nil
}

type SendPropertyChangedNotificationMessage struct {
	BaseMessage
	AdapterId string    `json:"adapterId"`
	DeviceId  string    `json:"deviceId"`
	Property  *Property `json:"property"`
}

func (a *AddonManager) SendPropertyChangedNotification(property *Property) error {
	e := SendPropertyChangedNotificationMessage{
		AdapterId: property.Device.Adapter.Id,
		DeviceId:  property.Device.Id,
		Property:  property,
	}

	e.PluginId = a.PluginId

	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	a.PluginSend(PROPERTY_CHANGED, data)
	return nil
}

type SendActionStatusNotificationMessage struct {
	BaseMessage
	AdapterId string  `json:"adapterId"`
	DeviceId  string  `json:"deviceId"`
	Action    *Action `json:"action"`
}

func (a *AddonManager) SendActionStatusNotification(action *Action) error {
	e := SendActionStatusNotificationMessage{
		AdapterId: action.Device.Adapter.Id,
		DeviceId:  action.Device.Id,
		Action:    action,
	}

	e.PluginId = a.PluginId

	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	a.PluginSend(ACTION_STATUS, data)
	return nil
}

type SendEventNotificationMessage struct {
	BaseMessage
	AdapterId string `json:"adapterId"`
	DeviceId  string `json:"deviceId"`
	Event     *Event `json:"event"`
}

func (a *AddonManager) SendEventNotification(event *Event) error {
	e := SendEventNotificationMessage{
		AdapterId: event.Device.Adapter.Id,
		DeviceId:  event.Device.Id,
		Event:     event,
	}

	e.PluginId = a.PluginId

	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	a.PluginSend(EVENT, data)
	return nil
}

type SendConnectedNotificationMessage struct {
	BaseMessage
	AdapterId string `json:"adapterId"`
	DeviceId  string `json:"deviceId"`
	Connected bool   `json:"connected"`
}

func (a *AddonManager) SendConnectedNotification(device *Device, connected bool) error {
	e := SendConnectedNotificationMessage{
		AdapterId: device.Adapter.Id,
		DeviceId:  device.Id,
		Connected: connected,
	}

	e.PluginId = a.PluginId

	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	a.PluginSend(CONNECTED, data)
	return nil
}

func (a *AddonManager) PluginSend(messageType string, data []byte) error {
	log.Printf("ManagerSend: %s\n", string(data))
	m := mangos.NewMessage(len(data))
	m.Body = data
	err := a.IpcClient.PluginSocket.SendMsg(m)
	if err != nil {
		return err
	}

	return nil
}

func (a *AddonManager) ManagerSend(data []byte) error {
	log.Printf("ManagerSend: %s\n", string(data))
	m := mangos.NewMessage(len(data))
	m.Body = data
	err := a.IpcClient.ManagerSocket.SendMsg(m)
	if err != nil {
		return err
	}

	return nil
}

func (a *AddonManager) ManagerRecv() {
	log.Println("Create new socket")
	managerSocket := a.IpcClient.ManagerSocket

	for {
		msg, err := managerSocket.Recv()
		if err != nil {
		}

		if msg != nil {
			log.Printf("ManagerRecv: %s\n", string(msg))
			a.IpcClient.RegisterPluginSocket("gateway.plugin.golang")
		}
	}
}

func (a *AddonManager) ManagerClient(nworkers int) {
	log.Printf("Starting %d ManagerClient\n", nworkers)
	a.ManagerRecv()
}

func (a *AddonManager) SendRegistrationMessage() {
	log.Println("Sending registration message")
	registrationMessage := fmt.Sprintf(`{
		"messageType": "registerPlugin",
		"data": {
			"pluginId": "%s"
		}
	}`, a.PluginId)

	err := a.ManagerSend([]byte(registrationMessage))
	if err != nil {
		log.Println(err)
	}
}
