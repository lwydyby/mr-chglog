package main

import (
	"log"
	"os"

	"github.com/lwydyby/mr-chglog/cmd"
)

func main() {
	app := cmd.CreateApp(cmd.AppAction)
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
