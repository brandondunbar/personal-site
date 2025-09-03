// cmd/web/main.go
package main

import (
	"fmt"
	"net/http"
)

func routes() http.Handler {
	mux := http.NewServeMux()
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
