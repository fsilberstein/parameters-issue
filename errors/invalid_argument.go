package errors

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"
)

type errInvalidArgument struct {
	error
}

// NewInvalidArgument creates a special error that, when processed by GoKit's DefaultErrorEncoder, will translate to a
// fully-fledged ReST HTTP response:
// - HTTP status code 400
// - JSON response body like '{ "error" : "Invalid argument: amount cannot be negative" }'
func NewInvalidArgument(msg string) error {
	return errInvalidArgument{stderrors.New(fmt.Sprintf("Invalid argument: %s", msg))}
}

// MarshalJSON lets GoKit's DefaultErrorEncoder set the proper HTTP response body
func (e errInvalidArgument) MarshalJSON() ([]byte, error) {
	outputBody := map[string]interface{}{}
	outputBody["error"] = e.error.Error()
	return json.Marshal(outputBody)
}

// StatusCode lets GoKit's DefaultErrorEncoder set the proper HTTP status code in response
func (e errInvalidArgument) StatusCode() int {
	return http.StatusBadRequest
}
