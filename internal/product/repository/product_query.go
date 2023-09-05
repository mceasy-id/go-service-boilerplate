package repository

const isProductKeyExists = `
SELECT
	p.uuid
FROM product p
WHERE p.company_id = $1
	AND LOWER(p.name) = LOWER($2)

`

const insertProduct = `
INSERT INTO product (
	company_id,
	uuid, 
	name, 
	description, 
	price,
	created_on,
	created_by,
	updated_on,
	updated_by
) 
values (
	:company_id,
	:uuid,
	:name,
	:description,
	:price,
	:created_on,
	:created_by,
	:updated_on,
	:updated_by
)
RETURNING uuid
`

const getProductByUUID = `
SELECT 
	p.company_id,
	p.uuid,
	p.name,
	p.description,
	p.price,
	p.created_on,
	p.created_by,
	p.updated_on,
	p.updated_by
FROM product p
WHERE p.uuid = $1
`

const updateProductByUUID = `
UPDATE product
SET 
	name = $3,
	description = $4,
	price = $5,
	updated_on = $6,
	updated_by = $7
WHERE uuid = $1
	AND company_id = $2
RETURNING uuid
`

const deleteProductByUUID = `
DELETE FROM product
WHERE uuid = $1 AND company_id = $2
RETURNING uuid
`
