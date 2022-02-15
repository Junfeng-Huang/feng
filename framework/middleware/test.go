package middleware

import (
	"feng/framework"
	"fmt"
)

func Test1() framework.ControllerHandler {
	// 使用函数回调
	return func(c *framework.Context) {
		fmt.Println("middleware pre test1")
		c.Next()
		fmt.Println("middleware post test1")
	}
}

func Test2() framework.ControllerHandler {
	// 使用函数回调
	return func(c *framework.Context) {
		fmt.Println("middleware pre test2")
		c.Next()
		fmt.Println("middleware post test2")
	}
}

func Test3() framework.ControllerHandler {
	// 使用函数回调
	return func(c *framework.Context) {
		fmt.Println("middleware pre test3")
		c.Next()
		fmt.Println("middleware post test3")
	}
}
