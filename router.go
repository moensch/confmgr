package confmgr

import (
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/moensch/confmgr/backends"
	"net/http"
)

func (c *ConfMgr) NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range c.RouteDefinitions() {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = c.ClientHandler(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

const ReqScope = 0

type HandlerFuncBackend func(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend)

func handlerDecorate(f HandlerFuncBackend) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, ReqScope, ScopeFromHeaders(r.Header))

		b := backendFactory.NewBackend()
		defer b.Close()
		f(w, r, b)
		context.Clear(r)
	})
}
