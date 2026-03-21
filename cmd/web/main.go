package main

import (
	"github.com/kingsukhoi/ipv6-tester/pkg/router"
)

func main() {
	e := router.NewRouter()

	err := e.Start(":1323")

	if err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
