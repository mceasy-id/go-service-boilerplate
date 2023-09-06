package optional

import (
	"database/sql"
	"database/sql/driver"
)

func NewString(values ...string) String {
	var s String
	s.SetEmpty()

	if len(values) > 0 {
		s.Set(values[0])
	}

	return s
}

type String struct {
	Option[string]
}

func (s String) Value() (driver.Value, error) {
	str, ok := s.Get()
	if !ok || !s.IsValueSet() {
		return nil, nil
	}
	return str, nil
}

func (s *String) Scan(value interface{}) error {
	sqlStr := sql.NullString{}
	err := sqlStr.Scan(value)
	if err != nil {
		return err
	}

	if sqlStr.Valid {
		s.Set(sqlStr.String)
	}
	return nil
}
