// cmd/web/main.go
package main

import (
	"net/http"

	"github.com/brandondunbar/personal-site/internal/config"
)

func main() {
	rt := config.LoadRuntime()

	println("Server environment:", rt.Env)
	println("Server listening on", rt.BaseURL)

	app, err := NewApp()
	if err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(rt.Addr, app.Routes()); err != nil {
		panic(err)
	}
}

