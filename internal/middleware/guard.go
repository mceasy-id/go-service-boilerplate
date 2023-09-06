package middleware

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"mceasy/service-demo/config"
	"mceasy/service-demo/internal/identity/identityentities"
	"mceasy/service-demo/pkg/apperror"
	"mceasy/service-demo/pkg/httpclient"
	"mceasy/service-demo/pkg/observability/instrumentation"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
)

func GuardMiddleware(cfg config.Config) func(*fiber.Ctx) error {
	client := httpclient.NewWithoutLog()

	return func(c *fiber.Ctx) error {
		_, span := instrumentation.NewTraceSpan(
			c.UserContext(),
			"GuardMiddleware.IdentityCheck",
		)

		defer span.End()

		var (
			err   error
			token string
		)

		token = c.Cookies("resource")
		if tokenHeader, err := getTokenFromHeader(c); err == nil && tokenHeader != "" {
			token = tokenHeader
		}

		if token != "" {
			err = validateJWT(token, cfg.Authentication.Key, c)
			if err == nil {
				span.End()
				return c.Next()
			}

			err = validateOauth2(cfg, c, token, client)
			if err == nil {
				span.End()
				return c.Next()
			}
		}
		if err == nil && token == "" {
			err = apperror.Unauthorized()
		}
		return apperror.HttpHandleError(c, errors.Wrap(err, "GuardMiddleware"))
	}
}

func validateOauth2(cfg config.Config, ctx *fiber.Ctx, token string, client *retryablehttp.Client) error {
	err := checkIdentityScope(cfg, token, client)
	if err != nil {
		return err
	}

	err = getIdentityAuthorization(cfg, ctx, token, client)
	if err != nil {
		return err
	}
	return nil
}

func validateJWT(token string, key string, c *fiber.Ctx) error {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return err
	}

	if parsedToken.Valid {
		mapClaim := parsedToken.Claims.(jwt.MapClaims)
		var authCredential identityentities.Credential
		authCredential.CompanyId = uint64(mapClaim["companyId"].(float64))
		authCredential.UserId = uint64(mapClaim["userId"].(float64))
		authCredential.UserName = mapClaim["name"].(string)

		if err := authCredential.Validate(); err != nil {
			return apperror.Unauthorized()
		}
		c.Locals(identityentities.KeyAuthCredential, authCredential)
		return nil
	}

	return err

}

func checkIdentityScope(cfg config.Config, token string, client *retryablehttp.Client) error {
	req, err := http.NewRequest(http.MethodGet, cfg.ExternalURI.Scope, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(&retryablehttp.Request{
		Request: req,
	})
	if err != nil {
		if strings.Contains(err.Error(), "giving up after") {
			return apperror.GatewayTimeout()
		}
		return err
	}
	defer resp.Body.Close()

	var payload struct {
		Data []string `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		return err
	}

	scopingOrder := map[string]bool{
		"tms:create": true,
		"tms:delete": true,
		"tms:read":   true,
		"tms:update": true,
	}

	counter := 0
	for _, scope := range payload.Data {
		if _, ok := scopingOrder[scope]; ok {
			counter++
		}
	}
	if counter != len(scopingOrder) {
		return apperror.Unauthorized()
	}

	return nil
}

func getIdentityAuthorization(cfg config.Config, ctx *fiber.Ctx, token string, client *retryablehttp.Client) error {
	parsedUrl, err := url.Parse(cfg.ExternalURI.Identity)
	if err != nil {
		return err
	}

	queryString := parsedUrl.Query()
	queryString.Set("kind", "POST")
	queryString.Set("action", "/api/oauth/token")
	parsedUrl.RawQuery = queryString.Encode()

	req, err := http.NewRequest(http.MethodGet, parsedUrl.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("namespace", "identity")

	resp, err := client.Do(&retryablehttp.Request{
		Request: req,
	})
	if err != nil {
		if strings.Contains(err.Error(), "giving up after") {
			return apperror.GatewayTimeout()
		}
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return apperror.NotFound()
	}

	if resp.StatusCode == http.StatusForbidden {
		return apperror.Forbidden()
	}

	var payload struct {
		Data struct {
			AuthorizedResource *struct {
				Resource identityentities.Credential `json:"resource"`
			} `json:"authorized_resource"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		return err
	}

	if payload.Data.AuthorizedResource == nil {
		return apperror.Unauthorized()
	}

	if err := payload.Data.AuthorizedResource.Resource.Validate(); err != nil {
		return apperror.Unauthorized()
	}

	ctx.Locals(identityentities.KeyAuthCredential, payload.Data.AuthorizedResource.Resource)
	return nil
}

func getTokenFromHeader(c *fiber.Ctx) (string, error) {
	auth := c.Get("Authorization")
	if auth == "" {
		return "", nil
	}
	parts := strings.Split(auth, " ")
	if len(parts) < 2 {
		return "", apperror.Unauthorized()
	}
	if parts[0] != "Bearer" {
		return "", apperror.Forbidden()
	}

	return parts[1], nil
}
