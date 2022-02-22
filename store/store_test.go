package store

import (
	"context"
	"customer/models"
	"database/sql"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"reflect"
	"testing"
)

func InitializeDb() (*sql.DB, sqlmock.Sqlmock, *gofr.Context, store) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Println(err)
	}
	g := gofr.Gofr{DataStore: datastore.DataStore{ORM: db}}
	ctx := gofr.NewContext(nil, nil, &g)
	ctx.Context = context.Background()
	store := New()
	return db, mock, ctx, store
}

func TestStore_Get(t *testing.T) {
	db, mock, ctx, store := InitializeDb()
	defer db.Close()

	customer1 := []models.Customer{{
		ID: 1, Name: "Divya", Age: 22, Salary: 30000,
	}}

	rows := sqlmock.NewRows([]string{"id", "name", "age", "salary", "scanError"}).AddRow(1, "Divya", 22, 30000, "scanError")
	query := "SELECT * FROM customer"
	tests := []struct {
		desc     string
		expected []models.Customer
		err      error
		mock     interface{}
	}{
		{"success", customer1, nil,
			mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"ID", "Name", "Age", "Salary"}).AddRow(1, "Divya", 22, 30000))},
		{"internal server error", nil, errors.DB{Err: errors.Error("db error")},
			mock.ExpectQuery(query).WillReturnError(errors.Error("db error"))},
		{"scan error", nil, errors.Error("scan error"), mock.ExpectQuery(query).WillReturnRows(rows)},
	}

	for i, tc := range tests {

		t.Run(tc.desc, func(t *testing.T) {
			res, err := store.Get(ctx)
			if !reflect.DeepEqual(err, tc.err) {
				t.Errorf("TEST[%d] Expected%v\nGot%v", i+1, tc.err, err)
			}
			if !reflect.DeepEqual(res, tc.expected) {
				t.Errorf("TEST[%d] Expected%v\nGot%v", i+1, tc.expected, res)
			}
		})
	}
}

func TestStore_GetByID(t *testing.T) {
	db, mock, ctx, store := InitializeDb()
	defer db.Close()

	customer1 := models.Customer{ID: 1, Name: "Divya", Age: 22, Salary: 30000}
	query := "SELECT * FROM customer where id=?"
	tests := []struct {
		desc     string
		id       int
		expected models.Customer
		err      error
		mock     interface{}
	}{
		{"success", 1, customer1, nil,
			mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"ID", "Name", "Age", "Salary"}).AddRow(1, "Divya", 22, 30000))},
		{"internal server error", 1, models.Customer{}, errors.DB{Err: errors.Error("db error")},
			mock.ExpectQuery(query).WithArgs(1).WillReturnError(errors.Error("db error"))},
		{"ID not found", 5, models.Customer{}, sql.ErrNoRows,
			mock.ExpectQuery(query).WithArgs(5).WillReturnError(sql.ErrNoRows)},
	}

	for i, tc := range tests {

		t.Run(tc.desc, func(t *testing.T) {
			res, err := store.GetByID(ctx, tc.id)
			//assert.Equal(t, tc.err, err, "TEST[%d] Expected%v\nGot%v", i+1, tc.err, err)
			if !reflect.DeepEqual(err, tc.err) {
				t.Errorf("TEST[%d] Expected%v\nGot%v", i+1, tc.err, err)
			}
			if !reflect.DeepEqual(res, tc.expected) {
				t.Errorf("TEST[%d] Expected%v\nGot%v", i+1, tc.expected, res)
			}
		})
	}
}

func TestStore_Create(t *testing.T) {
	db, mock, ctx, store := InitializeDb()
	defer db.Close()

	customer1 := models.Customer{ID: 0, Name: "Divya", Age: 22, Salary: 30000}

	query := "INSERT INTO customer (name,age,salary) VALUES(?,?,?)"
	tests := []struct {
		desc     string
		input    models.Customer
		expected models.Customer
		err      error
		mock     interface{}
	}{
		{"success", customer1, customer1, nil,
			mock.ExpectExec(query).WithArgs("Divya", 22, 30000).WillReturnResult(sqlmock.NewResult(1, 1))},
		{"internal server error", customer1, models.Customer{}, errors.DB{Err: errors.Error("db error")},
			mock.ExpectExec(query).WillReturnError(errors.Error("db error"))},
	}

	for i, tc := range tests {

		t.Run(tc.desc, func(t *testing.T) {
			res, err := store.Create(ctx, tc.input)
			if !reflect.DeepEqual(err, tc.err) {
				t.Errorf("TEST[%d] Expected %v\nGot %v", i+1, tc.err, err)
			}
			if !reflect.DeepEqual(res, tc.expected) {
				t.Errorf("TEST[%d] Expected %v\nGot %v", i+1, tc.expected, res)
			}
		})
	}
}

