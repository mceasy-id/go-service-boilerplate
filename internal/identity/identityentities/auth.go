package identityentities

import "github.com/invopop/validation"

const (
	KeyAuthCredential = "authCredential"
)

type Credential struct {
	UserName  string `json:"name,omitempty"`
	UserId    uint64 `json:"userId,omitempty"`
	CompanyId uint64 `json:"companyId,omitempty"`
}

func (c Credential) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.UserName, validation.Required),
		validation.Field(&c.UserId, validation.Required),
		validation.Field(&c.CompanyId, validation.Required),
	)
}
