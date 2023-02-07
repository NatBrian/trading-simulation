package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/NatBrian/Stockbit-Golang-Challenge/application"
	"github.com/NatBrian/Stockbit-Golang-Challenge/infrastructure"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	app, err := application.SetupApp(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error: SetupApp")
		panic(err)
	}
	defer func() {
		cancel()
	}()

	infrastructure.ConsumeMessages(app)

	httpRouter := infrastructure.ServeHTTP(app)

	log.Info().Msg(fmt.Sprintf("ListenAndServe: " + app.Config.Constants.HTTPHost + ":" + app.Config.Constants.HTTPPort))
	err = http.ListenAndServe(app.Config.Constants.HTTPHost+":"+app.Config.Constants.HTTPPort, httpRouter)
	if err != nil {
		log.Error().Err(err).Msg("error: ListenAndServe")
		return
	}
}
