# jsonstore  :convenience_store:

[![GoDoc](https://godoc.org/github.com/schollz/jsonstore?status.svg)](https://godoc.org/github.com/schollz/jsonstore)

*JSONStore* is a Go-library for a simple thread-safe in-memory JSON key-store with persistent backend.
It's made for those times where you don't need a RDBMS like [MySQL](https://www.mysql.com/),
or a NoSQL like [MongoDB](https://www.mongodb.com/) - basically when you just need a simple keystore.
A really simple keystore. *JSONStore* is used in those times you don't need a distributed keystore
like [etcd](https://coreos.com/etcd/docs/latest/), or
a remote keystore [Redis](https://redis.io/) or a local keystore like [Bolt](https://github.com/boltdb/bolt).
Its really for those times where you just need a JSON file, hence *JSONStore*.

Its very easy to use:

```golang
fs.Init()                 // initialize data file
fs.Set("data", 1234)      // sets a key to 1234 and saves it
data, _ := fs.Get("data")
fmt.Println(data)         // prints 1234
_, err := fs.Get("nothing")
fmt.Println(err)          // throws 'not found' error
```

It will automatically save it to a file with Gzip compression, `data.json.gz`.
You can see your JSONStore file easily,

```
$ zcat data.json.gz
{
  "data":1234
}
```

# License

MIT
