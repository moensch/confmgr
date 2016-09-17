package backend

type ConfigBackend interface {
	Exists(string) (bool, error)
	GetType(string) (int, error)
	DeleteKey(string) error
	GetHash(string) (map[string]string, error)
	SetHash(string, map[string]string) error
	GetHashField(string, string) (string, error)
	SetHashField(string, string, string) error
	HashFieldExists(string, string) (bool, error)
	GetList(string) ([]string, error)
	SetList(string, []string) error
	GetListIndex(string, int64) (string, error)
	ListIndexExists(string, int64) (bool, error)
	GetString(string) (string, error)
	SetString(string, string) error
	ListKeys(string) ([]string, error)
	ListAppend(string, string) error
	Check() error
	Close()
}

type ConfigBackendFactory interface {
	NewBackend() ConfigBackend
}
