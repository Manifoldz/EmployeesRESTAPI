package entities

type Department struct {
	Id        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name" binding:"required"`
	Phone     string `json:"phone" db:"phone"`
	CompanyId int    `json:"company_id" db:"company_id" binding:"required"`
}
