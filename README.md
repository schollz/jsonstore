# jsonstore  :convenience_store:

[![GoDoc](https://godoc.org/github.com/schollz/jsonstore?status.svg)](https://godoc.org/github.com/schollz/jsonstore)

*JSONStore* is a Go-library for a simple thread-safe in-memory JSON key-store with persistent backend.
It's made for those times where you don't need a RDBMS like [MySQL](https://www.mysql.com/),
or a NoSQL like [MongoDB](https://www.mongodb.com/) - basically when you just need a simple keystore.
A really simple keystore. *JSONStore* is used in those times you don't need a distributed keystore
like [etcd](https://coreos.com/etcd/docs/latest/), or
a remote keystore [Redis](https://redis.io/) or a local keystore like [Bolt](https://github.com/boltdb/bolt).
Its really for those times where you just need a JSON file, hence *JSONStore*.

## Usage

Here's an example:

```golang
package main

import (
	"fmt"

	"github.com/schollz/jsonstore"
)

var fs jsonstore.JSONStore

func main() {
	// initialize data file
	fs.Init()

	// set a key to any object you want
	type Human struct {
		Name   string
		Height float64
	}
	fs.Set("human:1", Human{"Dante", 5.4})

	// get the data back via an interface
	get, err := fs.Get("human:1")
	if err != nil {
		fmt.Println(err)
	}
	// convert the object from interface
	fmt.Println(get.(Human))

	// get the data of a non-existent object
	_, err = fs.Get("nothing") // throws 'not found' error
	if err != nil {
		fmt.Println(err)
	}

	// get data from a lot of objects
	fs.Set("human:2", Human{"Da Vinci", 5.2})
	fs.Set("human:3", Human{"Einstein", 5.43})
	fs.Set("NumberOfHumans", "3")
	get, err = fs.Get("human:*")
	for key, val := range get.(map[string]interface{}) {
		fmt.Println(key, val.(Human))
	}
}
```

It will automatically save it to a file with Gzip compression, `data.json.gz`.
You can see your JSONStore file easily,

```
$ zcat data.json.gz
{
 "NumberOfHumans": "3",
 "human:1": {
  "Name": "Dante",
  "Height": 5.4
 },
 "human:2": {
  "Name": "Da Vinci",
  "Height": 5.2
 },
 "human:3": {
  "Name": "Einstein",
  "Height": 5.43
 }
}
```

**JSONStore** in the wild:

- [schollz/urls](https://github.com/schollz/urls) - URL shortening

# License

MIT
