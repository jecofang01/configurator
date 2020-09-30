package configurator

import (
	"os"
	"strings"
)

type envProvider struct {
	prefix string
}

func NewENVProvider(prefix string) *envProvider {
	return &envProvider{
		prefix: strings.ToUpper(prefix),
	}
}

func (p envProvider) Provide(v interface{}, si StructInfo) error {
	for _, fi := range si.Fields() {
		k := p.normalize(fi.ENVKey())
		if k == "" {
			continue
		}
		val, ok := os.LookupEnv(k)
		if !ok {
			continue
		}
		_ = val // TODO set value
	}
	return nil
}

func (p envProvider) normalize(key string) string {
	if key == "" {
		return ""
	}
	if p.prefix == "" || strings.HasPrefix(key, p.prefix) {
		return key
	}
	return strings.Join([]string{p.prefix, key}, "_")
}
