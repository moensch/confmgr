package confmgr

import (
	"encoding/json"
	"github.com/moensch/confmgr"
	"github.com/moensch/confmgr/vars"
	"testing"
)

func TestExisting(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	res := srv.ExistingKeys("hash", vars.TYPE_HASH, make(map[string]string))
	jsonblob, _ := json.MarshalIndent(res, "", "  ")
	t.Logf("%s\n", string(jsonblob))
}

func TestLookupString(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	res, err := srv.LookupString("string", make(map[string]string))
	if err != nil {
		t.Fatalf("ERROR: Cannot get hash: %s", err)
	}
	jsonblob, _ := json.MarshalIndent(res, "", "  ")
	t.Logf("%s\n", string(jsonblob))
}
func TestLookupHash(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	res, err := srv.LookupHash("hash", make(map[string]string))
	if err != nil {
		t.Fatalf("ERROR: Cannot get hash: %s", err)
	}
	jsonblob, _ := json.MarshalIndent(res, "", "  ")
	t.Logf("%s\n", string(jsonblob))
}

func TestLookupList(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	res, err := srv.LookupList("array", make(map[string]string))
	if err != nil {
		t.Fatalf("ERROR: Cannot get hash: %s", err)
	}
	jsonblob, _ := json.MarshalIndent(res, "", "  ")
	t.Logf("%s\n", string(jsonblob))
}
