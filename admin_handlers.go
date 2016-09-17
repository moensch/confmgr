package confmgr

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/moensch/confmgr/backends"
	"github.com/moensch/confmgr/vars"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func (c *ConfMgr) HandleAdminGetKeyType(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	if !strings.HasPrefix(keyName, c.Config.Main.KeyPrefix) {
		keyName = c.Config.Main.KeyPrefix + keyName
	}

	keytype, err := b.GetType(keyName)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
		return
	}

	if keytype == vars.TYPE_NOT_FOUND {
		SendErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Key %s not found", keyName))
		return
	}
	resp := StringKeyResponse{"string", TypeToString(keytype)}

	SendResponse(w, r, resp)
}

func (c *ConfMgr) HandleAdminKeyStore(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	if !strings.HasPrefix(keyName, c.Config.Main.KeyPrefix) {
		keyName = c.Config.Main.KeyPrefix + keyName
	}
	log.Infof("Storing key %s", keyName)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
		return
	}

	if err := r.Body.Close(); err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
		return
	}
	err = c.SaveKeyFromJSON(keyName, body, b)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *ConfMgr) HandleAdminListHashFields(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	if !strings.HasPrefix(keyName, c.Config.Main.KeyPrefix) {
		keyName = c.Config.Main.KeyPrefix + keyName
	}

	value, err := b.GetHash(keyName)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
		return
	}

	var resp ListKeyResponse
	resp.Type = "list"
	resp.Data = make([]string, len(value))
	idx := 0
	for k, _ := range value {
		resp.Data[idx] = k
		idx++
	}

	SendResponse(w, r, resp)
}

func (c *ConfMgr) HandleAdminListKeys(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend) {
	resp, err := c.ListKeys("", b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
		return
	}

	SendResponse(w, r, resp)
}

func (c *ConfMgr) HandleAdminListKeysFiltered(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend) {
	reqVars := mux.Vars(r)
	filter := reqVars["filter"]
	if !strings.HasPrefix(filter, c.Config.Main.KeyPrefix) {
		filter = c.Config.Main.KeyPrefix + filter
	}

	resp, err := c.ListKeys(filter, b)

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
		return
	}

	SendResponse(w, r, resp)
}

func (c *ConfMgr) ListKeys(filter string, b backend.ConfigBackend) (ListKeyResponse, error) {
	var resp ListKeyResponse
	keys, err := b.ListKeys(filter)

	if err != nil {
		return resp, err
	}

	resp.Type = "list"
	resp.Data = make([]string, len(keys))

	for idx, key := range keys {
		resp.Data[idx] = strings.TrimPrefix(key, c.Config.Main.KeyPrefix)
	}

	sort.Strings(resp.Data)

	return resp, err
}

func (c *ConfMgr) HandleAdminKeyDelete(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	if !strings.HasPrefix(keyName, c.Config.Main.KeyPrefix) {
		keyName = c.Config.Main.KeyPrefix + keyName
	}

	err := b.DeleteKey(keyName)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *ConfMgr) HandleListAppend(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]

	if !strings.HasPrefix(keyName, c.Config.Main.KeyPrefix) {
		keyName = c.Config.Main.KeyPrefix + keyName
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
		return
	}

	if err := r.Body.Close(); err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
		return
	}

	log.Infof("List append to %s: '%s'", keyName, body)
	err = c.ListAppendFromJSON(keyName, body, b)

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *ConfMgr) ListAppendFromJSON(keyName string, jsondata []byte, b backend.ConfigBackend) error {
	var request DataRequest

	if err := json.Unmarshal(jsondata, &request); err != nil {
		return err
	}

	return b.ListAppend(keyName, request.Data)
}
func (c *ConfMgr) HandleAdminSetHashField(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	fieldName := reqVars["fieldName"]
	if !strings.HasPrefix(keyName, c.Config.Main.KeyPrefix) {
		keyName = c.Config.Main.KeyPrefix + keyName
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
		return
	}

	if err := r.Body.Close(); err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
		return
	}

	log.Infof("Set hfield %s/%s to '%s'", keyName, fieldName, body)
	err = c.SetHashFieldFromJSON(keyName, fieldName, body, b)

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *ConfMgr) SetHashFieldFromJSON(keyName string, fieldName string, jsondata []byte, b backend.ConfigBackend) error {
	var request DataRequest

	if err := json.Unmarshal(jsondata, &request); err != nil {
		return err
	}

	log.Infof("Settings %s/%s to '%s'", keyName, fieldName, request.Data)
	return b.SetHashField(keyName, fieldName, request.Data)
}

