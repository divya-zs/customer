package service

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"customer/models"
	"customer/store"
)

type customer struct {
	store store.ServiceIn
}

func New(c store.ServiceIn) customer {
	return customer{store: c}
}

func (c customer) Get(ctx *gofr.Context) ([]models.Customer, error) {
	res, err := c.store.Get(ctx)
	if err != nil {
		return []models.Customer{}, errors.DB{Err: errors.Error("db error")}
	}
	return res, nil
}

func (c customer) GetByID(ctx *gofr.Context, id int) (models.Customer, error) {
	res, err := c.store.GetByID(ctx, id)
	if err != nil {
		return models.Customer{}, err
	}
	return res, nil
}

func (c customer) Create(ctx *gofr.Context, customer models.Customer) (models.Customer, error) {
	res, err := c.store.Create(ctx, customer)
	if err != nil {
		return models.Customer{}, err
	}
	return res, nil
}

func (c customer) Update(ctx *gofr.Context, customer models.Customer) (models.Customer, error) {
	res, err := c.store.Update(ctx, customer.ID, customer)
	if err != nil {
		return models.Customer{}, err
	}
	return res, nil
}

func (c customer) Delete(ctx *gofr.Context, id int) error {
	err := c.store.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (c customer) Patch(ctx *gofr.Context, id int, customer models.Customer) (models.Customer, error) {
	return c.store.Patch(ctx, id, customer)
}
