package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	basicHandler := http.HandlerFunc(TriggerPanic)
	mux.Handle("/panic/plain", basicHandler)
	mux.Handle("/panic/recover", RecoveryMiddleware(basicHandler))
	mux.Handle("/panic/crash", http.HandlerFunc(TriggerCrash))
	mux.Handle("/panic/crashwithrecovery", http.HandlerFunc(TriggerCrashWithRecovery))
	mux.Handle("/headers/posts", http.HandlerFunc(OnlyPosts))
	mux.Handle("/headers/json", http.HandlerFunc(OnlyJSON))
	mux.Handle("/body/toolarge", http.HandlerFunc(BodyTooLarge))
	mux.Handle("/json/basic", http.HandlerFunc(BasicJsonRequest))
	mux.Handle("/json/validate", http.HandlerFunc(ValidatedJsonRequest))
	http.ListenAndServe(":8080", mux)
}
