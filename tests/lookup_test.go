package confmgr

import (
	"encoding/json"
	"github.com/moensch/confmgr"
	"github.com/moensch/confmgr/vars"
	"testing"
)

func TestExisting(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	res := srv.ExistingKeys("db_policy", vars.TYPE_HASH)
	jsonblob, _ := json.MarshalIndent(res, "", "  ")
	t.Logf("%s\n", string(jsonblob))
}

/*
func TestLookupString(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	res := srv.LookupString("string")
	jsonblob, _ := json.MarshalIndent(res, "", "  ")
	t.Logf("%s\n", string(jsonblob))
}
func TestLookupKeyHash(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	res, _ := srv.LookupKey("db_policy")
	jsonblob, _ := json.MarshalIndent(res, "", "  ")
	t.Logf("%s\n", string(jsonblob))
}

func TestLookupKeyString(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	res, _ := srv.LookupKey("string")
	jsonblob, _ := json.MarshalIndent(res, "", "  ")
	t.Logf("%s\n", string(jsonblob))
}
*/
