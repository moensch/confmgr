package confmgr

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func (c *ConfMgr) RouteDefinitions() Routes {
	return Routes{
		Route{
			"Index",
			"GET",
			"/",
			c.Index,
		},
		Route{
			"AdminListKeys",
			"GET",
			"/admin/keys",
			c.HandleListKeys,
		},
		Route{
			"AdminListKeysFiltered",
			"GET",
			"/admin/keys/{filter}",
			c.HandleListKeysFiltered,
		},
		Route{
			"AdminListHashFields",
			"GET",
			"/admin/util/hashfields/{keyName}",
			c.HandleListHashFields,
		},
		Route{
			"AdminGetKeyType",
			"GET",
			"/admin/util/type/{keyName}",
			c.HandleGetKeyType,
		},
		Route{
			"AdminGetKey",
			"GET",
			"/admin/key/{keyName}",
			c.KeyGet,
		},
		Route{
			"AdminKeyStore",
			"POST",
			"/admin/key/{keyName}",
			c.KeyStore,
		},
		Route{
			"AdminKeyDelete",
			"DELETE",
			"/admin/key/{keyName}",
			c.HandleKeyDelete,
		},
		Route{
			"AdminGetKeyHashField",
			"GET",
			"/admin/key/{keyName}/{fieldName}",
			c.KeyGetHashField,
		},
		Route{
			"AdminSetKeyHashField",
			"POST",
			"/admin/key/{keyName}/{fieldName}",
			c.HandleSetHashField,
		},
		Route{
			"AdminGetKeyListIndex",
			"GET",
			"/admin/key/{keyName}/index/{listIndex:[0-9]+}",
			c.KeyGetListIndex,
		},
		Route{
			"LookupHashKey",
			"GET",
			"/hash/{keyName}",
			c.HandleLookupHash,
		},
		Route{
			"LookupStringKey",
			"GET",
			"/string/{keyName}",
			c.HandleLookupString,
		},
		Route{
			"LookupListKey",
			"GET",
			"/list/{keyName}",
			c.HandleLookupList,
		},
		Route{
			"LookupHashField",
			"GET",
			"/string/{keyName}/{fieldName}",
			c.HandleLookupHashField,
		},
		Route{
			"LookupListIndex",
			"GET",
			"/string/{keyName}/index/{listIndex}",
			c.HandleLookupListIndex,
		},
	}
}
