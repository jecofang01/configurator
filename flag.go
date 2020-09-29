package configurator

type flagProvider struct{}

func NewFlagProvider() *flagProvider {
	return &flagProvider{}
}

func (p flagProvider) Provide(v interface{}) error {
	return nil
}
