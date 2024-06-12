package entities

import (
	"errors"
	"reflect"
)

type Employee struct {
	Id           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Surname      string `json:"surname" db:"surname" `
	Phone        string `json:"phone" db:"phone"`
	CompanyId    int    `json:"company_id" db:"company_id"`
	DepartmentId int    `json:"department_id" db:"department_id"`
}

type EmployeeInputAndResponse struct {
	Id         int        `json:"id"`
	Name       string     `json:"name" binding:"required"`
	Surname    string     `json:"surname" binding:"required"`
	Phone      string     `json:"phone"`
	CompanyId  int        `json:"company_id" binding:"required"`
	Passport   Passport   `json:"passport" binding:"required"`
	Department Department `json:"department" binding:"required"`
}

type UpdateEmployeeInput struct {
	Id         *int                   `json:"id"`
	Name       *string                `json:"name"`
	Surname    *string                `json:"surname"`
	Phone      *string                `json:"phone"`
	CompanyId  *int                   `json:"company_id"`
	Passport   *PassportUpdateInput   `json:"passport"`
	Department *DepartmentUpdateInput `json:"department"`
}

func (i UpdateEmployeeInput) Validate() error {
	v := reflect.ValueOf(i)
	for j := 0; j < v.NumField(); j++ {
		field := v.Field(j)
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			// найдено ненулевое поле, значит не все поля nil
			return nil
		}
	}
	// если дошли до этого момента, значит все поля nil
	return errors.New("all fields are nil, check correct fields name")
}
