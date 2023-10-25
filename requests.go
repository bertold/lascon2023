package main

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
)

func OnlyPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"error": "method not allowed"}`))
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "got your POST"}`))
	}
}

func OnlyJSON(w http.ResponseWriter, r *http.Request) {
	// Content-Type may use different casing, and may contain additional
	// information as charset. For example:
	// Content-Type: application/json; charset=utf-8
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	w.Header().Add("Content-Type", "application/json")

	if err != nil || contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		_, _ = w.Write([]byte(`{"error": "JSON content is expected, but received: ` + contentType + `"}`))
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "thanks for sending JSON"}`))
	}
}

const maxRequestBodySize = 1024

func BodyTooLarge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	defer r.Body.Close()

	content := make([]byte, maxRequestBodySize+1)
	readBytes, err := r.Body.Read(content)

	w.Header().Add("Content-Type", "application/json")
	if err == nil || errors.Is(err, io.EOF) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"message": "got your POST of length: %d"}`, readBytes)))
	} else {
		if _, ok := err.(*http.MaxBytesError); ok {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			_, _ = w.Write([]byte(`{"message": "request body is too large"}`))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error": "failed to parse request body: ` + err.Error() + `"}`))
		}
	}
}
