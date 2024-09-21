package http

import (
	"encoding/json"
	"io"
	"net/http"
)

var err error

// Reads the entirety of the given request's body and unmarshalls it into
// the given pointer to the JSON schema
func ReadUnmarshalRequestBody(request *http.Request, schema any) error {
	var requestBodyBytes []byte
	requestBodyBytes, err = io.ReadAll(request.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(requestBodyBytes, schema)

	if err != nil {
		return err
	}

	return nil
}
