package main

import (
	"fmt"
	"github.com/sol1corejz/goph-keeper/cmd/client/keepercli/cmd"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
)

func main() {
	// Вывод информации о версии сборки
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	cmd.Execute()
}
