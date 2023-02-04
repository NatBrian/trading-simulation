package application

import (
	"github.com/NatBrian/Stockbit-Golang-Challenge/controller"
	"github.com/NatBrian/Stockbit-Golang-Challenge/service"
)

type Dependency struct {
	StockController controller.IStockController
}

func SetupDependency(app App) Dependency {
	stockService := service.StockService{}

	stockController := controller.StockController{
		Config:       app.Config,
		StockService: stockService,
	}

	return Dependency{
		StockController: &stockController,
	}
}
