package confmgr

import (
	"net/http"
)

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
			handlerDecorate(c.Index),
		},
		Route{
			"HandleAdminListKeys",
			"GET",
			"/admin/keys",
			handlerDecorate(c.HandleAdminListKeys),
		},
		Route{
			"HandleAdminListKeysFiltered",
			"GET",
			"/admin/keys/{filter}",
			handlerDecorate(c.HandleAdminListKeysFiltered),
		},
		Route{
			"HandleAdminListHashFields",
			"GET",
			"/admin/util/hashfields/{keyName}",
			handlerDecorate(c.HandleAdminListHashFields),
		},
		Route{
			"HandleAdminGetKeyType",
			"GET",
			"/admin/util/type/{keyName}",
			handlerDecorate(c.HandleAdminGetKeyType),
		},
		Route{
			"HandleAdminGetKey",
			"GET",
			"/admin/key/{keyName}",
			handlerDecorate(c.HandleAdminKeyGet),
		},
		Route{
			"HandleAdminKeyStore",
			"POST",
			"/admin/key/{keyName}",
			handlerDecorate(c.HandleAdminKeyStore),
		},
		Route{
			"HandleAdminKeyDelete",
			"DELETE",
			"/admin/key/{keyName}",
			handlerDecorate(c.HandleAdminKeyDelete),
		},
		Route{
			"HandleAdminGetHashField",
			"GET",
			"/admin/key/{keyName}/{fieldName}",
			handlerDecorate(c.HandleAdminGetHashField),
		},
		Route{
			"HandleAdminListAppend",
			"PATCH",
			"/admin/key/append/{keyName}",
			handlerDecorate(c.HandleAdminListAppend),
		},
		Route{
			"HandleAdminSetHashField",
			"POST",
			"/admin/key/{keyName}/{fieldName}",
			handlerDecorate(c.HandleAdminSetHashField),
		},
		Route{
			"HandleAdminGetListIndex",
			"GET",
			"/admin/key/{keyName}/index/{listIndex:[0-9]+}",
			handlerDecorate(c.HandleAdminGetListIndex),
		},
		Route{
			"HandleLookupHash",
			"GET",
			"/hash/{keyName}",
			handlerDecorate(c.HandleLookupHash),
		},
		Route{
			"HandleLookupString",
			"GET",
			"/string/{keyName}",
			handlerDecorate(c.HandleLookupString),
		},
		Route{
			"HandleLookupList",
			"GET",
			"/list/{keyName}",
			handlerDecorate(c.HandleLookupList),
		},
		Route{
			"HandleLookupHashField",
			"GET",
			"/string/{keyName}/{fieldName}",
			handlerDecorate(c.HandleLookupHashField),
		},
		Route{
			"HandleLookupListIndex",
			"GET",
			"/string/{keyName}/index/{listIndex}",
			handlerDecorate(c.HandleLookupListIndex),
		},
	}
}
