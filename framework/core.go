package framework

import (
	"context"
	"feng/framework/container"
	"feng/framework/provider/app"
	"feng/framework/provider/config"
	"feng/framework/provider/env"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Core represent core struct
type Core struct {
	router      map[string]*Tree    // all routers
	middlewares []ControllerHandler // 从core这边设置的中间件
	container   container.Container
	closeWait   time.Duration
}

// 初始化core结构
func NewCore() *Core {
	// 初始化路由
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()
	core := &Core{router: router, container: container.NewFengContainer(), closeWait: 3 * time.Second}
	return core
}

// 初始化core并且bind默认服务
func Default() *Core {
	core := NewCore()
	// TODO: Bind default provider
	core.Bind(&app.FengAppProvider{})
	core.Bind(&env.FengEnvProvider{})
	core.Bind(&config.FengConfigProvider{})

	return core
}

// 注册中间件
func (c *Core) Use(middlewares ...ControllerHandler) {
	c.middlewares = middlewares
}

// === http method 封装

// 匹配GET 方法, 增加路由规则
func (c *Core) GET(url string, handlers ...ControllerHandler) {
	// 将core的middleware 和 handlers结合起来
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["GET"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 匹配POST 方法, 增加路由规则
func (c *Core) POST(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["POST"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 匹配PUT 方法, 增加路由规则
func (c *Core) PUT(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["PUT"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 匹配DELETE 方法, 增加路由规则
func (c *Core) DELETE(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["DELETE"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// ==== http method 封装完成

func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

// 匹配路由，如果没有匹配到，返回nil
func (c *Core) FindRouteNodeByRequest(request *http.Request) *node {
	// uri 和 method 全部转换为大写，保证大小写不敏感
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)

	// 查找第一层map
	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.root.matchNode(uri)
	}
	return nil
}

// core封装container
func (core *Core) Bind(provider container.ServiceProvider) error {
	return core.container.Bind(provider)
}

func (core *Core) IsBind(key string) bool {
	return core.container.IsBind(key)
}

// closeWait
func (core *Core) GetCloseWait() time.Duration {
	return core.closeWait
}

func (core *Core) SetCloseWait(t time.Duration) {
	core.closeWait = t
}

// 直接启动框架
func (core *Core) Run(addr ...string) {
	var address string
	switch len(addr) {
	case 0:
		address = ":8080"
	case 1:
		address = addr[0]
	default:
		panic("too many parameters")
	}
	server := &http.Server{
		Handler: core,
		Addr:    address,
	}

	// 这个goroutine是启动服务的goroutine
	go func() {
		server.ListenAndServe()
	}()

	//  设置等待信号量。可以不在main里面设置，在其他Goroutine也可以收到信息
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 等待信号
	<-quit

	// 调用Server.Shutdown graceful结束
	timeoutCtx, cancel := context.WithTimeout(context.Background(), core.closeWait)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}

// 所有请求都进入这个函数, 这个函数负责路由分发
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	// 封装自定义context
	ctx := NewContext(request, response, c.container)
	// 寻找路由
	node := c.FindRouteNodeByRequest(request)
	if node == nil {
		// 如果没有找到，这里打印日志
		ctx.SetStatus(404).Json("not found")
		return
	}

	ctx.SetHandlers(node.handlers)

	// 设置路由参数
	params := node.parseParamsFromEndNode(request.URL.Path)
	ctx.SetParams(params)

	// 调用路由函数，如果返回err 代表存在内部错误，返回500状态码
	ctx.Next()
}
