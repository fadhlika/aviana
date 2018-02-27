package main

import (
	"github.com/fadhlika/aviana/app"
)

func main() {
	app := app.App{}
	app.Init()
	app.Run()
}
