package configurator

import (
	"strings"
)

type ConfiguratorOptions struct {
	enableFile    bool
	filename      string
	enableENV     bool
	envPrefix     string
	enableFlag    bool
	enableDefault bool
}

type ConfiguratorOption func(*ConfiguratorOptions)

func WithFileProvider(filename string) ConfiguratorOption {
	return func(co *ConfiguratorOptions) {
		co.enableFile = true
		co.filename = filename
	}
}

func WithENVProvider(prefix string) ConfiguratorOption {
	return func(co *ConfiguratorOptions) {
		co.enableENV = true
		co.envPrefix = prefix
	}
}

func WithFlagProvider() ConfiguratorOption {
	return func(co *ConfiguratorOptions) {
		co.enableFlag = true
	}
}

func WithDefaultProvider() ConfiguratorOption {
	return func(co *ConfiguratorOptions) {
		co.enableDefault = true
	}
}

type Provider interface {
	Provide(interface{}, StructInfo) error
}

func NewConfigurator(options ...ConfiguratorOption) *Configurator {
	opts := &ConfiguratorOptions{
		enableFile:    true,
		filename:      "./config/config.yaml",
		enableENV:     false,
		envPrefix:     "",
		enableFlag:    false,
		enableDefault: false,
	}
	for _, fn := range options {
		fn(opts)
	}

	providers := make([]Provider, 0, 4)
	if opts.enableFile && strings.TrimSpace(opts.filename) != "" {
		providers = append(providers, NewFileProvider(opts.filename))
	}
	if opts.enableENV {
		providers = append(providers, NewENVProvider(opts.envPrefix))
	}
	if opts.enableFlag {
		providers = append(providers, NewFlagProvider())
	}
	if opts.enableDefault {
		providers = append(providers, NewDefaultProvider())
	}

	return &Configurator{
		providers: providers,
	}
}

type Configurator struct {
	providers []Provider
}

func (c *Configurator) Load(v interface{}) error {
	si, err := getStructInfo(v, nil)
	if err != nil {
		return err
	}
	for _, p := range c.providers {
		if err := p.Provide(v, si); err != nil {
			return err
		}
	}
	return nil
}
