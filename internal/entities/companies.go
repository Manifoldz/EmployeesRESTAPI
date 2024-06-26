package entities

type Company struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name" binding:"required"`
}

type CompanyUpdateInput struct {
	Id   *int    `json:"id" db:"id"`
	Name *string `json:"name" db:"name"`
}
