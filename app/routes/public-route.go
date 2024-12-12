package routes

import (
	"creditlimit-connector/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	creditController := controllers.InitCreditController()

	route := app.Group("/connector/v1")

	route.Get("/availableCreditLimits", creditController.CheckCreditLimit)
	route.Get("/adjustCreditLimits", creditController.AdjustCreditLimit)
	route.Get("/isSbaBusiness", creditController.CheckOperationalDate)
}
