package repository

var getProducts = `SELECT * FROM product`

var storeProduct = `INSERT INTO product (product_uuid, name, description, price) values ($1, $2, $3, $4)`
