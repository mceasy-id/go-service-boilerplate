package entities

import "time"

type UpdateProduct struct {
	CompanyId   int64     `db:"company_id"`
	UUID        string    `db:"uuid"`
	Name        string    `db:"name"`
	Description string    `db:"name"`
	Price       int64     `db:"price"`
	UpdatedOn   time.Time `db:"updated_on"`
	UpdatedBy   string    `db:"updated_by"`
}
