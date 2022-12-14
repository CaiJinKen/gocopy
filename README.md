# gocopy

`NewFrom` is a deep copy func by `reflect`, and wihtout third-party dependencies. Supported types as flow:

- [x] signed/unsigned int
- [x] float
- [x] string
- [x] bool
- [x] map
- [x] slice / array
- [x] ptr
- [x] complex
- [x] struct
- [x] interface
- [x] chan
- [x] func
- [x] unsafePointer

**Note:** The func `NewFrom` returns itself when input is a chan or func or unsafePointer.

`Update` is a update tool by `reflect`, assign value from src to dst when src and dst has same field name and data type

### Install

```bash
go get github.com/CaiJinKen/gocopy
```

### Example

**NewFrom**:

```go
package main

import (
  "fmt"
  
  "github.com/CaiJinKen/gocopy"
)

type Config struct {
	Id    uint32
	Name  string
	Used  bool
	Pets  []*Dog
	PetMp map[string]*Dog
}

type Dog struct {
	Name string
	Age  uint8
}

func main() {
	conf := Config{
		Id:    10,
		Name:  "tiger",
		Used:  true,
		Pets:  []*Dog{{Name: "k1", Age: 1}, {Name: "k2", Age: 2}},
		PetMp: map[string]*Dog{"v1": {Name: "v1", Age: 3}, "v2": {Name: "v2", Age: 4}},
	}

	data := gocopy.NewFrom(conf) // data=>main.Config
	data = gocopy.NewFrom(&conf) // data=>&main.Config
	// data will intact when conf changed
	
	// handle data
	fmt.Printf("data %v\n",data)
}
```

**Update**:

```go
package main

import (
  "fmt"
  
  "github.com/CaiJinKen/gocopy"
)

type Config struct {
	Id    uint32
	Name  string
	Used  bool
	Pets  []*Dog
	PetMp map[string]*Dog
}

type Dog struct {
	Name string
	Age  uint8
}

func main() {
	conf := Config{
		Id:    10,
		Name:  "tiger",
		Used:  true,
		Pets:  []*Dog{{Name: "k1", Age: 1}, {Name: "k2", Age: 2}},
		PetMp: map[string]*Dog{"v1": {Name: "v1", Age: 3}, "v2": {Name: "v2", Age: 4}},
	}
  
  data := &Config

	err := gocopy.Update(conf, data)
	err = gocopy.Update(&conf, data)
	// dst change when conf reference data changed
	
	// handle data
	fmt.Printf("data %v\n",data)
}
```

