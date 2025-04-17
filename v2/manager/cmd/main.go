package main

import (
	"fmt"
	"hash_manager/internal/app"
)

func main() {
	config, err := app.LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
	app, err := app.NewApp(config)
	if err != nil {
		panic(err)
	}
	app.StartApp()
}
