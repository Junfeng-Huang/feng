package demo

import (
	"feng/app/provider/demo"
	"feng/framework"
	"feng/framework/contract"
)

func UserLoginController(c *framework.Context) {
	// test
	DemoService := c.MustMake(demo.Key).(demo.Service)
	f := DemoService.GetFoo()
	// 输出结果
	c.SetOkStatus().Json(f)
	AppService := c.MustMake(contract.AppKey).(contract.App)
	c.Json("\n" + AppService.BaseFolder())
}
