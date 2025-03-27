package main

import (
	"fmt"
	"hash_manager/internal/app"
	"log"
)

func main() {
	defer func() { log.Println("Finish manager") }()

	config, err := app.LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
	app := app.NewApp(config)
	app.StartApp()
}
