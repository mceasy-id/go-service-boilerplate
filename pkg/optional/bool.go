package optional

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

func NewBool(vs ...bool) Bool {
	var b Bool
	b.SetEmpty()

	if len(vs) > 0 {
		b.Set(vs[0])
	}

	return b
}

type Bool struct {
	value bool
	isSet bool
}

func (b *Bool) IsValueSet() bool {
	return b.isSet
}

func (b *Bool) Set(v bool) {
	b.value = v
	b.isSet = true
}
func (b *Bool) SetEmpty() {
	b.value = false
	b.isSet = true
}

func (b Bool) Get() (bool, bool) {
	if !b.IsPresent() {
		var zero bool
		return zero, false
	}

	return b.value, true
}
func (o Bool) MustGet() bool {
	if !o.IsPresent() {
		panic("value is not present")
	}

	return o.value
}

func (o Bool) IsPresent() bool {
	return true
}

func (t Bool) IfPresent(fn func(bool)) Bool {
	if t.IsPresent() {
		fn(t.value)
	}
	return t
}

func (o Bool) MarshalJSON() ([]byte, error) {
	if o.IsPresent() {
		return json.Marshal(o.value)
	}

	return json.Marshal(nil)
}

func (o *Bool) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.value = false
		o.isSet = true
		return nil
	}
	var value bool
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.value = value
	o.isSet = true
	return nil
}

func (b Bool) Value() (driver.Value, error) {
	v, ok := b.Get()
	if !ok || !b.IsValueSet() {
		return nil, nil
	}
	return v, nil
}

func (b *Bool) Scan(value interface{}) error {
	sqlBool := sql.NullBool{}
	err := sqlBool.Scan(value)
	if err != nil {
		return err
	}

	if sqlBool.Valid {
		b.Set(sqlBool.Bool)
	}
	return nil
}
