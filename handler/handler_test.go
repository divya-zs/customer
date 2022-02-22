package handler

import (
	"bytes"
	"customer/mocks"
	"customer/models"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func connect(r *http.Request) *gofr.Context {
	app := gofr.New()
	w := httptest.NewRecorder()
	req := request.NewHTTPRequest(r)
	res := responder.NewContextualResponder(w, r)
	ctx := gofr.NewContext(res, req, app)
	return ctx
}

func TestHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockHandlerIn(ctrl)
	h := New(m)

	customer1 := []models.Customer{{
		ID: 1, Name: "Divya", Age: 22, Salary: 30000,
	}}

	tests := []struct {
		desc     string
		mocks    []*gomock.Call
		expected interface{}
		err      error
	}{
		{"get all", []*gomock.Call{m.EXPECT().Get(gomock.Any()).Return(customer1, nil)},
			customer1, nil},
		{"internal server error",
			[]*gomock.Call{m.EXPECT().Get(gomock.Any()).Return([]models.Customer{}, errors.DB{Err: errors.Error("db error")})},
			[]models.Customer{}, errors.DB{Err: errors.Error("db error")}},
	}

	for i, tc := range tests {

		r := httptest.NewRequest(http.MethodGet, "http://customer", nil)
		ctx := connect(r)

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

func TestHandler_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockHandlerIn(ctrl)
	h := New(m)

	customer1 := models.Customer{
		ID: 1, Name: "Divya", Age: 22, Salary: 30000,
	}

	tests := []struct {
		desc     string
		ID       string
		expected interface{}
		err      error
		mocks    []*gomock.Call
	}{
		{"get by ID", "1", customer1, nil,
			[]*gomock.Call{m.EXPECT().GetByID(gomock.Any(), 1).Return(customer1, nil)}},
		{"missing ID", "", nil, errors.MissingParam{Param: []string{"id"}}, nil},
		{"ID not found", "2", models.Customer{}, errors.EntityNotFound{Entity: "id"},
			[]*gomock.Call{m.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(models.Customer{}, errors.EntityNotFound{Entity: "id"})}},
		{"invalid ID", "s", nil, errors.InvalidParam{Param: []string{"id"}}, nil},
		{"internal server error", "1", models.Customer{}, errors.DB{Err: errors.Error("db error")},
			[]*gomock.Call{m.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(models.Customer{}, errors.DB{Err: errors.Error("db error")})}},
	}

	for i, tc := range tests {

		r := httptest.NewRequest(http.MethodGet, "http://customer", nil)
		ctx := connect(r)

		t.Run(tc.desc, func(t *testing.T) {

			ctx.SetPathParams(map[string]string{
				"id": tc.ID,
			})
			resp, err := h.GetByID(ctx)

			if !reflect.DeepEqual(tc.expected, resp) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.expected, resp)
			}
			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.err, err)
			}
		})
	}
}

func TestHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockHandlerIn(ctrl)
	h := New(m)

	customer1 := models.Customer{
		ID: 1, Name: "Divya", Age: 22, Salary: 30000,
	}

	c1 := []byte(`{"id":1, "name": "divya", "age": 22, "salary": 30000}`)
	c2 := []byte(``)
	c3 := []byte(`{"id":1, "name": "", "age": 22, "salary": 30000}`)
	c4 := []byte(`{"id":1, "age": 22, "salary": 30000}`)

	tests := []struct {
		desc     string
		body     []byte
		expected interface{}
		err      error
		mock     []*gomock.Call
	}{
		{"success", c1, customer1, nil,
			[]*gomock.Call{m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(customer1, nil)}},
		{"invalid body 1", c2, nil, errors.InvalidParam{Param: []string{"body"}}, nil},
		{"invalid body 2", c3, models.Customer{}, errors.InvalidParam{Param: []string{"body"}},
			[]*gomock.Call{m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(models.Customer{}, errors.InvalidParam{Param: []string{"body"}})}},
		{"invalid body 3", c4, models.Customer{}, errors.InvalidParam{Param: []string{"body"}},
			[]*gomock.Call{m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(models.Customer{}, errors.InvalidParam{Param: []string{"body"}})}},
		{"internal server error", c1, models.Customer{}, errors.DB{Err: errors.Error("db err")},
			[]*gomock.Call{m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(models.Customer{}, errors.DB{Err: errors.Error("db err")})}},
	}

	for i, tc := range tests {

		r := httptest.NewRequest(http.MethodPost, "http://customer", bytes.NewReader(tc.body))
		ctx := connect(r)

		t.Run(tc.desc, func(t *testing.T) {

			resp, err := h.Create(ctx)
			if !reflect.DeepEqual(tc.expected, resp) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.expected, resp)
			}
			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.err, err)
			}
		})
	}
}

func TestHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockHandlerIn(ctrl)
	h := New(m)

	customer1 := models.Customer{
		ID: 1, Name: "Divya", Age: 22, Salary: 30000,
	}
	c1 := []byte(`{"id":1, "name": "divya", "age": 22, "salary": 30000}`)
	c2 := []byte(``)
	c3 := []byte(`{"id":1, "age":21, "salary": 30000}`)

	tests := []struct {
		desc     string
		ID       string
		body     []byte
		expected interface{}
		err      error
		mock     []*gomock.Call
	}{
		{"success", "1", c1, customer1, nil,
			[]*gomock.Call{m.EXPECT().Update(gomock.Any(), gomock.Any()).Return(customer1, nil)}},
		{"invalid body 1", "1", c2, nil, errors.InvalidParam{Param: []string{"body"}},
			nil},
		{"invalid body", "1", c3, models.Customer{},
			errors.InvalidParam{Param: []string{"body"}},
			[]*gomock.Call{m.EXPECT().Update(gomock.Any(), gomock.Any()).Return(models.Customer{}, errors.InvalidParam{Param: []string{"body"}})}},
		{"internal server error", "1", c1, models.Customer{},
			errors.DB{Err: errors.Error("db err")},
			[]*gomock.Call{m.EXPECT().Update(gomock.Any(), gomock.Any()).Return(models.Customer{}, errors.DB{Err: errors.Error("db err")})}},
		{"ID not found", "1", c1, models.Customer{}, errors.EntityNotFound{Entity: "ID"},
			[]*gomock.Call{m.EXPECT().Update(gomock.Any(), gomock.Any()).Return(models.Customer{}, errors.EntityNotFound{Entity: "ID"})}},
		{"empty ID", "", c1, nil, errors.MissingParam{Param: []string{"id"}}, nil},
		{"empty ID", "s", c1, nil, errors.InvalidParam{Param: []string{"id"}}, nil},
	}

	for i, tc := range tests {

		r := httptest.NewRequest(http.MethodPost, "http://customer", bytes.NewReader(tc.body))
		ctx := connect(r)

		t.Run(tc.desc, func(t *testing.T) {
			ctx.SetPathParams(map[string]string{
				"id": tc.ID,
			})
			resp, err := h.Update(ctx)
			if !reflect.DeepEqual(tc.expected, resp) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.expected, resp)
			}
			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.err, err)
			}
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockHandlerIn(ctrl)
	h := New(m)

	tests := []struct {
		desc string
		ID   string
		err  error
		mock []*gomock.Call
	}{
		{"success", "1", nil,
			[]*gomock.Call{m.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)}},
		{"missing ID", "", errors.MissingParam{Param: []string{"id"}}, nil},
		{"invalid ID", "s", errors.InvalidParam{Param: []string{"id"}}, nil},
		{"ID not found", "1", errors.EntityNotFound{Entity: "ID"},
			[]*gomock.Call{m.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.EntityNotFound{Entity: "ID"})}},
		{"internal server error", "1", errors.DB{Err: errors.Error("db error")},
			[]*gomock.Call{m.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.DB{Err: errors.Error("db error")})}},
	}

	for i, tc := range tests {

		r := httptest.NewRequest(http.MethodPost, "http://customer", nil)
		ctx := connect(r)

		t.Run(tc.desc, func(t *testing.T) {
			ctx.SetPathParams(map[string]string{
				"id": tc.ID,
			})
			_, err := h.Delete(ctx)

			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.err, err)
			}
		})
	}
}

func TestHandler_Patch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockHandlerIn(ctrl)
	h := New(m)

	customer1 := models.Customer{
		ID: 1, Name: "Divya", Age: 22, Salary: 30000,
	}
	c1 := []byte(`{"name": "divya"}`)
	c2 := []byte(`{"divya"}`)
	c3 := []byte(`{"id":1, "age":21, "salary": 30000}`)

	tests := []struct {
		desc     string
		ID       string
		body     []byte
		expected interface{}
		err      error
		mock     []*gomock.Call
	}{
		{"success", "1", c1, customer1, nil,
			[]*gomock.Call{m.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any()).Return(customer1, nil)}},
		{"empty ID", "", c1, nil, errors.MissingParam{Param: []string{"id"}}, nil},
		{"invalid ID", "abc", c1, nil, errors.InvalidParam{Param: []string{"id"}},
			nil},
		//{"parsing error", "1", []byte(`?{`), nil, errors.InvalidParam{Param: []string{"body"}}, nil},
		{"invalid body 1", "1", c2, nil, errors.Error("patch body error"), nil},
		//{"unmarshall error", "1", []byte(`"name":"divya":`), nil, errors.Error("unmarshal error"), nil},
		{"restricted field ID", "1", c3, nil, errors.InvalidParam{Param: []string{"id"}},
			nil},
	}

	for i, tc := range tests {
		r := httptest.NewRequest(http.MethodPost, "http://customer", bytes.NewReader(tc.body))
		ctx := connect(r)

		t.Run(tc.desc, func(t *testing.T) {
			ctx.SetPathParams(map[string]string{
				"id": tc.ID,
			})
			resp, err := h.Patch(ctx)
			if !reflect.DeepEqual(tc.expected, resp) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.expected, resp)
			}
			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("TEST[%d], failed.\n%s\nExpected %v\nGot %v", i+1, tc.desc, tc.err, err)
			}
		})
	}
}
