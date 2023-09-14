package main

import (
	"mini-project/config"
	"mini-project/handlers"
	"mini-project/middleware"
	"mini-project/model"

	"github.com/labstack/echo/v4"
    "github.com/swaggo/echo-swagger"
    _ "mini-project/docs"
)

// @title Manufacturer Go API
// @version 1.0
// @description This is the Manufacturer Go API for managing equipment, users, and rental history.
// @termsOfService https://example.com/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT
func main() {
	db := config.InitDatabase()
	db.AutoMigrate(&model.User{}, &model.Equipment{}, &model.RentalHistory{})

	handlers.SetDB(db)

	e := echo.New()

	e.POST("/register", handlers.RegisterUserHandler)
	e.POST("/login", handlers.LoginUserHandler)

	e.POST("/top-up", handlers.TopUpUserHandler, middleware.JWTMiddleware)

	e.GET("/equipment", handlers.GetAllEquipmentHandler, middleware.JWTMiddleware)
	e.POST("/equipment", handlers.CreateEquipmentHandler, middleware.JWTMiddleware)
	e.PUT("/equipment/:id", handlers.UpdateEquipmentHandler, middleware.JWTMiddleware)
	e.DELETE("/equipment/:id", handlers.DeleteEquipmentHandler, middleware.JWTMiddleware)

	e.GET("/rental", handlers.GetAllRentalHistoryHandler, middleware.JWTMiddleware)
	e.POST("/rental", handlers.CreateRentalHistoryHandler, middleware.JWTMiddleware)
	e.PUT("/rental/:id", handlers.UpdateRentalHistoryHandler, middleware.JWTMiddleware)
	e.DELETE("/rental/:id", handlers.DeleteRentalHistoryHandler, middleware.JWTMiddleware)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Start(":8080")
}
