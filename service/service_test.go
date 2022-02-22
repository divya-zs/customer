package service

import (
	"context"
	"customer/mocks"
	"customer/models"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func connect(t *testing.T) (*gomock.Controller, customer, *mocks.MockServiceIn, *gofr.Gofr) {
	ctrl := gomock.NewController(t)

	m := mocks.NewMockServiceIn(ctrl)
	h := New(m)
	app := gofr.New()
	return ctrl, h, m, app
}

func TestCustomer_Get(t *testing.T) {
	ctrl, h, m, app := connect(t)
	defer ctrl.Finish()
	customer1 := []models.Customer{{
		ID: 1, Name: "Divya", Age: 22, Salary: 30000,
	}}

	tests := []struct {
		desc     string
		mocks    []*gomock.Call
		expected []models.Customer
		err      error
	}{
		{"get all", []*gomock.Call{m.EXPECT().Get(gomock.Any()).Return(customer1, nil)}, customer1, nil},
		{"internal server error",
			[]*gomock.Call{m.EXPECT().Get(gomock.Any()).Return([]models.Customer{}, errors.DB{Err: errors.Error("db error")})},
			[]models.Customer{}, errors.DB{Err: errors.Error("db error")}},
	}

	for i, tc := range tests {
		ctx := gofr.NewContext(nil, nil, app)
		ctx.Context = context.Background()

		t.Run(tc.desc, func(t *testing.T) {

			resp, err := h.Get(ctx)
			if !reflect.DeepEqual(tc.expected, resp) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.expected, resp)
			}
			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.err, err)
			}
		})
	}
}

func TestCustomer_GetByID(t *testing.T) {
	ctrl, h, m, app := connect(t)
	defer ctrl.Finish()

	customer1 := models.Customer{
		ID: 1, Name: "Divya", Age: 22, Salary: 30000,
	}

	tests := []struct {
		desc     string
		ID       int
		expected interface{}
		err      error
		mocks    []*gomock.Call
	}{
		{"get by ID", 1, customer1, nil,
			[]*gomock.Call{m.EXPECT().GetByID(gomock.Any(), 1).Return(customer1, nil)}},
		{"ID not found", 2, models.Customer{}, errors.EntityNotFound{Entity: "id"},
			[]*gomock.Call{m.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(models.Customer{}, errors.EntityNotFound{Entity: "id"})}},
		{"internal server error", 1, models.Customer{}, errors.DB{Err: errors.Error("db error")},
			[]*gomock.Call{m.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(models.Customer{}, errors.DB{Err: errors.Error("db error")})}},
	}

	for i, tc := range tests {
		ctx := gofr.NewContext(nil, nil, app)
		ctx.Context = context.Background()

		t.Run(tc.desc, func(t *testing.T) {
			resp, err := h.GetByID(ctx, tc.ID)

			if !reflect.DeepEqual(tc.expected, resp) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.expected, resp)
			}
			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.err, err)
			}
		})
	}
}

func TestCustomer_Create(t *testing.T) {
	ctrl, h, m, app := connect(t)
	defer ctrl.Finish()

	customer1 := models.Customer{
		ID: 1, Name: "Divya", Age: 22, Salary: 30000,
	}
	customer2 := models.Customer{
		ID: 1, Name: "", Age: 22, Salary: 3000,
	}

	tests := []struct {
		desc     string
		input    models.Customer
		expected interface{}
		err      error
		mock     []*gomock.Call
	}{
		{"success", customer1, customer1, nil,
			[]*gomock.Call{m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(customer1, nil)}},
		{"internal server error", customer2, models.Customer{}, errors.DB{Err: errors.Error("db err")},
			[]*gomock.Call{m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(models.Customer{}, errors.DB{Err: errors.Error("db err")})}},
		{"internal server error", customer1, models.Customer{}, errors.DB{Err: errors.Error("db err")},
			[]*gomock.Call{m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(models.Customer{}, errors.DB{Err: errors.Error("db err")})}},
	}

	for i, tc := range tests {
		ctx := gofr.NewContext(nil, nil, app)
		ctx.Context = context.Background()

		t.Run(tc.desc, func(t *testing.T) {

			resp, err := h.Create(ctx, tc.input)
			if !reflect.DeepEqual(tc.expected, resp) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.expected, resp)
			}
			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.err, err)
			}
		})
	}
}

