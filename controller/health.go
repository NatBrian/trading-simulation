package controller

import (
	"context"
	"net/http"

	"github.com/NatBrian/Stockbit-Golang-Challenge/helper"
	"github.com/redis/go-redis/v9"
)

type (
	HealthController struct {
		Context context.Context
		Redis   *redis.Client
	}

	IHealthController interface {
		Ping(w http.ResponseWriter, r *http.Request)
	}
)

func (hc *HealthController) Ping(w http.ResponseWriter, r *http.Request) {
	err := hc.Redis.Ping(hc.Context).Err()
	if err != nil {
		helper.ResponseFormatter(w, http.StatusInternalServerError, err)
		return
	}
	helper.ResponseFormatter(w, http.StatusOK, "OK")
}
