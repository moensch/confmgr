package confmgr

import (
	"fmt"
	"github.com/gorilla/context"
	"github.com/moensch/confmgr/backends"
	"net/http"
	"strings"
)

func (c *ConfMgr) Index(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend) {
	fmt.Fprintln(w, "Welcome!")
}

/*
 * Extract scope variables from x-cfg-FOO request headers
 */
func (c *ConfMgr) SetRequestScopeFromHeaders(headers http.Header) {
	c.RequestScope = make(map[string]string)
	for hdrname, hdrval := range headers {
		hdrname = strings.ToLower(hdrname)
		if strings.HasPrefix(hdrname, "x-cfg-") {
			scopevar := strings.TrimPrefix(hdrname, "x-cfg-")
			scopeval := strings.ToLower(hdrval[0])
			c.RequestScope[scopevar] = scopeval
		}
	}
}

func ScopeFromHeaders(headers http.Header) map[string]string {
	scope := make(map[string]string)
	for hdrname, hdrval := range headers {
		hdrname = strings.ToLower(hdrname)
		if strings.HasPrefix(hdrname, "x-cfg-") {
			scopevar := strings.TrimPrefix(hdrname, "x-cfg-")
			scopeval := strings.ToLower(hdrval[0])
			scope[scopevar] = scopeval
		}
	}

	return scope
}
func GetRequestScope(r *http.Request) map[string]string {
	return context.Get(r, ReqScope).(map[string]string)
}
