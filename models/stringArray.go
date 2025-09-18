package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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
