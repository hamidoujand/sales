package web

type Middleware func(HandlerFunc) HandlerFunc

func applyMiddleware(h HandlerFunc, mids ...Middleware) HandlerFunc {
	for i := len(mids) - 1; i >= 0; i-- {
		m := mids[i]
		if m != nil {
			h = m(h)
		}
	}
	return h
}
