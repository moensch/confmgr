package confmgr

import (
	"github.com/gorilla/mux"
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
