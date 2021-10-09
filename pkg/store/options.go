package store

type listKeysOptions struct {
	prefix *string
}

type ListKeysOption func(*listKeysOptions)

func getListKeysOptions(opts []ListKeysOption) *listKeysOptions {
	options := &listKeysOptions{}
	for _, applyOption := range opts {
		applyOption(options)
	}
	return options
}

func WithPrefix(prefix string) ListKeysOption {
	return func(options *listKeysOptions) {
		options.prefix = &prefix
	}
}