func TestStore_Update(t *testing.T) {
	db, mock, ctx, store := InitializeDb()
	defer db.Close()

	customer1 := models.Customer{ID: 1, Name: "Divya", Age: 22, Salary: 30000}

	query := "UPDATE customer SET name=?,age=?,salary=? WHERE id=?"
	tests := []struct {
		desc     string
		ID       int
		input    models.Customer
		expected models.Customer
		err      error
		mock     interface{}
	}{
		{"success", customer1.ID, customer1, customer1, nil,
			mock.ExpectExec(query).WithArgs("Divya", 22, 30000, 1).WillReturnResult(sqlmock.NewResult(1, 1))},
		{"internal server error", customer1.ID, customer1, models.Customer{},
			errors.DB{Err: errors.Error("db error")}, mock.ExpectExec(query).WillReturnError(errors.Error("db error"))},
		{"invalid id", 1, customer1, models.Customer{}, sql.ErrNoRows,
			mock.ExpectExec(query).WillReturnError(sql.ErrNoRows)},
	}

	for i, tc := range tests {

		t.Run(tc.desc, func(t *testing.T) {
			res, err := store.Update(ctx, tc.ID, tc.input)
			if !reflect.DeepEqual(err, tc.err) {
				t.Errorf("TEST[%d] Expected %v\nGot %v", i+1, tc.err, err)
			}
			if !reflect.DeepEqual(res, tc.expected) {
				t.Errorf("TEST[%d] Expected %v\nGot %v", i+1, tc.expected, res)
			}
		})
	}
}

func TestStore_Delete(t *testing.T) {
	db, mock, ctx, store := InitializeDb()
	defer db.Close()

	customer1 := models.Customer{ID: 1, Name: "Divya", Age: 22, Salary: 30000}

	query := "DELETE FROM customer where id=?"

	mock.ExpectExec(query).WithArgs(customer1.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(query).WillReturnError(errors.Error("db error"))
	mock.ExpectExec(query).WillReturnError(sql.ErrNoRows)

	tests := []struct {
		desc string
		ID   int
		err  error
	}{
		{"success", customer1.ID, nil},
		{"internal server error", customer1.ID, errors.DB{Err: errors.Error("db error")}},
		{"invalid id", 2, sql.ErrNoRows},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			err := store.Delete(ctx, tc.ID)
			if !reflect.DeepEqual(err, tc.err) {
				t.Errorf("TEST[%d] Expected %v\nGot %v", i+1, tc.err, err)
			}
		})
	}
}

func TestStore_Patch(t *testing.T) {
	db, mock, ctx, store := InitializeDb()
	defer db.Close()

	customer1 := models.Customer{ID: 1, Name: "Divya", Age: 22, Salary: 30000}

	query := "UPDATE customer SET name = ?, age = ?, salary = ? where id = ?"
	query1 := "UPDATE customer"
	tests := []struct {
		desc     string
		ID       int
		input    models.Customer
		expected models.Customer
		err      error
		mock     interface{}
	}{
		{"success", customer1.ID, customer1, customer1, nil,
			mock.ExpectExec(query).WithArgs("Divya", 22, 30000, 1).WillReturnResult(sqlmock.NewResult(1, 1))},
		{"internal server error", customer1.ID, customer1, models.Customer{},
			errors.DB{Err: errors.Error("db error")}, mock.ExpectExec(query).WillReturnError(errors.Error("db error"))},
		{"no values to patch", 1, models.Customer{}, models.Customer{}, nil,
			mock.ExpectExec(query1).WillReturnError(nil)},
	}

	for i, tc := range tests {

		t.Run(tc.desc, func(t *testing.T) {
			res, err := store.Patch(ctx, tc.ID, tc.input)
			if !reflect.DeepEqual(err, tc.err) {
				t.Errorf("TEST[%d] Expected %v\nGot %v", i+1, tc.err, err)
			}
			if !reflect.DeepEqual(res, tc.expected) {
				t.Errorf("TEST[%d] Expected %v\nGot %v", i+1, tc.expected, res)
			}
		})
	}
}
