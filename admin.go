package confmgr

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/moensch/confmgr/vars"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func (c *ConfMgr) HandleGetKeyType(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	if !strings.HasPrefix(keyName, c.Config.Main.KeyPrefix) {
		keyName = c.Config.Main.KeyPrefix + keyName
	}

	keytype, err := c.Backend.GetType(keyName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
		return
	}

	if keytype == vars.TYPE_NOT_FOUND {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Key %s not found\n", keyName)
		return
	}
	resp := StringKeyResponse{"string", TypeToString(keytype)}

	SendResponse(w, r, resp)
}

func (c *ConfMgr) KeyStore(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	if !strings.HasPrefix(keyName, c.Config.Main.KeyPrefix) {
		keyName = c.Config.Main.KeyPrefix + keyName
	}
	log.Printf("Storing key %s", keyName)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
		return
	}

	if err := r.Body.Close(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
		return
	}
	err = c.SaveKeyFromJSON(keyName, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *ConfMgr) HandleListHashFields(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	if !strings.HasPrefix(keyName, c.Config.Main.KeyPrefix) {
		keyName = c.Config.Main.KeyPrefix + keyName
	}

	value, err := c.Backend.GetHash(keyName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
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

func (c *ConfMgr) HandleListKeys(w http.ResponseWriter, r *http.Request) {
	resp, err := c.ListKeys("")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
		return
	}

	SendResponse(w, r, resp)
}

func (c *ConfMgr) HandleListKeysFiltered(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	filter := reqVars["filter"]
	if !strings.HasPrefix(filter, c.Config.Main.KeyPrefix) {
		filter = c.Config.Main.KeyPrefix + filter
	}

	resp, err := c.ListKeys(filter)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
		return
	}

	SendResponse(w, r, resp)
}

func (c *ConfMgr) ListKeys(filter string) (ListKeyResponse, error) {
	var resp ListKeyResponse
	keys, err := c.Backend.ListKeys(filter)

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

func (c *ConfMgr) HandleKeyDelete(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	if !strings.HasPrefix(keyName, c.Config.Main.KeyPrefix) {
		keyName = c.Config.Main.KeyPrefix + keyName
	}

	err := c.Backend.DeleteKey(keyName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *ConfMgr) HandleSetHashField(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	fieldName := reqVars["fieldName"]
	if !strings.HasPrefix(keyName, c.Config.Main.KeyPrefix) {
		keyName = c.Config.Main.KeyPrefix + keyName
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
		return
	}

	if err := r.Body.Close(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
		return
	}

	log.Printf("Set hfield %s/%s to '%s'", keyName, fieldName, body)
	err = c.SetHashFieldFromJSON(keyName, fieldName, body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Backend error: %s\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *ConfMgr) SetHashFieldFromJSON(keyName string, fieldName string, jsondata []byte) error {
	var request DataRequest

	if err := json.Unmarshal(jsondata, &request); err != nil {
		return err
	}

	log.Printf("Settings %s/%s to '%s'", keyName, fieldName, request.Data)
	return c.Backend.SetHashField(keyName, fieldName, request.Data)
}

func (c *ConfMgr) SaveKeyFromJSON(keyName string, jsondata []byte) error {
	var err error

	var keyEntry GenericRequest

	if err := json.Unmarshal(jsondata, &keyEntry); err != nil {
		return err
	}

	switch keyEntry.Type {
	case "string":
		c.StoreString(keyName, keyEntry.AsString())
	case "hash":
		c.StoreHash(keyName, keyEntry.AsHash())
	case "list":
		c.StoreList(keyName, keyEntry.AsList())
	default:
		log.Print("ERROR")
	}

	return err
}

func (c *ConfMgr) StoreString(keyName string, data StringKeyResponse) error {
	err := c.Backend.SetString(keyName, data.Data)
	return err
}

func (c *ConfMgr) StoreHash(keyName string, data HashKeyResponse) error {
	err := c.Backend.SetHash(keyName, data.Data)
	return err
}

func (c *ConfMgr) StoreList(keyName string, data ListKeyResponse) error {
	err := c.Backend.SetList(keyName, data.Data)
	return err
}

/*
 * Retrieve an absolute key (of any type)
 */
func (c *ConfMgr) KeyGet(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	if !strings.HasPrefix(keyName, c.Config.Main.KeyPrefix) {
		keyName = c.Config.Main.KeyPrefix + keyName
	}
	keytype, err := c.Backend.GetType(keyName)
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
	case vars.TYPE_STRING:
		value, err := c.Backend.GetString(keyName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Backend error: %s\n", err)
			return
		}
		resp = &StringKeyResponse{
			Type: TypeToString(keytype),
			Data: value,
		}
	case vars.TYPE_LIST:
		value, err := c.Backend.GetList(keyName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Backend error: %s\n", err)
			return
		}
		resp = &ListKeyResponse{
			Type: TypeToString(keytype),
			Data: value,
		}
	case vars.TYPE_HASH:
		value, err := c.Backend.GetHash(keyName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Backend error: %s\n", err)
			return
		}

		resp = &HashKeyResponse{
			Type: TypeToString(keytype),
			Data: value,
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unsupported key type: %s\n", TypeToString(keytype))
		return
	}

	SendResponse(w, r, resp)

}

func (c *ConfMgr) KeyGetHashField(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	if !strings.HasPrefix(keyName, c.Config.Main.KeyPrefix) {
		keyName = c.Config.Main.KeyPrefix + keyName
	}
	fieldName := reqVars["fieldName"]

	keytype, err := c.Backend.GetType(keyName)
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
	case vars.TYPE_HASH:
		exists, err := c.Backend.HashFieldExists(keyName, fieldName)
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
		value, err := c.Backend.GetHashField(keyName, fieldName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Backend error: %s\n", err)
			return
		}
		log.Printf("Retrieved hash field: Key: '%s'  / Field: '%s'", keyName, fieldName)

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

func (c *ConfMgr) KeyGetListIndex(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	keyName := reqVars["keyName"]
	if !strings.HasPrefix(keyName, c.Config.Main.KeyPrefix) {
		keyName = c.Config.Main.KeyPrefix + keyName
	}
	listIndex, _ := strconv.ParseInt(reqVars["listIndex"], 10, 64)

	keytype, err := c.Backend.GetType(keyName)
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
		exists, err := c.Backend.ListIndexExists(keyName, listIndex)
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
		value, err := c.Backend.GetListIndex(keyName, listIndex)
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
