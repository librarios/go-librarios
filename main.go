package main

import (
	"github.com/librarios/go-librarios/app"
)

func main() {
	filename := "config/librarios.yaml"
	app.StartServer(filename)
}
