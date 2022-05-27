package demo

import (
	"feng/framework"
	"feng/framework/middleware"
	"time"
)

func Register(core *framework.Core) {
	// 静态路由+HTTP方法匹配
	// core.Get("/user/login", middleware.Test3(), UserLoginController)
	// 批量通用前缀
	core.GET("/test", middleware.Timeout(1*time.Second), middleware.Test2())
	subjectApi := core.Group("/subject")
	{
		subjectApi.Use(middleware.Test1())
		// 动态路由
		subjectApi.DELETE("/:id", SubjectDelController)
		subjectApi.PUT("/:id", SubjectUpdateController)
		subjectApi.GET("/:id", middleware.Test3(), SubjectGetController)
		subjectApi.GET("/list/all", SubjectListController)

		subjectInnerApi := subjectApi.Group("/info")
		{
			subjectInnerApi.GET("/name", SubjectNameController)
		}
	}
}
