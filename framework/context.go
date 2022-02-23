package framework

import (
	"context"
	"feng/framework/container"
	"net/http"
	"sync"
	"time"
)

// Context代表当前请求上下文
type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context

	// 是否超时标记位
	hasTimeout bool
	writerMux  *sync.Mutex

	// 当前请求的handler链条
	handlers []ControllerHandler
	index    int // 当前请求调用到调用链的哪个节点

	params map[string]string // url路由匹配的参数

	// Container
	container container.Container
}

// NewContext 初始化一个Context
func NewContext(r *http.Request, w http.ResponseWriter, container container.Container) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		writerMux:      &sync.Mutex{},
		index:          -1,
		container:      container,
	}
}

// #region base function

func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.writerMux
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.responseWriter
}

func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
}

func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

// 为context设置handlers
func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	ctx.handlers = handlers
}

// 设置参数
func (ctx *Context) SetParams(params map[string]string) {
	ctx.params = params
}

// 核心函数，调用context的下一个函数
func (ctx *Context) Next() {
	ctx.index++
	if ctx.index < len(ctx.handlers) {
		ctx.handlers[ctx.index](ctx)
	}
}

// #endregion

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

// #region implement context.Context
func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.BaseContext().Value(key)
}

// context封装container
func (ctx *Context) Make(key string) (interface{}, error) {
	return ctx.container.Make(key)
}

func (ctx *Context) MustMake(key string) interface{} {
	return ctx.container.MustMake(key)
}

func (ctx *Context) MakeNew(key string, params []interface{}) (interface{}, error) {
	return ctx.container.MakeNew(key, params)
}
