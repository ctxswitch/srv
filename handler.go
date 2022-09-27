package srv

type Handler interface {
	ServeTCP(ResponseWriter, *Request)
}

type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeTCP(w ResponseWriter, r *Request) {
	f(w, r)
}
