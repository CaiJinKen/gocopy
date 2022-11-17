package gocopy

import (
	"testing"
)

var _src = Config{
	ID:        1,
	Hosts:     []string{"1,2,3"},
	AccessKey: "key1",
	SecretKey: "key2",
	Extra:     map[string]*Dog{"1": {Name: "dog1"}, "2": {Name: "dog2"}},
	Pets:      []*Dog{{Name: "dog3"}, {Name: "dog4"}},
}

type Config struct {
	ID        int
	Hosts     []string
	AccessKey string
	SecretKey string
	Extra     map[string]*Dog
	Pets      []*Dog
}
type Dog struct {
	Name string
}

func (c Config) Equal(v Config) bool {
	if c.ID != v.ID || c.AccessKey != v.AccessKey || c.SecretKey != v.SecretKey ||
		len(c.Hosts) != len(v.Hosts) || len(c.Extra) != len(v.Extra) {
		return false
	}
	for i := range c.Hosts {
		if c.Hosts[i] != v.Hosts[i] {
			return false
		}
	}
	for key, value := range c.Extra {
		if v.Extra[key] != value {
			return false
		}
	}

	return true
}

func TestNewFrom(t *testing.T) {
	value := NewFrom(_src)

	dst, ok := value.(Config)
	if !ok {
		t.Error("NewFrom copy err")
	}
	if _src.Equal(dst) {
		t.Error("NewFrom copy err")
	}
	_src.Hosts = append(_src.Hosts, "", "", "", "", "")
	if len(_src.Hosts) == len(dst.Hosts) || _src.Pets[0] == dst.Pets[0] || _src.Extra["1"] == dst.Extra["1"] {
		t.Error("NewFrom copy err")
	}
}

func TestUpdate(t *testing.T) {
	dst := &Config{}
	if err := Update(_src, dst); err != nil {
		t.Errorf("Update err:%+v", err)
	}
	if !_src.Equal(*dst) {
		t.Error("Update err")
	}

	dst = &Config{}
	if err := Update(&_src, dst); err != nil {
		t.Errorf("Update err:%+v", err)
	}
	if !_src.Equal(*dst) {
		t.Error("Update err")
	}

}
