package optional

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

func NewTime(values ...time.Time) Time {
	var t Time
	t.SetEmpty()

	if len(values) > 0 {
		t.Set(values[0])
	}

	return t
}

type Time struct {
	Option[time.Time]
}

func (t *Time) Scan(value any) error {
	var sqlTime sql.NullTime
	err := sqlTime.Scan(value)
	if err != nil {
		return err
	}

	if sqlTime.Valid {
		t.Set(sqlTime.Time)
	}

	return nil
}

func (t Time) Value() (driver.Value, error) {
	timeValue, ok := t.Get()
	if !ok {
		return nil, nil
	}
	return timeValue, nil
}
