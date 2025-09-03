// cmd/web/main.go
package main

import (
	"fmt"
	"net/http"
)

func routes() http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		// Plain text OK; works for uptime/load balancers.
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Home
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, personal-site is running!")
	})

	return mux
}

func main() {
	addr := ":8080"
	fmt.Printf("Server listening on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, routes()); err != nil {
		panic(err)
	}
}
