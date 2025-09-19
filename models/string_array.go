package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type StringArray []string

func (sa StringArray) Value() (driver.Value, error) {
	if len(sa) == 0 {
		return nil, nil
	}

	return json.Marshal(sa)
}

func (sa *StringArray) Scan(src any) error {
	if src == nil {
		*sa = nil
		return nil
	}

	bytes, ok := src.([]byte)
	if !ok {
		return errors.New("[Error] src cannot be case to []byte")
	}

	return json.Unmarshal(bytes, sa)
}

func (sa StringArray) ValidateUniqueItems() validation.Rule {
	return validation.By(func(value interface{}) error {
		seen := make(map[string]bool, len(sa))
		for _, v := range sa {
			if _, exists := seen[v]; exists {
				return errors.New("contains duplicate value " + v)
			}
			seen[v] = true
		}
		return nil
	})
}
