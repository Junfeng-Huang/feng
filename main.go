package main

import (
	"feng/app/appHttp"
	"feng/app/provider/demo"
	"feng/framework"
	"feng/framework/middleware"
)

func main() {
	core := framework.Default()
	core.Bind(&demo.DemoServiceProvider{})

	core.Use(middleware.Cost())

	appHttp.Routes(core)
	core.Run("0.0.0.0:8888")
	// command.RunCommand()
}
