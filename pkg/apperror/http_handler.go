package apperror

import (
	// "encoding/json"

	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"mceasy/service-demo/constants"

	"mceasy/service-demo/pkg/observability"
	"mceasy/service-demo/pkg/observability/instrumentation"
	"mceasy/service-demo/pkg/resourceful"

	"github.com/golang-jwt/jwt/v4"
	"github.com/invopop/validation"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

var regexSnakeCase = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

type Error struct {
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func toSnakeCase(s string) string {
	return strings.ToLower(strings.Join(regexSnakeCase.FindAllString(s, -1), "_"))
}

func HttpHandleError(c *fiber.Ctx, err error) error {
	// Unauthorized Error
	var jwtErr *jwt.ValidationError
	if errors.As(err, &jwtErr) || err.Error() == "Missing or malformed JWT" {
		return c.Status(http.StatusUnauthorized).JSON(Error{
			Message: constants.UNAUTHORIZED_ERROR,
		})
	}

	// Path Parse Error
	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: constants.MALFORMED_BODY_ERROR,
		})
	}

	// Handle Http Error
	var appErr *AppError
	if errors.As(err, &appErr) {
		if errors.Is(appErr.Err, ErrBadRequest) {
			if appErr.errMap != nil {
				return c.Status(http.StatusBadRequest).JSON(Error{
					Message: appErr.Message,
					Errors:  appErr.errMap,
				})
			}

			return c.Status(http.StatusBadRequest).JSON(Error{
				Message: appErr.Message,
			})
		}

		if errors.Is(appErr.Err, ErrUnauthorized) {
			return c.SendStatus(http.StatusUnauthorized)
		}

		if errors.Is(appErr.Err, ErrForbiddenAccess) {
			if appErr.errMap != nil {
				return c.Status(http.StatusForbidden).JSON(Error{
					Message: appErr.Message,
					Errors:  appErr.errMap,
				})
			}

			return c.Status(http.StatusForbidden).JSON(Error{
				Message: appErr.Message,
			})
		}

		if errors.Is(appErr.Err, ErrNotFound) {
			return c.SendStatus(http.StatusNotFound)
		}

		if errors.Is(appErr.Err, ErrConflict) {
			return c.SendStatus(http.StatusConflict)
		}

		observability.SendErrorToTeams(c, err)
		if errors.Is(appErr.Err, ErrGateway) {
			return c.SendStatus(http.StatusGatewayTimeout)
		}

		return c.SendStatus(http.StatusInternalServerError)

	}
	var validatorError validation.Errors
	if errors.As(err, &validatorError) {
		mapErr := validationErrorMapping(validatorError)
		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: constants.VALIDATION_ERROR,
			Errors:  mapErr,
		})
	}

	// Validation goresourceful error
	var resourcefulErr resourceful.ValidationErrors
	if errors.As(err, &resourcefulErr) {
		validationErrors := make(map[string][]string)
		for _, val := range resourcefulErr {
			validationErrors[val.FieldName] = val.Errors
		}

		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: constants.VALIDATION_ERROR,
			Errors:  validationErrors,
		})
	}

	// JSON Format Error
	var jsonSyntaxErr *json.SyntaxError
	if errors.As(err, &jsonSyntaxErr) {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: constants.MALFORMED_BODY_ERROR,
		})
	}

	// Unmarshal Error
	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalErr) {
		var translatedType string
		switch unmarshalErr.Type.Name() {
		// REGEX *int*
		case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
			translatedType = "number"
		case "Time":
			translatedType = "date time"
		case "string":
			translatedType = "string"
		}

		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: constants.MALFORMED_BODY_ERROR,
			Errors: map[string][]string{
				unmarshalErr.Field: {fmt.Sprintf("the field must be a valid %s", translatedType)},
			},
		})
	}

	//time parse error
	var timeParseErr *time.ParseError
	if errors.As(err, &timeParseErr) {
		v := timeParseErr.Value
		if v == "" {
			v = "empty string (``)"
		}
		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: fmt.Sprintf("invalid time format on %s", v),
		})
	}

	// Query Parameter Error
	var fiberMultiErr fiber.MultiError
	if errors.As(err, &fiberMultiErr) {
		validationErrors := make(map[string][]string)

		for key, err := range fiberMultiErr {
			validationErrors[key] = append(validationErrors[key], err.Error())
		}
		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: constants.MALFORMED_QUERY_ERROR,
			Errors:  validationErrors,
		})
	}

	// Multipart Error
	if errors.Is(err, fasthttp.ErrNoMultipartForm) {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: "invalid multipart content-type",
		})
	}

	// Default Fiber Error
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return c.SendStatus(fiberErr.Code)
	}

	// Internal Server Error
	// Debug Mode for local env
	if os.Getenv("APP_ENV") == "local" {
		log.Println(err)
		log.Println(string(debug.Stack()))
	}

	if os.Getenv("TEAMS_WEBHOOK_URI") != "" {
		observability.SendErrorToTeams(c, err)
	}

	//TCP connection error
	var tcpErr *net.OpError
	if errors.As(err, &tcpErr) {
		log.Fatalf("unable to get tcp connection from %s, shutting down...", tcpErr.Addr.String())
	}

	span := trace.SpanFromContext(c.Context())
	instrumentation.RecordSpanError(span, err)
	return c.SendStatus(http.StatusInternalServerError)
}

func validationErrorMapping(validatorError validation.Errors) map[string][]string {
	mapErr := make(map[string][]string)
	for key, err := range validatorError {
		if errs, ok := err.(validation.Errors); ok {
			newMap := validationErrorMapping(errs)
			mapErr = mergeMapWithKey(key, mapErr, newMap)
		} else {
			mapErr[key] = append(mapErr[key], err.Error())
		}

	}
	return mapErr
}

func mergeMapWithKey(key string, maps ...map[string][]string) map[string][]string {
	res := make(map[string][]string)
	for _, m := range maps {
		for k, v := range m {
			mergedKey := key + "." + k
			res[mergedKey] = append(res[mergedKey], v...)
		}
	}
	if len(res) == 0 {
		return nil
	}
	return res
}
