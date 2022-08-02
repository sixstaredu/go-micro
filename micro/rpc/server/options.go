package server

type serverOptions struct {
	openssl bool
	cartFile string
	keyFile string
	wraps []HandlerWrapper
}

var defaultServerOptions = serverOptions{}

type ServerOption interface {
	apply(*serverOptions)
}

type EmptyServerOption struct{}

func (EmptyServerOption) apply(*serverOptions) {}

type funcServerOption struct {
	f func(*serverOptions)
}

func (fdo *funcServerOption) apply(do *serverOptions) {
	fdo.f(do)
}

func newFuncServerOption(f func(*serverOptions)) *funcServerOption {
	return &funcServerOption{
		f: f,
	}
}

func SetRSAKey(cartfile, keyfile string) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.cartFile = cartfile
		o.keyFile = keyfile

		if o.cartFile != ""&& o.keyFile != "" {
			o.openssl = true
		}
	})
}

func WithHandlerWrap(hw ...HandlerWrapper) ServerOption {
	return newFuncServerOption(func(options *serverOptions) {
		options.wraps = append(options.wraps, hw...)
	})
}