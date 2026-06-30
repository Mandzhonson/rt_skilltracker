package main

import (
	"core_service/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
