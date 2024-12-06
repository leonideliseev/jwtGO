package main

import (
	"log"

	"github.com/leonideliseev/jwtGO/internal/pkg/app"
)

func main() {
	ap, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	err = ap.Run()
	if err != nil {
		log.Fatal(err)
	}
}
