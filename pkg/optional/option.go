package optional

// https://github.com/markphelps/optional MIT License
import (
	"encoding/json"
)

type Option[T comparable] struct {
	value *T
	isSet bool
}

// IsValueSet check is value has been set programmaticly
// or set from external resource like form or json
func (o *Option[T]) IsValueSet() bool {
	return o.isSet
}

// Set set the value
func (o *Option[T]) Set(v T) Option[T] {
	o.value = &v
	o.isSet = true
	return *o
}

// SetEmpty set the value to nil pointer / default value
func (o *Option[T]) SetEmpty() Option[T] {
	var value *T
	o.value = value
	o.isSet = true
	return *o
}

// Get if its nil, it return default value and false
// if its not nil, it return value and true
//
// example:
//
//	func main() {
//		var stringOption optional.String
//		json.UnmarshalJSON(any, &stringOption)
//		v, ok := stringOption.Get()
//		if !ok {
//			//Do error handling
//		}
//		newValue := v
//	}
func (o Option[T]) Get() (T, bool) {
	if !o.IsPresent() {
		var zero T
		return zero, false
	}

	return *o.value, true
}

func (o Option[T]) GetOrDefault(defaultValue T) T {
	if !o.IsPresent() {
		return defaultValue
	}

	return *o.value
}

// MustGet same with Get but it will panic when value is nil
func (o Option[T]) MustGet() T {
	if !o.IsPresent() {
		panic("value is not present")
	}

	return *o.value
}

// IsPresent check if value nil or not
func (o Option[T]) IsPresent() bool {
	return o.value != nil
}

// IfPresent give you a callback value if value is not nil
//
// example:
//
//	func main() {
//		var stringOption optional.String
//		json.UnmarshalJSON(any, &stringOption)
//		stringOption.IfPresent(func(s string) {
//			//set to another entity
//		})
//	}
func (t Option[T]) IfPresent(fn func(T)) Option[T] {
	if t.IsPresent() {
		fn(*t.value)
	}
	return t
}

// MarshalJSON If value is nil, we  marshal it to nil
// if value is not nil, we marshal it to its value
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.IsPresent() {
		return json.Marshal(o.value)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON to check if the value is set to nil
// or the key is not sent
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.SetEmpty()
		return nil
	}
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.value = &value
	o.isSet = true
	return nil
}
