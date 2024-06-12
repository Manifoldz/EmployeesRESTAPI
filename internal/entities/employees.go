package entities

type Employee struct {
	Id           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Surname      string `json:"surname" db:"surname" `
	Phone        string `json:"phone" db:"phone"`
	CompanyId    int    `json:"company_id" db:"company_id"`
	DepartmentId int    `json:"department_id" db:"department_id"`
}

type EmployeeInputAndResponse struct {
	Name       string     `json:"name" binding:"required"`
	Surname    string     `json:"surname" binding:"required"`
	Phone      string     `json:"phone"`
	CompanyId  int        `json:"company_id" binding:"required"`
	Passport   Passport   `json:"passport" binding:"required"`
	Department Department `json:"department" binding:"required"`
}
