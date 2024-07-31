package main

import (
	"backend/config"
	"backend/internal/app"
)

func main() {
	err := config.FillConfig()
	if err != nil {
		panic(err)
	}

	app.Run(config.Cfg)
}