func TestCustomer_Update(t *testing.T) {
	ctrl, h, m, app := connect(t)
	defer ctrl.Finish()

	customer1 := models.Customer{
		ID: 1, Name: "Divya", Age: 22, Salary: 30000,
	}

	tests := []struct {
		desc     string
		input    models.Customer
		expected interface{}
		err      error
		mock     []*gomock.Call
	}{
		{"success", customer1, customer1, nil,
			[]*gomock.Call{m.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(customer1, nil)}},
		{"internal server error", customer1, models.Customer{}, errors.DB{Err: errors.Error("db err")},
			[]*gomock.Call{m.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Customer{}, errors.DB{Err: errors.Error("db err")})}},
		{"ID not found", customer1, models.Customer{}, errors.EntityNotFound{Entity: "ID"},
			[]*gomock.Call{m.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Customer{}, errors.EntityNotFound{Entity: "ID"})}},
	}

	for i, tc := range tests {
		ctx := gofr.NewContext(nil, nil, app)
		ctx.Context = context.Background()

		t.Run(tc.desc, func(t *testing.T) {

			resp, err := h.Update(ctx, tc.input)
			if !reflect.DeepEqual(tc.expected, resp) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.expected, resp)
			}
			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.err, err)
			}
		})
	}
}

func TestCustomer_Delete(t *testing.T) {
	ctrl, h, m, app := connect(t)
	defer ctrl.Finish()

	tests := []struct {
		desc string
		ID   int
		err  error
		mock []*gomock.Call
	}{
		{"success", 1, nil,
			[]*gomock.Call{m.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)}},
		{"invalid ID", 0, errors.InvalidParam{Param: []string{"id"}},
			[]*gomock.Call{m.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.InvalidParam{Param: []string{"id"}})}},
		{"ID not found", 1, errors.EntityNotFound{Entity: "ID"},
			[]*gomock.Call{m.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.EntityNotFound{Entity: "ID"})}},
		{"internal server error", 1, errors.DB{Err: errors.Error("db error")},
			[]*gomock.Call{m.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.DB{Err: errors.Error("db error")})}},
	}

	for i, tc := range tests {
		ctx := gofr.NewContext(nil, nil, app)
		ctx.Context = context.Background()

		t.Run(tc.desc, func(t *testing.T) {
			err := h.Delete(ctx, tc.ID)
			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.err, err)
			}
		})
	}
}

func TestCustomer_Patch(t *testing.T) {
	ctrl, h, m, app := connect(t)
	defer ctrl.Finish()
	customer1 := models.Customer{
		ID: 1, Name: "Divya", Age: 22, Salary: 30000,
	}
	tests := []struct {
		desc     string
		ID       int
		input    models.Customer
		expected models.Customer
		err      error
		mock     []*gomock.Call
	}{
		{"success", 1, customer1, customer1, nil,
			[]*gomock.Call{m.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any()).Return(customer1, nil)}},
		{"invalid ID", 0, models.Customer{}, models.Customer{},
			errors.InvalidParam{Param: []string{"id"}},
			[]*gomock.Call{m.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Customer{}, errors.InvalidParam{Param: []string{"id"}})}},
		{"ID not found", 1, customer1, models.Customer{}, errors.EntityNotFound{Entity: "ID"},
			[]*gomock.Call{m.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Customer{}, errors.EntityNotFound{Entity: "ID"})}},
		{"internal server error", 1, customer1, models.Customer{}, errors.DB{Err: errors.Error("db error")},
			[]*gomock.Call{m.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Customer{}, errors.DB{Err: errors.Error("db error")})}},
	}

	for i, tc := range tests {
		ctx := gofr.NewContext(nil, nil, app)
		ctx.Context = context.Background()

		t.Run(tc.desc, func(t *testing.T) {
			resp, err := h.Patch(ctx, tc.ID, tc.input)
			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.err, err)
			}
			if !reflect.DeepEqual(tc.expected, resp) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.expected, resp)
			}
		})
	}
}
