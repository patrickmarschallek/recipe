package recipe

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// Valid interface implements a validation function for a struct.
type Valid interface {
	Valid() error
}

// Validate calls the implemented Valid interface on the given struct
func Validate(v interface{}) error {
	obj, ok := v.(Valid)
	if !ok {
		return nil // no valid method
	}
	return obj.Valid()
}

// ValidateRequestLength checks the requests bodies content length.
func ValidateRequestLength(r *http.Request) error {
	if r.ContentLength == 0 {
		return NewErrorResponse("empty payload", http.StatusUnprocessableEntity)
	}
	return nil
}

// JSONDecodeAndValidate - entrypoint for deserialization and validation
// of the submission input
func JSONDecodeAndValidate(r *http.Request, v interface{}) error {
	if err := ValidateRequestLength(r); err != nil {
		return NewErrorResponse("empty payload", http.StatusUnprocessableEntity)
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	defer r.Body.Close()
	return Validate(v)
}

// XMLDecodeAndValidate - entrypoint for deserialization and validation
// of the submission input
func XMLDecodeAndValidate(r *http.Request, v interface{}) error {
	if err := xml.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	defer r.Body.Close()
	return Validate(v)
}
