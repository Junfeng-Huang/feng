package main

import (
	"feng/app/appHttp"
	"feng/app/provider/demo"
	"feng/framework"
)

func main() {
	core := framework.Default()
	core.Bind(&demo.DemoServiceProvider{})
	appHttp.Routes(core)
	core.Run("0.0.0.0:8888")

}
