package dockercloud

import (
	"fmt"
)

type HttpError struct {
	Status     string
	StatusCode int
	Message    []byte
}

func (e HttpError) Error() string {
	if Debug == true {
		return fmt.Sprintf("Failed API call: %s Message: %s", e.Status, e.Message)
	}

	return fmt.Sprintf("Failed API call: %s", e.Status)
}
