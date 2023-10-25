package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type MyRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var (
	ErrUnmarshalType    = &json.UnmarshalTypeError{}
	ErrSyntax           = &json.SyntaxError{}
	ErrInvalidUnmarshal = &json.InvalidUnmarshalError{}
)

func BasicJsonRequest(w http.ResponseWriter, r *http.Request) {
	var myRequest MyRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	w.Header().Add("Content-Type", "application/json")

	err := decoder.Decode(&myRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		switch {
		case errors.Is(err, ErrUnmarshalType),
			errors.Is(err, ErrSyntax),
			errors.Is(err, ErrInvalidUnmarshal),
			errors.Is(err, io.ErrUnexpectedEOF),
			errors.Is(err, io.EOF):
			w.Write([]byte(`{"error": "invalid JSON"}`))
		// unfortunately, this is the easiest way to check for an unknown field
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			w.Write([]byte(`{"error": "field is not declared in schema"}`))
		default:
			w.Write([]byte(`{"error": "failed to parse request"}`))
		}
		return
	}

	// check to see if multiple JSON objects are present
	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "unexpected data in request"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "Hello, ` +
		myRequest.FirstName + ` ` + myRequest.LastName + `"}`))
}

var schema = gojsonschema.NewReferenceLoader("file:///Users/bkolics/src/lascon2023/requestSchema.json")

func ValidatedJsonRequest(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	defer r.Body.Close()

	w.Header().Add("Content-Type", "application/json")

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "unexpected error parsing request"}`))
		return
	}

	loader := gojsonschema.NewBytesLoader(requestBody)
	result, err := gojsonschema.Validate(schema, loader)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid JSON request ` + err.Error() + `"}`))
		return
	}
	if !result.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "malformed JSON"}`))

		fmt.Fprintln(os.Stderr, "Validation errors")
		for _, desc := range result.Errors() {
			fmt.Fprintf(os.Stderr, "- %s\n", desc)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "thank you for your JSON"}`))
}
