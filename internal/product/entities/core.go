package entities

import "time"

type Product struct {
	CompanyId   int64     `db:"company_id"`
	UUID        string    `db:"uuid"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       int64     `db:"price"`
	CreatedOn   time.Time `db:"created_on"`
	CreatedBy   string    `db:"created_by"`
	UpdatedOn   time.Time `db:"updated_on"`
	UpdatedBy   string    `db:"updated_by"`
}
