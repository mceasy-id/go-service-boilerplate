package apperror_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"mceasy/service-demo/constants"
	"mceasy/service-demo/pkg/apperror"
	"mceasy/service-demo/pkg/optional"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/invopop/validation"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func newFiberApp(err error) *fiber.App {
	var fiberConfig fiber.Config
	fiberConfig.ErrorHandler = apperror.HttpHandleError
	fiberApp := fiber.New(fiberConfig)

	fiberApp.Get("/", func(c *fiber.Ctx) error {
		if err != nil {
			return err
		}

		return nil
	})

	return fiberApp
}

func TestHttp_HandleError(t *testing.T) {
	t.Run("should catch wrapped jwtErr", func(t *testing.T) {
		jwtToken := "eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiQWRtaW4iLCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6IkphdmFJblVzZSIsImV4cCI6MTY4NDcyNzE4MSwiaWF0IjoxNjg0NzI3MTgxfQ.O8fAddg82hc6FabLCeo9Yu94f5tlqxA69yAgGbRusaA"
		var registeredClaim jwt.RegisteredClaims
		_, jwtErr := jwt.ParseWithClaims(jwtToken, &registeredClaim, func(t *jwt.Token) (interface{}, error) {
			return "invalid key", nil
		})
		err := errors.Wrap(jwtErr, "wraping jwt cases")
		fiberApp := newFiberApp(err)
		request := httptest.NewRequest(fiber.MethodGet, "/", nil)

		response, err := fiberApp.Test(request)
		require.NoError(t, err)

		defer response.Body.Close()

		require.Equal(t, http.StatusUnauthorized, response.StatusCode)
		var respErr apperror.Error
		err = json.NewDecoder(response.Body).Decode(&respErr)
		require.NoError(t, err)
		require.Equal(t, constants.UNAUTHORIZED_ERROR, respErr.Message)
	})

	t.Run("should catch wrapped numErr", func(t *testing.T) {
		_, numErr := strconv.Atoi("AIF")
		err := errors.Wrap(numErr, "wraping numErr cases")
		fiberApp := newFiberApp(err)
		request := httptest.NewRequest(fiber.MethodGet, "/", nil)

		response, err := fiberApp.Test(request)
		require.NoError(t, err)

		defer response.Body.Close()

		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		var respErr apperror.Error
		err = json.NewDecoder(response.Body).Decode(&respErr)
		require.NoError(t, err)
		require.Equal(t, constants.MALFORMED_BODY_ERROR, respErr.Message)
	})

	t.Run("should catch wrapped appErr", func(t *testing.T) {
		appErr := apperror.BadRequestMap(map[string][]string{
			"field.1": {"invalid value"},
		})
		err := errors.Wrap(appErr, "wraping httpErr cases")
		fiberApp := newFiberApp(err)
		request := httptest.NewRequest(fiber.MethodGet, "/", nil)

		response, err := fiberApp.Test(request)
		require.NoError(t, err)

		defer response.Body.Close()

		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		var respErr apperror.Error
		err = json.NewDecoder(response.Body).Decode(&respErr)
		require.NoError(t, err)
		require.Equal(t, "bad_request", respErr.Message)
		require.Equal(t, map[string]any{
			"field.1": []any{"invalid value"},
		}, respErr.Errors)
	})

	t.Run("should catch ivanpop/validation", func(t *testing.T) {
		var myDummy dummyStruct
		validationErr := myDummy.Validate()
		err := errors.Wrap(validationErr, "wraping validationErr cases")
		fiberApp := newFiberApp(err)
		request := httptest.NewRequest(fiber.MethodGet, "/", nil)

		response, err := fiberApp.Test(request)
		require.NoError(t, err)

		defer response.Body.Close()

		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		var respErr apperror.Error
		err = json.NewDecoder(response.Body).Decode(&respErr)
		require.NoError(t, err)
		require.Equal(t, constants.VALIDATION_ERROR, respErr.Message)
		require.Equal(t, map[string]any{
			"Field1": []any{"cannot be blank"},
		}, respErr.Errors)
	})

	// t.Run("should catch resourceful error", func(t *testing.T) {
	// 	customerHttpV1.NewCustomerInstance()
	// 	resourcefulErr := resourceful.
	// 		NewResource[uuid.UUID, *entity.Customer](customerHttpV1.CustomerDefinition).
	// 		SetParam(resourceful.Parameter{
	// 			Filters: []string{"created_on lte 1"},
	// 		})
	// 	err := errors.Wrap(resourcefulErr, "wraping resourcefulErr cases")
	// 	fiberApp := newFiberApp(err)
	// 	request := httptest.NewRequest(fiber.MethodGet, "/", nil)

	// 	response, err := fiberApp.Test(request)
	// 	require.NoError(t, err)

	// 	defer response.Body.Close()

	// 	require.Equal(t, http.StatusBadRequest, response.StatusCode)
	// 	var respErr httperror.Error
	// 	err = json.NewDecoder(response.Body).Decode(&respErr)
	// 	require.NoError(t, err)
	// 	require.Equal(t, constants.VALIDATION_ERROR, respErr.Message)
	// 	require.Equal(t, map[string]any{
	// 		"filters.1": []any{"value must be a valid date"},
	// 	}, respErr.Errors)
	// })

	// t.Run("should catch nullable validation error", func(t *testing.T) {
	// 	type dummyStruct struct {
	// 		Field1 nullable.NullString `json:"field1" nullable:"required"`
	// 	}

	// 	mydummyType := dummyStruct{
	// 		Field1: nullable.NullString{
	// 			IsExists: true,
	// 		},
	// 	}
	// 	validationErr := nullable.ValidateStruct(mydummyType)
	// 	err := errors.Wrap(validationErr, "wraping nullable validation cases")
	// 	fiberApp := newFiberApp(err)
	// 	request := httptest.NewRequest(fiber.MethodGet, "/", nil)

	// 	response, err := fiberApp.Test(request)
	// 	require.NoError(t, err)

	// 	defer response.Body.Close()

	// 	require.Equal(t, http.StatusBadRequest, response.StatusCode)
	// 	var respErr httperror.Error
	// 	err = json.NewDecoder(response.Body).Decode(&respErr)
	// 	require.NoError(t, err)
	// 	require.Equal(t, constants.VALIDATION_ERROR, respErr.Message)
	// 	require.Equal(t, map[string]any{
	// 		"field1": []any{"value cannot be null"},
	// 	}, respErr.Errors)
	// })

	t.Run("should catch jsonSyntaxErr", func(t *testing.T) {
		var dummyStruct struct {
			Field1 string `json:"field1"`
		}

		jsonErr := json.Unmarshal([]byte(`
		{
			"field1":12,
		}`), &dummyStruct)
		err := errors.Wrap(jsonErr, "wraping jsonErr cases")
		fiberApp := newFiberApp(err)
		request := httptest.NewRequest(fiber.MethodGet, "/", nil)

		response, err := fiberApp.Test(request)
		require.NoError(t, err)

		defer response.Body.Close()

		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		var respErr apperror.Error
		err = json.NewDecoder(response.Body).Decode(&respErr)
		require.NoError(t, err)
		require.Equal(t, constants.MALFORMED_BODY_ERROR, respErr.Message)
	})
	t.Run("should catch unmarshalError", func(t *testing.T) {
		var dummyStruct struct {
			Field1 string `json:"field1"`
		}

		unmarshalErr := json.Unmarshal([]byte(`
		{
			"field1":12
		}`), &dummyStruct)
		err := errors.Wrap(unmarshalErr, "wraping unmarshalErr cases")
		fiberApp := newFiberApp(err)
		request := httptest.NewRequest(fiber.MethodGet, "/", nil)

		response, err := fiberApp.Test(request)
		require.NoError(t, err)

		defer response.Body.Close()

		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		var respErr apperror.Error
		err = json.NewDecoder(response.Body).Decode(&respErr)
		require.NoError(t, err)
		require.Equal(t, constants.MALFORMED_BODY_ERROR, respErr.Message)
		require.Equal(t, map[string]any{
			"field1": []any{"the field must be a valid string"},
		}, respErr.Errors)
	})

	t.Run("should catch timeParseErr", func(t *testing.T) {
		var dummyStruct struct {
			Field1 optional.Time `json:"field1"`
		}

		timeParseError := json.Unmarshal([]byte(`
		{
			"field1":""
		}`), &dummyStruct)
		err := errors.Wrap(timeParseError, "wraping unmarshalErr cases")
		fiberApp := newFiberApp(err)
		request := httptest.NewRequest(fiber.MethodGet, "/", nil)

		response, err := fiberApp.Test(request)
		require.NoError(t, err)

		defer response.Body.Close()

		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		var respErr apperror.Error
		err = json.NewDecoder(response.Body).Decode(&respErr)
		require.NoError(t, err)
		require.Contains(t, respErr.Message, "invalid time format")
	})

	t.Run("should catch query parameterErr", func(t *testing.T) {
		newCustomFiberApp := func() *fiber.App {
			var fiberConfig fiber.Config
			fiberConfig.ErrorHandler = apperror.HttpHandleError
			fiberApp := fiber.New(fiberConfig)

			fiberApp.Get("/", func(c *fiber.Ctx) error {
				var dummyStruct struct {
					Field1 uuid.UUID `query:"uuid"`
				}
				err := c.QueryParser(&dummyStruct)
				if err != nil {
					return errors.Wrap(err, "wrap query paramErr")
				}

				return nil
			})

			return fiberApp
		}

		fiberApp := newCustomFiberApp()
		request := httptest.NewRequest(fiber.MethodGet, "/?uuid=123", nil)

		response, err := fiberApp.Test(request)
		require.NoError(t, err)

		defer response.Body.Close()

		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		var respErr apperror.Error
		err = json.NewDecoder(response.Body).Decode(&respErr)
		require.NoError(t, err)
		require.Equal(t, constants.MALFORMED_QUERY_ERROR, respErr.Message)
		require.Equal(t, map[string]any{
			"uuid": []any{"schema: error converting value for \"uuid\". Details: invalid UUID length: 3"}},
			respErr.Errors,
		)
	})

}

type dummyStruct struct {
	Field1 string
}

func (d dummyStruct) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Field1, validation.Required),
	)
}
