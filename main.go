package main

import (
	"customer/handler"
	"customer/middleware"
	"customer/service"
	"customer/store"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	app := gofr.New()

	app.Server.UseMiddleware(middleware.OauthMiddleware)
	store := store.New()
	service := service.New(store)
	handler := handler.New(service)

	app.GET("/customer", handler.Get)
	app.GET("/customer/{id}", handler.GetByID)
	app.POST("/customer", handler.Create)
	app.PUT("/customer/{id}", handler.Update)
	app.DELETE("/customer/{id}", handler.Delete)
	app.PATCH("/customer/{id}", handler.Patch)
	app.Start()
}
