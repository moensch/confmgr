package confmgr

import (
	"errors"
	"fmt"
	"github.com/moensch/confmgr/vars"
	"regexp"
	"strconv"
	"strings"
)

func (c *ConfMgr) LookupString(keyName string) (LookupStringResponse, error) {
	var resp LookupStringResponse
	var err error

	for _, keyName := range c.ExistingKeys(keyName, vars.TYPE_STRING) {
		stringdata, err := c.Backend.GetString(keyName)
		stringdata = c.SubstituteValues(stringdata)

		if err != nil {
			return resp, err
		}

		resp.Data = ValueSource{stringdata, keyName}
	}

	resp.Type = TypeToString(vars.TYPE_STRING)

	return resp, err
}

func (c *ConfMgr) LookupHash(keyName string) (LookupHashResponse, error) {
	var resp LookupHashResponse
	var err error

	var hashes_to_merge = make([]map[string]ValueSource, 0)

	for _, keyName := range c.ExistingKeys(keyName, vars.TYPE_HASH) {
		hashdata, err := c.Backend.GetHash(keyName)

		if err != nil {
			return resp, err
		}

		var valuesource = make(map[string]ValueSource)
		for k, v := range hashdata {
			valuesource[k] = ValueSource{c.SubstituteValues(v), keyName}
		}
		hashes_to_merge = append(hashes_to_merge, valuesource)
	}

	// Merge found hashes
	for _, hash := range hashes_to_merge {
		if len(resp.Data) == 0 {
			resp.Data = hash
		} else {
			// Override all existing keys
			for k, v := range hash {
				resp.Data[k] = v
			}
		}
	}

	resp.Type = TypeToString(vars.TYPE_HASH)

	return resp, err
}

func (c *ConfMgr) LookupHashFieldByString(searchString string) string {
	// ${key/fieldname}
	hash_field_vars := regexp.MustCompile("\\${(\\S+?)/(\\S+?)}")
	matches := hash_field_vars.FindAllStringSubmatch(searchString, -1)
	if len(matches) > 0 {
		keyName := matches[0][1]
		fieldName := matches[0][2]
		resp, err := c.LookupHashField(keyName, fieldName)

		if err != nil {
			return ""
		}

		return resp.ToString()
	}
	return ""
}

func (c *ConfMgr) LookupHashField(keyName string, fieldName string) (LookupStringResponse, error) {
	var resp LookupStringResponse
	var err error

	for _, keyName := range c.ExistingKeys(keyName, vars.TYPE_HASH) {
		exists, err := c.Backend.HashFieldExists(keyName, fieldName)
		if err != nil {
			return resp, err
		}
		if exists {
			stringdata, err := c.Backend.GetHashField(keyName, fieldName)
			stringdata = c.SubstituteValues(stringdata)

			if err != nil {
				return resp, err
			}

			resp.Data = ValueSource{stringdata, keyName}
		}
	}

	resp.Type = TypeToString(vars.TYPE_STRING)
	return resp, err
}
func (c *ConfMgr) LookupList(keyName string) (LookupListResponse, error) {
	var resp LookupListResponse
	var err error

	for _, keyName := range c.ExistingKeys(keyName, vars.TYPE_LIST) {
		listdata, err := c.Backend.GetList(keyName)

		if err != nil {
			return resp, err
		}

		for _, entry := range listdata {
			valuesource := ValueSource{entry, keyName}
			resp.Data = append(resp.Data, valuesource)
		}
	}

	resp.Type = TypeToString(vars.TYPE_LIST)
	return resp, err
}
func (c *ConfMgr) LookupListIndex(keyName string, listIndex int64) (LookupStringResponse, error) {
	var resp LookupStringResponse
	var err error

	list, err := c.LookupList(keyName)
	if err != nil {
		return resp, err
	}

	if int(listIndex) >= len(list.Data) {
		return resp, errors.New("List index out of range")
	}
	resp.Data = list.Data[listIndex]

	resp.Type = TypeToString(vars.TYPE_STRING)
	return resp, err
}

func (c *ConfMgr) LookupListIndexByString(searchString string) string {
	// ${key/index/3}
	hash_field_vars := regexp.MustCompile("\\${(\\S+?)/index/(\\d+?)}")
	matches := hash_field_vars.FindAllStringSubmatch(searchString, -1)
	if len(matches) > 0 {
		keyName := matches[0][1]
		listIndex, _ := strconv.ParseInt(matches[0][2], 10, 64)
		resp, err := c.LookupListIndex(keyName, listIndex)

		if err != nil {
			return ""
		}

		return resp.ToString()
	}
	return ""
}

/**
 * Finds other keys like ${db_policy/host} in input string
 * and replaces them with a lookup value
 */
func (c *ConfMgr) SubstituteValues(input string) string {
	hash_field_vars := regexp.MustCompile("\\${(\\S+?/\\S+?)}")
	list_index_vars := regexp.MustCompile("\\${(\\S+?/index/\\S+?)}")

	input = hash_field_vars.ReplaceAllStringFunc(input, c.LookupHashFieldByString)
	input = list_index_vars.ReplaceAllStringFunc(input, c.LookupListIndexByString)
	return input
}

/*
 * Return all matches based on search path and partial key name
 */
func (c *ConfMgr) ExistingKeys(key string, wantedType int) []string {
	foundKeys := make([]string, 0)

	for _, path := range c.SearchPaths() {
		keyName := fmt.Sprintf("%s%s:%s", c.Config.Main.KeyPrefix, path, key)
		//log.Printf("Searching key: '%s'", keyName)
		keytype, _ := c.Backend.GetType(keyName)

		if keytype == wantedType {
			foundKeys = append(foundKeys, keyName)
		}
	}

	return foundKeys
}

/**
 * Returns the paths in reverse order as this is how all other functions
 * will consume it
 **/
func (c *ConfMgr) SearchPaths() []string {
	reqscope := c.GetRequestScope()
	//log.Printf("Scope: %q\n", reqscope)

	newKeyPaths := make([]string, 0)

	// Find all tokens in the configured paths
	token_re := regexp.MustCompile("%{(\\S+?)}")

PathsLoop:
	for _, path := range c.Config.Main.KeyPaths {
		//log.Printf("Key path: '%s'", path)
		matches := token_re.FindAllStringSubmatch(path, -1)
		if len(matches) > 0 {
			for _, match := range matches {
				token_str := match[0]
				token_name := match[1]
				//log.Printf("  token_str : %s\n", token_str)
				//log.Printf("  token_name: %s\n", token_name)

				// Check if this token is defined in the request scope
				if val, ok := reqscope[token_name]; ok {
					//log.Printf("   Token '%s' is defined in the request as '%s'", token_name, val)
					path = strings.Replace(path, token_str, val, 1)
					//log.Printf("      Modified path: %s", path)
				} else {
					// Cannot replace this token, ignore this path completely
					//log.Printf("      NOT USING THIS PATH")
					continue PathsLoop
				}
			}
		}

		// Poor man's prepend. Make new slice with current value at the beginning
		// Then iterate over existing entries and append them one by one
		var tmpslice = make([]string, 1)
		tmpslice[0] = path
		for _, entry := range newKeyPaths {
			tmpslice = append(tmpslice, entry)
		}
		newKeyPaths = tmpslice
	}

	return newKeyPaths
}
