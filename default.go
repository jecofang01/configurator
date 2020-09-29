package configurator

type defaultProvider struct{}

func NewDefaultProvider() *defaultProvider {
	return &defaultProvider{}
}

func (p defaultProvider) Provide(v interface{}) error {
	return nil
}
