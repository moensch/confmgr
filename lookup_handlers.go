package confmgr

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/moensch/confmgr/backends"
	"net/http"
	"strconv"
)

func (c *ConfMgr) HandleLookupHash(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]

	//log.Printf("Requesting hash lookup: %s", keyName)

	resp, err := c.LookupHash(keyName, GetRequestScope(r), b)
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

func (c *ConfMgr) HandleLookupString(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]

	//log.Printf("Requesting string lookup: %s", keyName)

	resp, err := c.LookupString(keyName, GetRequestScope(r), b)
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

func (c *ConfMgr) HandleLookupList(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]

	//log.Printf("Requesting list lookup: %s", keyName)

	resp, err := c.LookupList(keyName, GetRequestScope(r), b)
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

func (c *ConfMgr) HandleLookupHashField(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	fieldName := reqVars["fieldName"]

	//log.Printf("Requesting hash field lookup: %s/%s", keyName, fieldName)

	resp, err := c.LookupHashField(keyName, fieldName, GetRequestScope(r), b)
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

func (c *ConfMgr) HandleLookupListIndex(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	listIndex, _ := strconv.ParseInt(reqVars["listIndex"], 10, 64)

	//log.Printf("Requesting list index lookup: %s[%d]", keyName, listIndex)

	resp, err := c.LookupListIndex(keyName, listIndex, GetRequestScope(r), b)
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
