package main

import (
	"log"

	"github.com/fenfir/gateway-addon-golang/lib"
)

func main() {
	addonManager, err := lib.NewAddonManager("golang", true)
	if err != nil {
		log.Fatal(err)
	}

	addonManager.SendRegistrationMessage()

	addonManager.ManagerClient(1)
}
