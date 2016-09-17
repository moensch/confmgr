package confmgr

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/moensch/confmgr/vars"
	"net/http"
	"strings"
)

func TypeToString(keytype int) string {
	switch keytype {
	case vars.TYPE_NOT_FOUND:
		return "none"
	case vars.TYPE_STRING:
		return "string"
	case vars.TYPE_LIST:
		return "list"
	case vars.TYPE_HASH:
		return "hash"
	default:
		return "INVALID"
	}
}

type ValueSource struct {
	Value  string `json:"value"`
	Source string `json:"source"`
}

type LookupListResponse struct {
	Type string        `json:"type"`
	Data []ValueSource `json:"data"`
}

/*
func (r ValueSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Value)
}
*/
func (r LookupListResponse) ToString() string {
	var values = make([]string, len(r.Data))
	for idx, entry := range r.Data {
		values[idx] = entry.Value
	}
	return strings.Join(values, "\n")
}

func (r LookupListResponse) ToJsonString() (string, error) {
	var retval string
	jsonblob, err := json.Marshal(r)
	if err != nil {
		return retval, err
	}
	return string(jsonblob), err
}

type LookupHashResponse struct {
	Type string                 `json:"type"`
	Data map[string]ValueSource `json:"data"`
}

func (r LookupHashResponse) ToString() string {
	keys := make([]string, len(r.Data))
	var i int
	i = 0
	for field, value := range r.Data {
		keys[i] = fmt.Sprintf("%s: %s", field, value.Value)
		i++
	}
	return strings.Join(keys, "\n")
}

func (r LookupHashResponse) ToJsonString() (string, error) {
	var retval string
	jsonblob, err := json.Marshal(r)
	if err != nil {
		return retval, err
	}
	return string(jsonblob), err
}

type LookupStringResponse struct {
	Type string      `json:"type"`
	Data ValueSource `json:"data"`
}

func (r LookupStringResponse) ToString() string {
	return r.Data.Value
}

func (r LookupStringResponse) ToJsonString() (string, error) {
	var retval string
	jsonblob, err := json.Marshal(r)
	if err != nil {
		return retval, err
	}
	return string(jsonblob), err
}

type KeyResponse interface {
	ToString() string
	ToJsonString() (string, error)
}

type HashKeyResponse struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

func (r HashKeyResponse) ToJsonString() (string, error) {
	var retval string
	jsonblob, err := json.Marshal(r)
	if err != nil {
		return retval, err
	}
	return string(jsonblob), err
}

func (r HashKeyResponse) ToString() string {
	keys := make([]string, len(r.Data))
	var i int
	i = 0
	for field, value := range r.Data {
		keys[i] = fmt.Sprintf("%s: %s", field, value)
		i++
	}
	return strings.Join(keys, "\n")
}

type ListKeyResponse struct {
	Type string   `json:"type"`
	Data []string `json:"data"`
}

func (r ListKeyResponse) ToJsonString() (string, error) {
	var retval string
	jsonblob, err := json.Marshal(r)
	if err != nil {
		return retval, err
	}
	return string(jsonblob), err
}

func (r ListKeyResponse) ToString() string {
	return strings.Join(r.Data, "\n")
}

type StringKeyResponse struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func (r StringKeyResponse) ToString() string {
	return r.Data
}

func (r StringKeyResponse) ToJsonString() (string, error) {
	var retval string
	jsonblob, err := json.Marshal(r)
	if err != nil {
		return retval, err
	}
	return string(jsonblob), err
}

func SendResponse(w http.ResponseWriter, r *http.Request, resp KeyResponse) {
	switch r.Header.Get("Accept") {
	case "text/plain":
		stringresp := resp.ToString()
		SendTEXTResponse(w, stringresp)
	default:
		stringresp, _ := resp.ToJsonString()
		SendJSONResponse(w, stringresp)
	}
}

func SendErrorResponse(w http.ResponseWriter, code int, body string) {
	log.WithFields(log.Fields{
		"error": body,
		"code":  code,
	}).Error("HTTP ERROR")
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s\n", body)
}

func SendTEXTResponse(w http.ResponseWriter, resp string) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, resp)
}

func SendJSONResponse(w http.ResponseWriter, resp string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, resp)
}

type DataRequest struct {
	Data string `json:"data"`
}

type GenericRequest struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (r *GenericRequest) AsHash() HashKeyResponse {
	var newtype HashKeyResponse
	newtype.Type = "hash"
	newtype.Data = make(map[string]string)
	for k, v := range r.Data.(map[string]interface{}) {
		newtype.Data[k] = v.(string)
	}
	return newtype
}

func (r *GenericRequest) AsList() ListKeyResponse {
	var newtype ListKeyResponse
	newtype.Type = "list"
	newtype.Data = make([]string, len(r.Data.([]interface{})))
	for k, v := range r.Data.([]interface{}) {
		newtype.Data[k] = v.(string)
	}
	return newtype
}

func (r *GenericRequest) AsString() StringKeyResponse {
	var newtype StringKeyResponse
	newtype.Type = "string"
	newtype.Data = r.Data.(string)
	return newtype
}
