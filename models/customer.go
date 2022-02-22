package models

type Customer struct {
	ID     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Age    int    `json:"age,omitempty"`
	Salary int    `json:"salary,omitempty"`
}
