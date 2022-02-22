package main

import (
	"customer/handler"
	"customer/service"
	"customer/store"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"fmt"
	"github-lvs.corpzone.internalzone.com/mcafee/cnsr-gofr-csp-auth/validator"
)

func main() {
	app := gofr.New()
	opts := validator.Options{
		Keys: map[string]string{
			app.Config.Get("CSP_APP_KEY_CATALOG"): app.Config.Get("CSP_SHARED_KEY_CATALOG"),
		},
	}
	fmt.Println(opts)

	app.Server.UseMiddleware(validator.CSPAuth(app.Logger, opts))

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
