package confmgr

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/moensch/confmgr/backends"
	"github.com/moensch/confmgr/vars"
	"regexp"
	"strconv"
	"strings"
)

func (c *ConfMgr) LookupString(keyName string, scope map[string]string, b backend.ConfigBackend) (LookupStringResponse, error) {
	var resp LookupStringResponse
	var err error

	for _, keyName := range c.ExistingKeys(keyName, vars.TYPE_STRING, scope, b) {
		stringdata, err := b.GetString(keyName)
		if err != nil {
			return resp, err
		}
		stringdata, err = c.SubstituteValues(stringdata, scope, b)
		if err != nil {
			return resp, err
		}

		resp.Data = ValueSource{stringdata, keyName}
	}

	resp.Type = TypeToString(vars.TYPE_STRING)

	return resp, err
}

func (c *ConfMgr) LookupHash(keyName string, scope map[string]string, b backend.ConfigBackend) (LookupHashResponse, error) {
	var resp LookupHashResponse
	var err error

	var hashes_to_merge = make([]map[string]ValueSource, 0)

	for _, keyName := range c.ExistingKeys(keyName, vars.TYPE_HASH, scope, b) {
		hashdata, err := b.GetHash(keyName)

		if err != nil {
			return resp, err
		}

		var valuesource = make(map[string]ValueSource)
		for k, v := range hashdata {
			v, err := c.SubstituteValues(v, scope, b)
			if err != nil {
				return resp, err
			}
			valuesource[k] = ValueSource{v, keyName}
			if err != nil {
				return resp, err
			}
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

func (c *ConfMgr) LookupStringByString(searchString string, scope map[string]string, b backend.ConfigBackend) (string, error) {
	// ${key}
	string_vars := regexp.MustCompile("\\${(\\S+?)}")
	matches := string_vars.FindAllStringSubmatch(searchString, -1)
	if len(matches) > 0 {
		keyName := matches[0][1]
		resp, err := c.LookupString(keyName, scope, b)

		if err != nil {
			return "", err
		}

		return resp.ToString(), nil
	}
	// No match - returning original
	return searchString, nil
}

func (c *ConfMgr) LookupHashFieldByString(searchString string, scope map[string]string, b backend.ConfigBackend) (string, error) {
	// ${key/fieldname}
	hash_field_vars := regexp.MustCompile("\\${(\\S+?)/(\\S+?)}")
	matches := hash_field_vars.FindAllStringSubmatch(searchString, -1)
	if len(matches) > 0 {
		keyName := matches[0][1]
		fieldName := matches[0][2]
		resp, err := c.LookupHashField(keyName, fieldName, scope, b)

		if err != nil {
			return "", err
		}

		return resp.ToString(), nil
	}
	// No matches - return original
	return searchString, nil
}

func (c *ConfMgr) LookupHashField(keyName string, fieldName string, scope map[string]string, b backend.ConfigBackend) (LookupStringResponse, error) {
	var resp LookupStringResponse
	var err error

	var foundAny bool

	for _, keyName := range c.ExistingKeys(keyName, vars.TYPE_HASH, scope, b) {
		exists, err := b.HashFieldExists(keyName, fieldName)
		if err != nil {
			return resp, err
		}
		if exists {
			stringdata, err := b.GetHashField(keyName, fieldName)
			if err != nil {
				return resp, err
			}
			stringdata, err = c.SubstituteValues(stringdata, scope, b)
			if err != nil {
				return resp, err
			}

			foundAny = true

			resp.Data = ValueSource{stringdata, keyName}
		}
	}

	resp.Type = TypeToString(vars.TYPE_STRING)
	if !foundAny {
		return resp, fmt.Errorf("Unable to find hash field: %s/%s", keyName, fieldName)
	}
	return resp, err
}

func (c *ConfMgr) LookupList(keyName string, scope map[string]string, b backend.ConfigBackend) (LookupListResponse, error) {
	var resp LookupListResponse
	var err error

	for _, keyName := range c.ExistingKeys(keyName, vars.TYPE_LIST, scope, b) {
		listdata, err := b.GetList(keyName)

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

func (c *ConfMgr) LookupListIndex(keyName string, listIndex int64, scope map[string]string, b backend.ConfigBackend) (LookupStringResponse, error) {
	var resp LookupStringResponse
	var err error

	list, err := c.LookupList(keyName, scope, b)
	if err != nil {
		return resp, err
	}

	resp.Type = TypeToString(vars.TYPE_STRING)
	if int(listIndex) >= len(list.Data) {
		resp.Data = ValueSource{"", ""}
		err = fmt.Errorf("Cannot find list index %d in list %s (only has %d entries)", listIndex, keyName, len(list.Data))
	} else {
		resp.Data = list.Data[listIndex]
	}

	return resp, err
}

func (c *ConfMgr) LookupListIndexByString(searchString string, scope map[string]string, b backend.ConfigBackend) (string, error) {
	// ${key/index/3}
	hash_field_vars := regexp.MustCompile("\\${(\\S+?)/index/(\\d+?)}")
	matches := hash_field_vars.FindAllStringSubmatch(searchString, -1)
	if len(matches) > 0 {
		keyName := matches[0][1]
		listIndex, _ := strconv.ParseInt(matches[0][2], 10, 64)
		resp, err := c.LookupListIndex(keyName, listIndex, scope, b)

		if err != nil {
			return "", err
		}

		return resp.ToString(), nil
	}
	// No match - returning original
	return searchString, nil
}

/**
 * Finds other keys like ${db_policy/host} in input string
 * and replaces them with a lookup value
 */
func (c *ConfMgr) SubstituteValues(input string, scope map[string]string, b backend.ConfigBackend) (string, error) {
	log.Debugf("Performing substitution in: %s", input)
	replacements := make(map[string]string)

	// You know, in a perfect world I could just use
	//  regexp.ReplaceAllStringFunc().
	// My code is dirty, alas I cannot

	// Find anything that looks like a variable: ${somehash/somefield}
	var_re := regexp.MustCompile("\\${\\S+?}")
	matches := var_re.FindAllStringSubmatch(input, -1)
	if len(matches) > 0 {
		for _, match := range matches {
			if _, ok := replacements[match[0]]; ok {
				// Already got this one
				continue
			}
			var replace string
			var err error
			switch {
			case strings.Contains(match[0], "/index/"):
				log.Debugf("Substituting list index: %s", match[0])
				// Array index match: ${somearray/index/0}
				replace, err = c.LookupListIndexByString(match[0], scope, b)
			case strings.Contains(match[0], "/"):
				log.Debugf("Substituting hash field: %s", match[0])
				// Hash field match: ${somehash/somefield}
				replace, err = c.LookupHashFieldByString(match[0], scope, b)
			default:
				log.Debugf("Substituting string var: %s", match[0])
				// Hash field match: ${somestringvar}
				replace, err = c.LookupStringByString(match[0], scope, b)
			}
			if err != nil {
				log.Warnf("String substitute error: %s", err)
				return input, err
			}
			log.Debugf("  Replacement value: %s", replace)
			replacements[match[0]] = replace
		}
	}
	for search, replace := range replacements {
		input = strings.Replace(input, search, replace, -1)
	}
	return input, nil
}

/*
 * Return all matches based on search path and partial key name
 */
func (c *ConfMgr) ExistingKeys(key string, wantedType int, scope map[string]string, b backend.ConfigBackend) []string {
	foundKeys := make([]string, 0)

	for _, path := range c.SearchPaths(scope) {
		keyName := fmt.Sprintf("%s%s:%s", c.Config.Main.KeyPrefix, path, key)
		log.Debugf("Searching key: '%s'", keyName)
		keytype, _ := b.GetType(keyName)

		if keytype == wantedType {
			log.Debugf("Using key: %s", keyName)
			foundKeys = append(foundKeys, keyName)
		}
	}

	return foundKeys
}

/**
 * Returns the paths in reverse order as this is how all other functions
 * will consume it
 **/
func (c *ConfMgr) SearchPaths(reqscope map[string]string) []string {
	log.Debugf("SearchPaths scope: %q\n", reqscope)

	newKeyPaths := make([]string, 0)

	// Find all tokens in the configured paths
	token_re := regexp.MustCompile("%{(\\S+?)}")

PathsLoop:
	for _, path := range c.Config.Main.KeyPaths {
		log.Debugf("Config key path: '%s'", path)
		matches := token_re.FindAllStringSubmatch(path, -1)
		if len(matches) > 0 {
			for _, match := range matches {
				token_str := match[0]
				token_name := match[1]

				// Check if this token is defined in the request scope
				if val, ok := reqscope[token_name]; ok {
					log.Debugf("   Token '%s' is defined in the request as '%s'", token_name, val)
					path = strings.Replace(path, token_str, val, 1)
					log.Debugf("    Modified path: %s", path)
				} else {
					// Cannot replace this token, ignore this path completely
					log.Debugf("  Ignoring path %s because token %s is not set", path, token_name)
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
