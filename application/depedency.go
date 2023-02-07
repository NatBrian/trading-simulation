package application

import (
	"github.com/NatBrian/Stockbit-Golang-Challenge/controller"
	"github.com/NatBrian/Stockbit-Golang-Challenge/service"
)

type Dependency struct {
	HealthContoller controller.IHealthController
	StockController controller.IStockController
}

func SetupDependency(app App) Dependency {
	stockService := service.StockService{
		Config:        app.Config,
		KafkaProducer: app.Kafka.Producer,
		Context:       app.Context,
		Redis:         app.Redis,
	}

	stockController := controller.StockController{
		Config:        app.Config,
		StockService:  stockService,
		Context:       app.Context,
		KafkaConsumer: app.Kafka.Consumer,
	}

	healthController := controller.HealthController{
		Context: app.Context,
		Redis:   app.Redis,
	}

	return Dependency{
		StockController: &stockController,
		HealthContoller: &healthController,
	}
}
