package optional

import (
	"database/sql"
	"database/sql/driver"
)

func NewFloat32(values ...float32) Float32 {
	var f Float32
	f.SetEmpty()

	if len(values) > 0 {
		f.Set(values[0])
	}

	return f
}

type Float32 struct {
	Option[float32]
}

func (f Float32) Value() (driver.Value, error) {
	v, ok := f.Get()
	if !f.IsValueSet() || !ok {
		return nil, nil
	}

	return v, nil
}

func (f *Float32) Scan(value interface{}) error {
	sqlFloat64 := sql.NullFloat64{}

	err := sqlFloat64.Scan(value)
	if err != nil {
		return err
	}

	if sqlFloat64.Valid {
		f.Set(float32(sqlFloat64.Float64))
	}

	return nil
}

func NewFloat64(values ...float64) Float64 {
	var f Float64
	f.SetEmpty()

	if len(values) > 0 {
		f.Set(values[0])
	}

	return f
}

type Float64 struct {
	Option[float64]
}

func (f Float64) IsValueSet() bool {
	return f.Option.IsValueSet()
}

func (f Float64) Value() (driver.Value, error) {
	v, ok := f.Get()
	if !f.IsValueSet() || !ok {
		return nil, nil
	}

	return v, nil
}

func (f *Float64) Scan(value interface{}) error {
	sqlFloat64 := sql.NullFloat64{}

	err := sqlFloat64.Scan(value)
	if err != nil {
		return err
	}

	if sqlFloat64.Valid {
		f.Set(sqlFloat64.Float64)
	}

	return nil
}
