package repository

import (
	"context"
	"database/sql"
	"fmt"
	"mceasy/service-demo/internal/product"
	"mceasy/service-demo/internal/product/dtos"
	"mceasy/service-demo/internal/product/entities"
	"mceasy/service-demo/internal/product/tabledefinition"
	"mceasy/service-demo/pkg/apperror"
	"mceasy/service-demo/pkg/database"
	"mceasy/service-demo/pkg/observability/instrumentation"
	"mceasy/service-demo/pkg/resourceful"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func NewProductPGRepo(db *sqlx.DB) product.Repository {
	return &productRepo{
		conn: db,
		db:   db,
	}
}

type productRepo struct {
	conn *sqlx.DB
	db   database.Queryer
}

// Atomic implements product.Repository.
func (r *productRepo) Atomic(ctx context.Context, opt *sql.TxOptions, cb func(tx product.Repository) error) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"AtomicRepo",
	)
	defer span.End()

	tx, err := r.conn.BeginTxx(ctx, opt)
	if err != nil {
		return errors.Wrap(err, "productRepo.Atomic.BeginTxx")
	}
	newRepo := &productRepo{
		conn: r.conn,
		db:   tx,
	}

	err = cb(newRepo)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err : %v, rb err: %w", err, rbErr)
		}
		return errors.Wrap(err, "productRepo.Atomic")
	}

	return tx.Commit()
}

func (r *productRepo) StoreNewProduct(ctx context.Context, product entities.Product) (string, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"StoreNewProduct",
	)

	defer span.End()

	query, args, err := sqlx.Named(insertProduct, product)
	if err != nil {
		return "", errors.Wrap(err, "productRepo.StoreNewProduct.sqlxNamed")
	}

	query = r.db.Rebind(query)

	var returnedId string
	err = r.db.GetContext(
		ctx,
		&returnedId,
		query,
		args...,
	)
	if err != nil {
		return "", errors.Wrap(err, "productRepo.StoreNewProduct.Tx.GetContext")
	}

	if returnedId != product.UUID {
		return "", errors.Wrap(errors.New("Invalid operation"), "productRepo.StoreNewProduct.Tx.GetContext")
	}

	return product.UUID, nil

	// newUUID := uuid.New()

	// query := storeProduct
	// _, err := r.db.ExecContext(ctx,
	// 	query,
	// 	newUUID,
	// 	product.Name,
	// 	product.Description,
	// 	product.Price,
	// )

	// if err != nil {
	// 	return nil, err
	// }

	// var createdProduct models.Product
	// createdProduct.ID = newUUID

	// return &createdProduct, nil
}

func (r *productRepo) GetProductByUUID(ctx context.Context, productUUID string, companyId int64, options ...entities.GetProductOption) (entities.Product, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"GetProductByUUIDRepo",
	)
	defer span.End()

	query := getProductByUUID
	if len(options) >= 1 && options[0].PessimisticLocking {
		query += " FOR UPDATE"
	}

	var product entities.Product
	err := r.db.GetContext(ctx,
		&product,
		query,
		productUUID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Product{}, apperror.NotFound()
		}
		return entities.Product{}, errors.Wrap(err, "productRepo.GetProductByUUID.GetContext.getProductByUUID")
	}

	if companyId != 0 && companyId != product.CompanyId {
		return entities.Product{}, apperror.Forbidden()
	}

	return product, nil

}

func (r *productRepo) FindProductResourceful(ctx context.Context, resource *resourceful.Resource[uuid.UUID, dtos.ProductList]) (*resourceful.Resource[uuid.UUID, dtos.ProductList], error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"FindProductResourcefulRepo",
	)
	defer span.End()

	resource.Select([]*resourceful.Field{
		tabledefinition.Product.Field("uuid"),
	})

	query, args, err := resource.QueryAndArgs()
	if err != nil {
		return nil, errors.Wrap(err, "productRepo.FindProductResourceful.QueryAndArgsResourceful")
	}

	var uuids []uuid.UUID
	err = r.db.SelectContext(ctx, &uuids, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "productRepo.FindProductResourceful.SelectContextDB")
	}

	paginatedIds, err := resource.GetPaginatedResults(uuids)
	if err != nil {
		return nil, err
	}

	resource.Select([]*resourceful.Field{
		tabledefinition.Product.Field("uuid"),
		tabledefinition.Product.Field("name"),
		tabledefinition.Product.Field("description"),
	})

	query, args, err = resource.PopulateQueryArgs(tabledefinition.Product.Field("uuid"), paginatedIds)
	if err != nil {
		return nil, errors.Wrap(err, "productRepo.FindProductResourceful.PopulateQueryArgsResourceful")
	}

	var products []dtos.ProductList
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "productRepo.FindProductResourceful.QueryContextDB")
	}

	for rows.Next() {
		var product dtos.ProductList

		err := rows.Scan(
			&product.UUID,
			&product.Name,
			&product.Description,
		)
		if err != nil {
			return nil, errors.Wrap(err, "productRepo.FindProductResourceful.QueryContextDB")
		}

		products = append(products, product)
	}

	resource.SetResult(resourceful.Result[uuid.UUID, dtos.ProductList]{Ids: uuids, PaginatedResult: products})
	return resource, nil

}

func (r *productRepo) IsProductKeyExists(ctx context.Context, payload entities.StoreProduct, companyId int64) (bool, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"IsProductKeyExistsRepo",
	)
	defer span.End()

	var productUUID string
	err := r.db.GetContext(ctx,
		&productUUID,
		isProductKeyExists,
		companyId,
		payload.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, errors.Wrap(err, "productRepo.IsproductKeyExists.GetContext.isProductKeyExists")
	}

	return true, nil

}

func (r *productRepo) UpdateProductByUUID(ctx context.Context, product entities.UpdateProduct) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"UpdateProductByUUIDRepo",
	)
	defer span.End()

	var returnedProductUUID string
	err := r.db.GetContext(
		ctx,
		&returnedProductUUID,
		updateProductByUUID,
		product.UUID,
		product.CompanyId,
		product.Name,
		product.Description,
		product.Price,
		product.UpdatedOn,
		product.UpdatedBy,
	)
	if err != nil {
		return err
	}

	if returnedProductUUID != product.UUID {
		return apperror.NotFound()
	}

	return nil
}

func (r *productRepo) DeleteProductByUUID(ctx context.Context, productUUID string, companyId int64) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"DeleteProductByUUIDRepo",
	)

	defer span.End()

	var returnedProductUUID string
	err := r.db.GetContext(
		ctx,
		&returnedProductUUID,
		deleteProductByUUID,
		productUUID,
		companyId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.AppError{}
		}
		return err
	}

	return nil
}
