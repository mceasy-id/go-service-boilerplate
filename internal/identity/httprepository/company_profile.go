package httprepository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"mceasy/service-demo/internal/identity/identityentities"
	"mceasy/service-demo/pkg/apperror"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
)

func (repo *httpRepo) FindCompanyProfile(ctx context.Context, companyId uint64) (*identityentities.CompanyProfile, error) {
	req, err := http.NewRequest(http.MethodGet, repo.cfg.ExternalURI.MasterData.CompanyProfile+fmt.Sprintf("/%d", companyId), nil)
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
			return nil, errors.Wrap(apperror.GatewayTimeout(), "HttpClient.GetCompanyProfile.GatewayTimeout")
		}
		return nil, errors.Wrap(err, "HttpClient.GetCompanyProfile.GetHttp")
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
			return nil, errors.Wrap(err, "HttpClient.GetCompanyProfile.JSONDecoder")
		}
		return nil, errors.Wrap(errors.New(errorData.Message), "HttpClient.GetCompanyProfile.Response")
	}

	var data struct {
		CompanyProfile identityentities.CompanyProfile `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, errors.Wrap(err, "HttpClient.GetCompanyProfile.JSONDecoder")
	}

	return &data.CompanyProfile, nil
}
