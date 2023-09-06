package identityentities

type UserPosition struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}
type User struct {
	Id        uint64       `json:"id"`
	CompanyId uint64       `json:"company_id"`
	Name      string       `json:"name"`
	Position  UserPosition `json:"position"`
}

