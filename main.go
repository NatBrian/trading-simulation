package main

import (
	"log"
	"net/http"

	"github.com/NatBrian/Stockbit-Golang-Challenge/application"
	"github.com/NatBrian/Stockbit-Golang-Challenge/infrastructure"
)

func main() {
	app, err := application.SetupApp()
	if err != nil {
		log.Println("error: SetupApp", err)
		panic(err)
	}

	httpRouter := infrastructure.ServeHTTP(app)

	log.Println("ListenAndServe: " + app.Config.HTTPHost + ":" + app.Config.HTTPPort)
	err = http.ListenAndServe(app.Config.HTTPHost+":"+app.Config.HTTPPort, httpRouter)
	if err != nil {
		log.Println("error: ListenAndServe", err)
		return
	}
}
