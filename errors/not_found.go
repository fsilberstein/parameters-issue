package errors

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"
)

type errNotFound struct {
	error
}

// NewNotFoundError creates a special error that, when processed by GoKit's DefaultErrorEncoder, will translate to a
// fully-fledged ReST HTTP response:
// - HTTP status code 400
// - JSON response body like '{ "error" : "Not found: user" }'
func NewNotFoundError(msg string) error {
	return errNotFound{stderrors.New(fmt.Sprintf("Not found: %s", msg))}
}

// MarshalJSON lets GoKit's DefaultErrorEncoder set the proper HTTP response body
func (e errNotFound) MarshalJSON() ([]byte, error) {
	outputBody := map[string]interface{}{}
	outputBody["error"] = e.error.Error()
	return json.Marshal(outputBody)
}

// StatusCode lets GoKit's DefaultErrorEncoder set the proper HTTP status code in response
func (e errNotFound) StatusCode() int {
	return http.StatusNotFound
}
