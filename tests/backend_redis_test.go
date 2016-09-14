package confmgr

import (
	"github.com/moensch/confmgr/backends/redis"
	"github.com/moensch/confmgr/vars"
	"testing"
)

func TestString(t *testing.T) {
	b := redis.Init()

	str, err := b.GetString("cfg:test:string")

	if err != nil {
		t.Fatalf("ERROR: Cannot get string: %s", err)
	}
	t.Logf("Retrieved string: '%s'", str)

	_, err = b.GetString("doesnotexist")
	if err == nil {
		t.Fatal("Expected error but none occurred")
	}
}

func TestType(t *testing.T) {
	b := redis.Init()

	testdata := make(map[string]int)
	testdata["cfg:test:string"] = vars.TYPE_STRING
	testdata["cfg:test:array"] = vars.TYPE_LIST
	testdata["cfg:test:hash"] = vars.TYPE_HASH
	testdata["notfound"] = vars.TYPE_NOT_FOUND

	for keyname, expected := range testdata {
		t.Logf("Testing type for key '%s'", keyname)
		actual, err := b.GetType(keyname)

		if err != nil {
			t.Fatalf("ERROR: Cannot check type: %s", err)
		}

		t.Logf("  Expected: %d", expected)
		t.Logf("  Actual  : %d", actual)

		if actual != expected {
			t.Fail()
		}
	}
}

func TestExists(t *testing.T) {
	b := redis.Init()

	exists, err := b.Exists("thiswillneverexist")

	if err != nil {
		t.Fatalf("ERROR: Cannot check exists: %s", err)
	}
	if exists == true {
		t.Log("key thiswillneverexist should return false on exists")
		t.Fail()
	} else {
		t.Log("absent key returned false")
	}

	exists, err = b.Exists("cfg:test:string")

	if err != nil {
		t.Fatalf("ERROR: Cannot check exists: %s", err)
	}
	if exists == false {
		t.Log("key cfg:test:string should return false on exists")
		t.Fail()
	} else {
		t.Log("present key returned true")
	}
}

func TestHashFieldExist(t *testing.T) {
	b := redis.Init()

	type TestEntry struct {
		Key    string
		Field  string
		Expect bool
	}

	testdata := []TestEntry{
		TestEntry{
			"cfg:test:hash",
			"field",
			true,
		},
		TestEntry{
			"cfg:test:hash",
			"noexist",
			false,
		},
		TestEntry{
			"invalidkey",
			"noexist",
			false,
		},
	}

	for idx, e := range testdata {
		t.Logf("Test %d: Testing key '%s', field '%s'", idx, e.Key, e.Field)

		exists, err := b.HashFieldExists(e.Key, e.Field)

		if err != nil {
			t.Fatalf("ERROR: Cannot get hash field: %s", err)
		}

		t.Logf("  Expected: '%t'", e.Expect)
		t.Logf("  Actual  : '%t'", exists)

		if exists != e.Expect {
			t.Fail()
		}
	}
}

func TestListIndexExist(t *testing.T) {
	b := redis.Init()

	type TestEntry struct {
		Key    string
		Index  int64
		Expect bool
	}

	testdata := []TestEntry{
		TestEntry{
			"cfg:test:array",
			0,
			true,
		},
		TestEntry{
			"cfg:test:array",
			1,
			true,
		},
		TestEntry{
			"cfg:test:array",
			2,
			true,
		},
		TestEntry{
			"cfg:test:array",
			3,
			false,
		},
		TestEntry{
			"cfg:test:array",
			-2,
			false,
		},
		TestEntry{
			"invalidkey",
			1,
			false,
		},
	}

	for idx, e := range testdata {
		t.Logf("Test %d: Testing key '%s', index '%d'", idx, e.Key, e.Index)

		exists, err := b.ListIndexExists(e.Key, e.Index)

		if err != nil {
			t.Fatalf("ERROR: Cannot get list index: %s", err)
		}

		t.Logf("  Expected: '%t'", e.Expect)
		t.Logf("  Actual  : '%t'", exists)

		if exists != e.Expect {
			t.Fail()
		}
	}
}

func TestHash(t *testing.T) {
	b := redis.Init()

	hash, err := b.GetHash("cfg:test:hash")

	if err != nil {
		t.Fatalf("ERROR: Cannot get hash: %s", err)
	}
	for hashKey, hashVal := range hash {
		t.Logf("%s => %s\n", hashKey, hashVal)
	}

}

func TestHashField(t *testing.T) {
	b := redis.Init()

	str, err := b.GetHashField("cfg:test:hash", "field")

	if err != nil {
		t.Fatalf("ERROR: Cannot get hash field: %s", err)
	}
	t.Logf("Retrieved hash field: '%s'", str)

}

func TestListIndex(t *testing.T) {
	b := redis.Init()

	str, err := b.GetListIndex("cfg:test:array", 0)

	if err != nil {
		t.Fatalf("ERROR: Cannot get list index: %s", err)
	}
	t.Logf("Retrieved list index: '%s'", str)
}

func TestList(t *testing.T) {
	b := redis.Init()

	strlist, err := b.GetList("cfg:test:array")

	if err != nil {
		t.Fatalf("ERROR: Cannot get list index: %s", err)
	}
	for pos, entry := range strlist {
		t.Logf(" List entry %d: '%s'", pos, entry)
	}
}

func TestListKeys(t *testing.T) {
	b := redis.Init()

	keys, err := b.ListKeys("")

	if err != nil {
		t.Fatalf("Cannot list keys: %s", err)
	}

	for pos, entry := range keys {
		t.Logf(" key %d: '%s'", pos, entry)
	}
}

func TestListKeysFilter(t *testing.T) {
	b := redis.Init()

	keys, err := b.ListKeys("*test*")

	if err != nil {
		t.Fatalf("Cannot list keys: %s", err)
	}

	for pos, entry := range keys {
		t.Logf(" key %d: '%s'", pos, entry)
	}
}
