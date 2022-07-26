package config

import (
	"feng/framework/container"
	"feng/framework/contract"
	"feng/framework/util"
	"path/filepath"
)

type FengConfigProvider struct {
	folder string // 文件夹。三种获取方式：初始化指定 > 环境变量 > App服务。
}

// Register registe a new function for make a service instance
func (provider *FengConfigProvider) Register(c container.Container) container.NewInstance {
	return NewFengConfig
}

// Boot will called when the service instantiate
func (provider *FengConfigProvider) Boot(c container.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *FengConfigProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *FengConfigProvider) Params(c container.Container) []interface{} {
	envService := c.MustMake(contract.EnvKey).(contract.Env)
	env := envService.AppEnv()
	var envFolder string // configFolder下的AppEnv路径
	// 配置文件夹地址
	switch true {
	case provider.folder != "":
		envFolder = filepath.Join(provider.folder, env)
	case c.IsBind(contract.AppKey):
		appService := c.MustMake(contract.AppKey).(contract.App)
		configFolder := appService.ConfigFolder()
		envFolder = filepath.Join(configFolder, env)
	default:
		execDirectory := util.GetExecDirectory()
		configFolder := filepath.Join(execDirectory, "config")
		envFolder = filepath.Join(configFolder, env)
	}
	return []interface{}{c, envFolder, envService.All()}
}

/// Name define the name for this service
func (provider *FengConfigProvider) Name() string {
	return contract.ConfigKey
}
