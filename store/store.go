package store

import (
	"customer/models"
	"database/sql"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"fmt"
)

type store struct{}

func New() store {
	return store{}
}

func (s store) Get(ctx *gofr.Context) ([]models.Customer, error) {

	rows, err := ctx.DB().QueryContext(ctx, "SELECT * FROM customer")
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	var res []models.Customer

	for rows.Next() {
		var customer models.Customer
		err := rows.Scan(&customer.ID, &customer.Name, &customer.Age, &customer.Salary)

		if err != nil {
			return nil, errors.Error("scan error")
		}
		res = append(res, customer)
	}
	return res, nil

}

func (s store) GetByID(ctx *gofr.Context, id int) (models.Customer, error) {
	var customer models.Customer
	rows := ctx.DB().QueryRowContext(ctx, "SELECT * FROM customer where id=?", id)
	err := rows.Scan(&customer.ID, &customer.Name, &customer.Age, &customer.Salary)
	if err == sql.ErrNoRows {
		return models.Customer{}, sql.ErrNoRows
	}
	if err != nil {
		return models.Customer{}, errors.DB{Err: err}
	}
	return customer, nil
}

func (s store) Create(ctx *gofr.Context, customer models.Customer) (models.Customer, error) {
	_, err := ctx.DB().ExecContext(ctx, "INSERT INTO customer (name,age,salary) VALUES(?,?,?)",
		customer.Name, customer.Age, customer.Salary)
	if err != nil {
		return models.Customer{}, errors.DB{Err: err}
	}
	return customer, nil
}

func (s store) Update(ctx *gofr.Context, id int, customer models.Customer) (models.Customer, error) {
	_, err := ctx.DB().ExecContext(ctx, "UPDATE customer SET name=?,age=?,salary=? WHERE id=?",
		customer.Name, customer.Age, customer.Salary, id)
	if err == sql.ErrNoRows {
		return models.Customer{}, sql.ErrNoRows
	}
	if err != nil {
		return models.Customer{}, errors.DB{Err: err}
	}
	return customer, nil
}

func (s store) Delete(ctx *gofr.Context, id int) error {
	_, err := ctx.DB().ExecContext(ctx, "DELETE FROM customer where id=?", id)
	if err == sql.ErrNoRows {
		return sql.ErrNoRows
	}
	if err != nil {
		return errors.DB{Err: err}
	}
	return nil
}

func (s store) Patch(ctx *gofr.Context, id int, customer models.Customer) (models.Customer, error) {
	query := "UPDATE customer"
	set, qp := setClause(customer)

	// No value is passed for update
	if qp == nil {
		return models.Customer{}, nil
	}

	query = fmt.Sprintf("%v %v where id = ?", query, set)

	qp = append(qp, id)

	_, err := ctx.DB().ExecContext(ctx, query, qp...)
	if err == sql.ErrNoRows {
		return models.Customer{}, sql.ErrNoRows
	}
	if err != nil {
		return models.Customer{}, errors.DB{Err: err}
	}
	customer.ID = id
	return customer, nil
}

func setClause(s models.Customer) (set string, filed []interface{}) {
	set = `SET`

	if s.Name != "" {
		set += " name = ?,"
		filed = append(filed, s.Name)
	}

	if s.Age != 0 {
		set += " age = ?,"
		filed = append(filed, s.Age)
	}

	if s.Salary != 0 {
		set += " salary = ?"
		filed = append(filed, s.Salary)
	}

	if set == "SET" {
		return "", nil
	}

	return set, filed
}
