package srv

type middleware interface {
	Middleware(handler Handler) Handler
}

type MiddlewareFunc func(Handler) Handler

func (mw MiddlewareFunc) Middleware(handler Handler) Handler {
	return mw(handler)
}
