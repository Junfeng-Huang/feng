package app

import (
	"feng/framework/container"
	"feng/framework/contract"
)

// FengAppProvider 提供App的具体实现方法
type FengAppProvider struct {
	BaseFolder string
}

// Register 注册FengApp方法
func (f *FengAppProvider) Register(container container.Container) container.NewInstance {
	return NewFengApp
}

// Boot 启动调用
func (f *FengAppProvider) Boot(container container.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (f *FengAppProvider) IsDefer() bool {
	return false
}

// Params 获取初始化参数
func (f *FengAppProvider) Params(container container.Container) []interface{} {
	return []interface{}{container, f.BaseFolder}
}

// Name 获取字符串凭证
func (f *FengAppProvider) Name() string {
	return contract.AppKey
}
