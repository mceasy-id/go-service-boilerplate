package v1

import (
	"mceasy/service-demo/config"
	"mceasy/service-demo/internal/identity/identityentities"
	"mceasy/service-demo/internal/product"
	"mceasy/service-demo/internal/product/dtos"
	"mceasy/service-demo/internal/product/entities"
	"mceasy/service-demo/pkg/observability/instrumentation"
	"mceasy/service-demo/pkg/resourceful"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func NewProductHandler(config config.Config, productUC product.UseCase) product.Handlers {
	return &productHandler{
		config:    config,
		productUC: productUC,
	}
}

type productHandler struct {
	config    config.Config
	productUC product.UseCase
}

func (h *productHandler) Index(c *fiber.Ctx) error {
	ctx, span := instrumentation.NewTraceSpan(
		c.UserContext(),
		"IndexHandler",
	)
	defer span.End()

	var request dtos.IndexRequest
	err := c.QueryParser(&request)
	if err != nil {
		return err
	}

	err = request.Validate()
	if err != nil {
		return err
	}

	authCredential := c.Locals(identityentities.KeyAuthCredential).(identityentities.Credential)

	resource := resourceful.NewResource[uuid.UUID, dtos.ProductList](ProductDefinition)

	resource.Parameter = &resourceful.Parameter{
		Limit:   request.Limit,
		Page:    request.Page,
		Search:  request.Search,
		Filters: request.Filters,
		Sorts:   request.Sorts,
	}

	resourceProduct, err := h.productUC.Index(ctx, authCredential.CompanyId, resource)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(resourceProduct.Response())

	// products, err := h.productUC.Index(c.UserContext())
	// if err != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(err)
	// }

	// return c.JSON(products)
}

func (h *productHandler) Store(c *fiber.Ctx) error {
	ctx, span := instrumentation.NewTraceSpan(
		c.UserContext(),
		"StoreHandler",
	)
	defer span.End()

	var createProduct dtos.CreateProductRequest
	err := c.BodyParser(&createProduct)
	if err != nil {
		return err
	}

	err = createProduct.Mod().Validate()
	if err != nil {
		return err
	}

	productId, err := h.productUC.Store(
		ctx,
		c.Locals(identityentities.KeyAuthCredential).(identityentities.Credential),
		entities.StoreProduct{
			Name:        createProduct.Name,
			Description: createProduct.Description,
			Price:       createProduct.Price,
		},
	)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(
		entities.ResponseData{Data: map[string]any{"id": productId}},
	)

}

func (h *productHandler) Show(c *fiber.Ctx) error {
	ctx, span := instrumentation.NewTraceSpan(
		c.UserContext(),
		"ShowHandler",
	)
	defer span.End()

	var param struct {
		ProductUUID uuid.UUID `params:"productUUID"`
	}

	err := c.ParamsParser(&param)

	if err != nil {
		return err
	}

	companyId := c.
		Locals(identityentities.KeyAuthCredential).(identityentities.Credential).CompanyId

	product, err := h.productUC.Show(ctx,
		param.ProductUUID.String(),
		int64(companyId),
	)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(entities.ResponseData{
		Data: dtos.NewProductResponse(product),
	})
}

func (h *productHandler) Delete(c *fiber.Ctx) error {
	ctx, span := instrumentation.NewTraceSpan(
		c.UserContext(),
		"DeleteHandler",
	)
	defer span.End()

	var param struct {
		ProductUUID uuid.UUID `params:"productId"`
	}

	err := c.ParamsParser(&param)
	if err != nil {
		return err
	}

	companyId := c.Locals(
		identityentities.KeyAuthCredential,
	).(identityentities.Credential).CompanyId

	err = h.productUC.Delete(ctx, param.ProductUUID.String(), int64(companyId))
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)

}

func (h *productHandler) Update(c *fiber.Ctx) error {
	ctx, span := instrumentation.NewTraceSpan(
		c.UserContext(),
		"UpdateHandler",
	)
	defer span.End()

	var param struct {
		ProductUUID uuid.UUID `params:"productId"`
	}

	err := c.ParamsParser(&param)
	if err != nil {
		return err
	}

	var updateProduct dtos.UpdateProductRequest
	err = c.BodyParser(&updateProduct)
	if err != nil {
		return err
	}

	err = updateProduct.Mod().Validate()
	if err != nil {
		return err
	}

	requestCredential := c.Locals(
		identityentities.KeyAuthCredential,
	).(identityentities.Credential)

	err = h.productUC.Update(
		ctx,
		requestCredential,
		param.ProductUUID.String(),
		entities.UpdateProduct{
			CompanyId:   int64(requestCredential.CompanyId),
			UUID:        param.ProductUUID.String(),
			Name:        updateProduct.Name,
			Description: updateProduct.Description,
			Price:       updateProduct.Price,
			UpdatedOn:   time.Now(),
			UpdatedBy:   requestCredential.UserName,
		},
	)

	return nil
}
