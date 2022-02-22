package service

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"customer/models"
)

type HandlerIn interface {
	Get(ctx *gofr.Context) ([]models.Customer, error)
	GetByID(ctx *gofr.Context, id int) (models.Customer, error)
	Create(ctx *gofr.Context, customer models.Customer) (models.Customer, error)
	Update(ctx *gofr.Context, customer models.Customer) (models.Customer, error)
	Delete(ctx *gofr.Context, id int) error
	Patch(ctx *gofr.Context, id int, customer models.Customer) (models.Customer, error)
}
