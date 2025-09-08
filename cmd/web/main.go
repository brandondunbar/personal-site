// cmd/web/main.go
package main

import (
	"net/http"
)

func main() {
	addr := ":8080"
	println("Server listening on http://localhost" + addr)

	app, err := NewApp()
	if err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(addr, app.Routes()); err != nil {
		panic(err)
	}
}

