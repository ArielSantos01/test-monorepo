package apigateway

type options struct {
	pklConfig map[string]any
}

type Option func(*options)

func WithPklConfig(pklConfig map[string]any) Option {
	return func(o *options) {
		o.pklConfig = pklConfig
	}
}
