package main

import (
	"fmt"
	"net/http"
	"rate_limit"
)

func main() {
	rl := rate_limit.New(1)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Hello")
	})

	f := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !rl.Available() {
				http.NotFound(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	pipeline := f(mux)
	http.ListenAndServe(":8080", pipeline)
}
