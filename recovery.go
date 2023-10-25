package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				fmt.Println("*** Detected a panic in the HTTP handler ***")
				debug.PrintStack()

				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "internal server error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func SimpleRecovery() {
	if rvr := recover(); rvr != nil {
		fmt.Println("*** Gracefully handling a panic ***")
		debug.PrintStack()
	}
}

type PanicHelper interface {
	PanicFn() string
}

func TriggerPanic(w http.ResponseWriter, r *http.Request) {
	var panicHelper PanicHelper
	fmt.Println(panicHelper.PanicFn())
}

func TriggerCrash(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "I will crash the app"}`))

	go func() {
		var panicHelper PanicHelper
		fmt.Println(panicHelper.PanicFn())
	}()
}

func TriggerCrashWithRecovery(w http.ResponseWriter, r *http.Request) {
	go func() {
		defer SimpleRecovery()
		var panicHelper PanicHelper
		fmt.Println(panicHelper.PanicFn())
	}()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "I will recover from a panic"}`))
}
