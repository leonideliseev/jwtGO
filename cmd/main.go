package main

import "github.com/leonideliseev/jwtGO/internal/pkg/app"

func main() {
	ap := app.NewApp()

	ap.Run()
}