func (c *ConfMgr) SaveKeyFromJSON(keyName string, jsondata []byte, b backend.ConfigBackend) error {
	var err error

	var keyEntry GenericRequest

	if err := json.Unmarshal(jsondata, &keyEntry); err != nil {
		return err
	}

	switch keyEntry.Type {
	case "string":
		c.StoreString(keyName, keyEntry.AsString(), b)
	case "hash":
		c.StoreHash(keyName, keyEntry.AsHash(), b)
	case "list":
		c.StoreList(keyName, keyEntry.AsList(), b)
	default:
		log.Print("ERROR")
	}

	return err
}

func (c *ConfMgr) StoreString(keyName string, data StringKeyResponse, b backend.ConfigBackend) error {
	err := b.SetString(keyName, data.Data)
	return err
}

func (c *ConfMgr) StoreHash(keyName string, data HashKeyResponse, b backend.ConfigBackend) error {
	err := b.SetHash(keyName, data.Data)
	return err
}

func (c *ConfMgr) StoreList(keyName string, data ListKeyResponse, b backend.ConfigBackend) error {
	err := b.SetList(keyName, data.Data)
	return err
}

/*
 * Retrieve an absolute key (of any type)
 */
func (c *ConfMgr) HandleAdminKeyGet(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	if !strings.HasPrefix(keyName, c.Config.Main.KeyPrefix) {
		keyName = c.Config.Main.KeyPrefix + keyName
	}
	keytype, err := b.GetType(keyName)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
		return
	}

	var resp KeyResponse

	switch keytype {
	case vars.TYPE_NOT_FOUND:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Key %s not found\n", keyName)
		return
	case vars.TYPE_STRING:
		value, err := b.GetString(keyName)
		if err != nil {
			SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
			return
		}
		resp = &StringKeyResponse{
			Type: TypeToString(keytype),
			Data: value,
		}
	case vars.TYPE_LIST:
		value, err := b.GetList(keyName)
		if err != nil {
			SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
			return
		}
		resp = &ListKeyResponse{
			Type: TypeToString(keytype),
			Data: value,
		}
	case vars.TYPE_HASH:
		value, err := b.GetHash(keyName)
		if err != nil {
			SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
			return
		}

		resp = &HashKeyResponse{
			Type: TypeToString(keytype),
			Data: value,
		}
	default:
		SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unsupported key type: %s", TypeToString(keytype)))
		return
	}

	SendResponse(w, r, resp)

}

func (c *ConfMgr) HandleAdminGetHashField(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	if !strings.HasPrefix(keyName, c.Config.Main.KeyPrefix) {
		keyName = c.Config.Main.KeyPrefix + keyName
	}
	fieldName := reqVars["fieldName"]

	keytype, err := b.GetType(keyName)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
		return
	}

	var resp KeyResponse

	switch keytype {
	case vars.TYPE_NOT_FOUND:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Key %s not found\n", keyName)
		return
	case vars.TYPE_HASH:
		exists, err := b.HashFieldExists(keyName, fieldName)
		if err != nil {
			SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
			return
		}
		if exists == false {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Key %s not found\n", keyName)
			return
		}
		value, err := b.GetHashField(keyName, fieldName)
		if err != nil {
			SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Backend error: %s", err))
			return
		}
		log.Infof("Retrieved hash field: Key: '%s'  / Field: '%s'", keyName, fieldName)

		resp = &StringKeyResponse{
			Type: TypeToString(vars.TYPE_STRING),
			Data: value,
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unsupported key type: %s\n", TypeToString(keytype))
		return
	}

	SendResponse(w, r, resp)
}

func (c *ConfMgr) HandleAdminGetListIndex(w http.ResponseWriter, r *http.Request, b backend.ConfigBackend) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	if !strings.HasPrefix(keyName, c.Config.Main.KeyPrefix) {
		keyName = c.Config.Main.KeyPrefix + keyName
	}
	listIndex, _ := strconv.ParseInt(reqVars["listIndex"], 10, 64)

	keytype, err := b.GetType(keyName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
		return
	}

	var resp KeyResponse

	switch keytype {
	case vars.TYPE_NOT_FOUND:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Key %s not found\n", keyName)
		return
	case vars.TYPE_LIST:
		exists, err := b.ListIndexExists(keyName, listIndex)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Backend error: %s\n", err)
			return
		}
		if exists == false {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Key %s not found\n", keyName)
			return
		}
		value, err := b.GetListIndex(keyName, listIndex)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Backend error: %s\n", err)
			return
		}

		resp = &StringKeyResponse{
			Type: TypeToString(vars.TYPE_STRING),
			Data: value,
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unsupported key type: %s\n", TypeToString(keytype))
		return
	}

	SendResponse(w, r, resp)
}
