package core

import "context"

type ModuleConfig struct {
	Name string
}

type Module interface {
	Initialize(config *ModuleConfig) error
	Execute(ctx context.Context) error
	Stop() error
}
