package httprepository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"mceasy/service-demo/config"
	"mceasy/service-demo/internal/identity"
	"mceasy/service-demo/internal/identity/identityentities"
	"mceasy/service-demo/pkg/apperror"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
)

func New(cfg config.Config, client *retryablehttp.Client) identity.Repository {
	return &httpRepo{
		cfg:  cfg,
		conn: client,
	}
}

type httpRepo struct {
	cfg  config.Config
	conn *retryablehttp.Client
}

// FindUserById implements masterdata.Repository.
func (repo *httpRepo) FindDriverById(ctx context.Context, id uint64) (*identityentities.User, error) {
	req, err := http.NewRequest(http.MethodGet, repo.cfg.ExternalURI.MasterData.User+fmt.Sprintf("/%d", id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", repo.cfg.ExternalURI.MasterData.GetBearerToken())
	otelhttptrace.Inject(ctx, req)

	resp, err := repo.conn.Do(&retryablehttp.Request{
		Request: req,
	})
	if err != nil {
		if strings.Contains(err.Error(), "giving up after") {
			return nil, errors.Wrap(apperror.GatewayTimeout(), "HttpClient.FindUserById.GatewayTimeout")
		}
		return nil, errors.Wrap(err, "HttpClient.FindUserById.GetHttp")
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		var errorData struct {
			Message string `json:"message"`
		}
		err = json.NewDecoder(resp.Body).Decode(&errorData)
		if err != nil {
			return nil, errors.Wrap(err, "HttpClient.FindUserById.JSONDecoder")
		}
		return nil, errors.Wrap(errors.New(errorData.Message), "HttpClient.FindUserById.Response")
	}

	var data struct {
		User identityentities.User `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, errors.Wrap(err, "HttpClient.FindUserById.JSONDecoder")
	}

	if data.User.Position.Name != identityentities.DriverPositionName {
		return nil, nil
	}

	return &data.User, nil
}
