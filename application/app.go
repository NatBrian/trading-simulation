package application

import (
	"log"

	"github.com/NatBrian/Stockbit-Golang-Challenge/config"
)

// App contains application instances
type App struct {
	Config config.Config
}

func SetupApp() (App, error) {
	var app App

	log.Println("Setup App")

	loadConfig, err := config.LoadConfig()
	if err != nil {
		log.Println("error: LoadConfig", err)
		return App{}, err
	}
	app.Config = loadConfig

	return app, nil
}
