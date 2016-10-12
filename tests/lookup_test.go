package confmgr

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/moensch/confmgr"
	"github.com/moensch/confmgr/backends/redis"
	"github.com/moensch/confmgr/vars"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
	b = redis.ConfigBackendRedis{}
	b.Conn, _ = redigo.Dial("tcp", ":6379")
}

func TestExisting(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	res := srv.ExistingKeys("hash", vars.TYPE_HASH, make(map[string]string), b)
	jsonblob, _ := json.MarshalIndent(res, "", "  ")
	t.Logf("%s\n", string(jsonblob))
}

func TestLookupString(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	res, err := srv.LookupString("string", make(map[string]string), b)
	if err != nil {
		t.Fatalf("ERROR: Cannot get hash: %s", err)
	}
	jsonblob, _ := json.MarshalIndent(res, "", "  ")
	t.Logf("%s\n", string(jsonblob))
}
func TestLookupHash(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	res, err := srv.LookupHash("hash", make(map[string]string), b)
	if err != nil {
		t.Fatalf("ERROR: Cannot get hash: %s", err)
	}
	jsonblob, _ := json.MarshalIndent(res, "", "  ")
	t.Logf("%s\n", string(jsonblob))
}

func TestLookupHashField(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	res, err := srv.LookupHashField("hash", "field1", make(map[string]string), b)
	if err != nil {
		t.Fatalf("ERROR: Cannot get hash field: %s", err)
	}
	jsonblob, _ := json.MarshalIndent(res, "", "  ")
	t.Logf("%s\n", string(jsonblob))
}

func TestLookupList(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	res, err := srv.LookupList("array", make(map[string]string), b)
	if err != nil {
		t.Fatalf("ERROR: Cannot get hash: %s", err)
	}
	jsonblob, _ := json.MarshalIndent(res, "", "  ")
	t.Logf("%s\n", string(jsonblob))
}

func TestSubstituteHashField(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	expected := "myvalue"

	res, err := srv.LookupHashField("otherhash", "simple", make(map[string]string), b)
	if err != nil {
		t.Fatalf("ERROR: Cannot get hash: %s", err)
	}

	actual := res.ToString()
	if actual != expected {
		t.Fatalf("Fail: Returned string %s did not match expected %s", actual, expected)
	}
}

func TestSubstituteHashFieldMulti(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	expected := "hello myvalue world myvalue2 goodbye entry2 and testing"
	res, err := srv.LookupHashField("otherhash", "multi", make(map[string]string), b)
	if err != nil {
		t.Fatalf("ERROR: Cannot get hash: %s", err)
	}

	actual := res.ToString()
	if actual != expected {
		t.Fatalf("Fail: Returned string %s did not match expected %s", actual, expected)
	}
}

func TestSubstituteStringRecurse(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	expected := "myvalue"

	res, err := srv.LookupString("recurse", make(map[string]string), b)
	if err != nil {
		t.Fatalf("ERROR: Cannot get hash: %s", err)
	}

	actual := res.ToString()
	if actual != expected {
		t.Fatalf("Fail: Returned string %s did not match expected %s", actual, expected)
	}
}

func TestSubstituteListIndexNotFound(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	_, err := srv.LookupString("otherstring", make(map[string]string), b)
	if err == nil {
		t.Fatalf("ERROR: Cannot get string: %s", err)
		t.Fatalf("Expected error but no error was generated")
	}
}

func TestSubstituteHashFieldNotFound(t *testing.T) {
	srv, _ := confmgr.NewConfMgr()

	_, err := srv.LookupString("fieldnotfound", make(map[string]string), b)
	if err == nil {
		t.Fatalf("Expected error but no error was generated")
	}
}
