package identityentities

type CompanyProfileHeadOfficeAddress struct {
	FullName string `json:"full_name"`
}

type CompanyProfile struct {
	Name              string                          `json:"name"`
	LogoURL           string                          `json:"logo_url"`
	HeadOfficeAddress CompanyProfileHeadOfficeAddress `json:"head_office_address"`
	Email             string                          `json:"email"`
	PhoneNumber       string                          `json:"phone_number"`
	FaxNumber         string                          `json:"fax_number"`
}
