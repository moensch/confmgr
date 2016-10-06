package confmgr

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/moensch/confmgr/backends"
	"net/http"
	"time"
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
		start := time.Now()
		scope := ScopeFromHeaders(r.Header)
		context.Set(r, ReqScope, scope)

		b := BackendFactory.NewBackend()
		defer b.Close()
		f(w, r, b)
		log.WithFields(log.Fields{
			"method": r.Method,
			"uri":    r.RequestURI,
			"client": r.RemoteAddr,
			"time":   time.Since(start),
			"scope":  scope,
		}).Info("Request")
		context.Clear(r)
	})
}
