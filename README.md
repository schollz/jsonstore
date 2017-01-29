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

```golang
fs.Load("humans.json")

// set a key to any object you want
type Human struct {
  Name   string
  Height float64
}
fs.Set("human:1", Human{"Dante", 5.4})

// get the data back via an interface
var human Human
err := fs.Get("human:1", &human)
if err != nil {
  fmt.Println(err)
}
```

The datastore on disk is automatically Gziped (which can be toggled), and results in `humans.json.gz`:

```bash
$ zcat humans.json.gz
{
"human:1": "{\"Name\":\"Dante\",\"Height\":5.4}"
}
```

**JSONStore** in the wild:

- [schollz/urls](https://github.com/schollz/urls) - URL shortening

# License

MIT
