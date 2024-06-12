package entities

type Passport struct {
	Id     int    `json:"id" db:"id"`
	Type   string `json:"type" db:"type" binding:"required"`
	Number string `json:"number" db:"number" binding:"required"`
	//EmployeeId int    `json:"employee_id" db:"employee_id" binding:"required"`
}
