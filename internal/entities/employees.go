package entities

type Employee struct {
	Id           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name" binding:"required"`
	Surname      string `json:"surname" db:"surname" binding:"required"`
	Phone        string `json:"phone" db:"phone"`
	CompanyId    int    `json:"company_id" db:"company_id" binding:"required"`
	DepartmentId int    `json:"department_id" db:"department_id" binding:"required"`
}
