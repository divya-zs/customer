package handler

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"customer/models"
	"customer/service"

	jsonpatch "github.com/evanphx/json-patch"
)

type Handler struct {
	service service.HandlerIn
}

func New(h service.HandlerIn) Handler {
	return Handler{service: h}
}

func (h Handler) Get(ctx *gofr.Context) (interface{}, error) {
	return h.service.Get(ctx)
}

func (h Handler) GetByID(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	if id == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}
	uid, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}
	return h.service.GetByID(ctx, uid)

}

func (h Handler) Create(ctx *gofr.Context) (interface{}, error) {
	var customer models.Customer
	if err := ctx.Bind(&customer); err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}
	return h.service.Create(ctx, customer)
}

func (h Handler) Update(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	if id == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}
	uid, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}
	var customer models.Customer
	if err := ctx.Bind(&customer); err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}
	customer.ID = uid
	return h.service.Update(ctx, customer)
}

func (h Handler) Delete(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	if id == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}
	uid, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}
	err = h.service.Delete(ctx, uid)
	return nil, err
}

func (h Handler) Patch(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	if id == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	uid, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	body, err1 := ioutil.ReadAll(ctx.Request().Body)
	if err1 != nil {
		return nil, errors.Error("parsing error")
	}

	resBytes, _ := json.Marshal(models.Customer{})

	patch, err3 := jsonpatch.MergePatch(resBytes, body)
	if err3 != nil {
		return nil, errors.Error("patch body error")
	}

	var customer models.Customer
	err = json.Unmarshal(patch, &customer)
	if err != nil {
		return nil, errors.Error("unmarshal error")
	}

	if customer.ID != 0 {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}
	return h.service.Patch(ctx, uid, customer)
}
