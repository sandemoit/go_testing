package models

type EmployeeRequest struct {
	Badge string `json:"badge"`
	Name  string `json:"name"`
}
