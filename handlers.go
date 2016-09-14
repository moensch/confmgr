package confmgr

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (c *ConfMgr) Index(w http.ResponseWriter, r *http.Request) {
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

func (c *ConfMgr) GetRequestScope() map[string]string {
	return c.RequestScope
	/*
		scopevars := make(map[string]string)
		scopevars["pod"] = "ml2"
		scopevars["fqdn"] = "samtest"
		scopevars["site"] = "latisys"
		scopevars["group"] = "group1"

		return scopevars
	*/
}

func (c *ConfMgr) HandleLookupHash(w http.ResponseWriter, r *http.Request) {
	c.SetRequestScopeFromHeaders(r.Header)
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]

	log.Printf("Requesting hash lookup: %s", keyName)

	resp, err := c.LookupHash(keyName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
		return
	}

	if len(resp.Data) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Key %s not found\n", keyName)
		return
	}

	SendResponse(w, r, resp)
}

func (c *ConfMgr) HandleLookupString(w http.ResponseWriter, r *http.Request) {
	c.SetRequestScopeFromHeaders(r.Header)
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]

	log.Printf("Requesting string lookup: %s", keyName)

	resp, err := c.LookupString(keyName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
		return
	}

	if resp.Data.Source == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Key %s not found\n", keyName)
		return
	}

	SendResponse(w, r, resp)
}

func (c *ConfMgr) HandleLookupList(w http.ResponseWriter, r *http.Request) {
	c.SetRequestScopeFromHeaders(r.Header)
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]

	log.Printf("Requesting list lookup: %s", keyName)

	resp, err := c.LookupList(keyName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
		return
	}

	if len(resp.Data) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Key %s not found\n", keyName)
		return
	}

	SendResponse(w, r, resp)
}

func (c *ConfMgr) HandleLookupHashField(w http.ResponseWriter, r *http.Request) {
	c.SetRequestScopeFromHeaders(r.Header)
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	fieldName := reqVars["fieldName"]

	log.Printf("Requesting hash field lookup: %s/%s", keyName, fieldName)

	resp, err := c.LookupHashField(keyName, fieldName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
		return
	}

	if resp.Data.Source == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Key %s not found\n", keyName)
		return
	}

	SendResponse(w, r, resp)
}

func (c *ConfMgr) HandleLookupListIndex(w http.ResponseWriter, r *http.Request) {
	c.SetRequestScopeFromHeaders(r.Header)
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	listIndex, _ := strconv.ParseInt(reqVars["listIndex"], 10, 64)

	log.Printf("Requesting list index lookup: %s[%d]", keyName, listIndex)

	resp, err := c.LookupListIndex(keyName, listIndex)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
		return
	}

	if resp.Data.Source == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Key %s not found\n", keyName)
		return
	}

	SendResponse(w, r, resp)
}