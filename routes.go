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
			"HandleAdminListKeys",
			"GET",
			"/admin/keys",
			c.HandleAdminListKeys,
		},
		Route{
			"HandleAdminListKeysFiltered",
			"GET",
			"/admin/keys/{filter}",
			c.HandleAdminListKeysFiltered,
		},
		Route{
			"HandleAdminListHashFields",
			"GET",
			"/admin/util/hashfields/{keyName}",
			c.HandleAdminListHashFields,
		},
		Route{
			"HandleAdminGetKeyType",
			"GET",
			"/admin/util/type/{keyName}",
			c.HandleAdminGetKeyType,
		},
		Route{
			"HandleAdminGetKey",
			"GET",
			"/admin/key/{keyName}",
			c.HandleAdminKeyGet,
		},
		Route{
			"HandleAdminKeyStore",
			"POST",
			"/admin/key/{keyName}",
			c.HandleAdminKeyStore,
		},
		Route{
			"HandleAdminKeyDelete",
			"DELETE",
			"/admin/key/{keyName}",
			c.HandleAdminKeyDelete,
		},
		Route{
			"HandleAdminGetHashField",
			"GET",
			"/admin/key/{keyName}/{fieldName}",
			c.HandleAdminGetHashField,
		},
		Route{
			"HandleAdminSetHashField",
			"POST",
			"/admin/key/{keyName}/{fieldName}",
			c.HandleAdminSetHashField,
		},
		Route{
			"HandleAdminGetListIndex",
			"GET",
			"/admin/key/{keyName}/index/{listIndex:[0-9]+}",
			c.HandleAdminGetListIndex,
		},
		Route{
			"HandleLookupHash",
			"GET",
			"/hash/{keyName}",
			c.HandleLookupHash,
		},
		Route{
			"HandleLookupString",
			"GET",
			"/string/{keyName}",
			c.HandleLookupString,
		},
		Route{
			"HandleLookupList",
			"GET",
			"/list/{keyName}",
			c.HandleLookupList,
		},
		Route{
			"HandleLookupHashField",
			"GET",
			"/string/{keyName}/{fieldName}",
			c.HandleLookupHashField,
		},
		Route{
			"HandleLookupListIndex",
			"GET",
			"/string/{keyName}/index/{listIndex}",
			c.HandleLookupListIndex,
		},
	}
}
