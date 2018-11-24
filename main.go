package main

import (
	"log"
)

func ConnectAddon() {
	addonManager, err := NewAddonManager("golang", true)
	if err != nil {
		log.Fatal(err)
	}

	addonManager.SendRegistrationMessage()

	addonManager.ManagerClient(1)
}
