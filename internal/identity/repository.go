package identity

import (
	"context"

	"mceasy/service-demo/internal/identity/identityentities"
)

type Repository interface {
	FindCompanyProfile(ctx context.Context, companyId uint64) (*identityentities.CompanyProfile, error)
	FindDriverById(ctx context.Context, id uint64) (*identityentities.User, error)
}
