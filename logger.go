package confmgr

import (
	"net/http"
)

func (c *ConfMgr) ClientHandler(inner http.Handler, name string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inner.ServeHTTP(w, r)
	})
}
