package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sol1corejz/goph-keeper/configs"
)

func main() {
	fmt.Println("GophKeeper Client v0.1")

	_, err := configs.LoadClientConfig("configs/client_config.yaml")
	if err != nil {
		log.Info("Failed to load client config", err.Error())
	}
}
