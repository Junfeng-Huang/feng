package app

import (
	"errors"
	"feng/framework/container"
	"feng/framework/util"
	"flag"
	"path/filepath"
)

// FengApp 代表feng框架的App实现
type FengApp struct {
	container  container.Container // 服务容器
	baseFolder string              // 基础路径
}

// Version 实现版本
func (f FengApp) Version() string {
	return "0.0.3"
}

// BaseFolder 表示基础目录，可以代表开发场景的目录，也可以代表运行时候的目录
func (f *FengApp) BaseFolder() string {
	if f.baseFolder != "" {
		return f.baseFolder
	}

	// 如果没有设置，则使用参数
	var baseFolder string
	flag.StringVar(&baseFolder, "base_folder", "", "base_folder参数, 默认为当前路径")
	flag.Parse()
	if baseFolder != "" {
		f.baseFolder = baseFolder
		return baseFolder
	}

	// 如果参数也没有，使用默认的当前路径
	f.baseFolder = util.GetExecDirectory()
	return f.baseFolder
}

// ConfigFolder  表示配置文件地址
func (f FengApp) ConfigFolder() string {
	return filepath.Join(f.BaseFolder(), "config")
}

// LogFolder 表示日志存放地址
func (f FengApp) LogFolder() string {
	return filepath.Join(f.StorageFolder(), "log")
}

func (f FengApp) HttpFolder() string {
	return filepath.Join(f.BaseFolder(), "appHttp")
}

func (f FengApp) ConsoleFolder() string {
	return filepath.Join(f.BaseFolder(), "console")
}

func (f FengApp) StorageFolder() string {
	return filepath.Join(f.BaseFolder(), "storage")
}

// ProviderFolder 定义业务自己的服务提供者地址
func (f FengApp) ProviderFolder() string {
	return filepath.Join(f.BaseFolder(), "provider")
}

// MiddlewareFolder 定义业务自己定义的中间件
func (f FengApp) MiddlewareFolder() string {
	return filepath.Join(f.HttpFolder(), "middleware")
}

// CommandFolder 定义业务定义的命令
func (f FengApp) CommandFolder() string {
	return filepath.Join(f.ConsoleFolder(), "command")
}

// RuntimeFolder 定义业务的运行中间态信息
func (f FengApp) RuntimeFolder() string {
	return filepath.Join(f.StorageFolder(), "runtime")
}

// TestFolder 定义测试需要的信息
func (f FengApp) TestFolder() string {
	return filepath.Join(f.BaseFolder(), "test")
}

// NewFengApp 初始化FengApp
func NewFengApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}

	// 有两个参数，一个是容器，一个是baseFolder
	container := params[0].(container.Container)
	baseFolder := params[1].(string)
	return &FengApp{baseFolder: baseFolder, container: container}, nil
}
