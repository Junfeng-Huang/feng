package appHttp

import (
	demo "feng/app/appHttp/module/demo/control"
	"feng/framework"
)

// 注册路由规则
func Routes(core *framework.Core) {
	demo.Register(core)
}
