package resourceful

import (
	"strconv"
	"strings"
	"time"
)

const (
	NUMERIC = "numeric"
	BOOLEAN = "boolean"
	STRING  = "string"
	DATE    = "date"
)

func sanitizeFilterValueByType(value, filterType string) ([]string, bool) {
	value, _ = strings.CutPrefix(value, "(")
	value, _ = strings.CutSuffix(value, ")")

	value = " " + value + " "
	lo := -1
	var quoted byte
	separatedValues := make([]string, 0)
	for hi := 0; hi < len(value); hi++ {
		if quoted == value[hi] {
			if !((quoted == '"' || quoted == '\'') && hi+1 < len(value) && value[hi+1] != ' ') || quoted == ' ' {
				separatedValues = append(separatedValues, value[lo:hi])
				quoted = 0
			}

			if hi+1 < len(value) && value[hi+1] == ' ' {
				hi++
			}
		}
		if (value[hi] == '"' || value[hi] == '\'') && quoted == 0 && hi+1 < len(value) {
			quoted = value[hi]
			lo = hi + 1
		}

		if value[hi] == ' ' && hi+1 < len(value) && value[hi+1] != '"' && value[hi+1] != '\'' && quoted == 0 {
			quoted = value[hi]
			lo = hi + 1
		}

	}

	if len(separatedValues) == 0 {
		return nil, false
	}

	sanitizedValues := make([]string, 0, len(separatedValues))
	for _, separatedValue := range separatedValues {
		separatedValue := strings.TrimSpace(separatedValue)

		if separatedValue == "" || separatedValue == "()" || separatedValue == "(" || separatedValue == ")" {
			return nil, false
		}

		switch filterType {
		case NUMERIC:
			_, err := strconv.ParseFloat(separatedValue, 64)
			if err != nil {
				return nil, false
			}

		case BOOLEAN:
			if separatedValue != "true" && separatedValue != "false" {
				return nil, false
			}
		case DATE:
			_, err := time.Parse(time.RFC3339, separatedValue)
			if err != nil {
				return nil, false
			}
		}

		sanitizedValues = append(sanitizedValues, separatedValue)
	}

	return sanitizedValues, true
}
