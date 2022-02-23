package env

import (
	"feng/framework/container"
	"feng/framework/contract"
)

type FengEnvProvider struct {
	Folder string
}

// Register registe a new function for make a service instance
func (provider *FengEnvProvider) Register(c container.Container) container.NewInstance {
	return NewFengEnv
}

// Boot will called when the service instantiate
func (provider *FengEnvProvider) Boot(c container.Container) error {
	app := c.MustMake(contract.AppKey).(contract.App)
	provider.Folder = app.BaseFolder()
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *FengEnvProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *FengEnvProvider) Params(c container.Container) []interface{} {
	return []interface{}{provider.Folder}
}

/// Name define the name for this service
func (provider *FengEnvProvider) Name() string {
	return contract.EnvKey
}
