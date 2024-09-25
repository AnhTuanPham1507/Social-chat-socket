package utils

import (
	"encoding/json"
)

func ParseByteArray(r []byte, x interface{}) error {
	if err := json.Unmarshal(r, x); err != nil {
		return err
	}
	return nil
}