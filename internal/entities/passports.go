package entities

type Passport struct {
	Id     int    `json:"-" db:"id"`
	Type   string `json:"type" db:"type" binding:"required"`
	Number string `json:"number" db:"number" binding:"required"`
	//EmployeeId int    `json:"employee_id" db:"employee_id" binding:"required"`
}

type PassportUpdateInput struct {
	Id     *int    `json:"-" db:"id"`
	Type   *string `json:"type" db:"type"`
	Number *string `json:"number" db:"number"`
}
