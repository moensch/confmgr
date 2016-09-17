package confmgr

import (
	"github.com/gorilla/context"
	"log"
	"net/http"
	"time"
)

const ReqScope = 0

func (c *ConfMgr) ClientHandler(inner http.Handler, name string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, ReqScope, ScopeFromHeaders(r.Header))
		start := time.Now()
		inner.ServeHTTP(w, r)
		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
		context.Clear(r)
	})
}
