package lib

import (
	"log"

	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pair"
	"nanomsg.org/go/mangos/v2/protocol/req"

	_ "nanomsg.org/go/mangos/v2/transport/all"
)

type IpcClient struct {
	ManagerSocket mangos.Socket
	PluginSocket  mangos.Socket
}

const IPC_ROOT = "ipc:///tmp/"

func NewIpcClient(pluginId string) (*IpcClient, error) {
	log.Println("Starting managerSocket")

	var managerSocket mangos.Socket
	var err error

	managerSocket, err = req.NewSocket()
	if err != nil {
		log.Fatal(err)
	}

	err = managerSocket.Dial(IPC_ROOT + "gateway.addonManager")
	if err != nil {
		log.Fatal(err)
	}

	i := new(IpcClient)
	i.ManagerSocket = managerSocket
	//defer i.ManagerSocket.Close()

	return i, nil
}

func (i *IpcClient) RegisterPluginSocket(ipcBaseAddr string) {
	log.Println("Starting pluginSocket")

	var pluginSocket mangos.Socket
	var err error

	pluginSocket, err = pair.NewSocket()
	if err != nil {
		log.Fatal(err)
	}

	err = pluginSocket.Dial(IPC_ROOT + ipcBaseAddr)
	if err != nil {
		log.Fatal(err)
	}

	i.PluginSocket = pluginSocket
	//defer i.PluginSocket.Close()
}
