package main

import "hash_worker/internal/app"

func main() {
	cfg, err := app.LoadConfig()
	if err != nil {
		panic(err)
	}
	app := app.NewApp(cfg)
	app.RunApp()
}
