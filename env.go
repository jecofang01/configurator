package configurator

type envProvider struct {
	prefix string
}

func NewENVProvider(prefix string) *envProvider {
	return &envProvider{
		prefix: prefix,
	}
}

func (p envProvider) Provide(v interface{}) error {
	return nil
}
