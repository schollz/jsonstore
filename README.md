# jsonstore  :convenience_store:

[![GoDoc](https://godoc.org/github.com/schollz/jsonstore?status.svg)](https://godoc.org/github.com/schollz/jsonstore)

*JSONStore* is a Go-library for a simple thread-safe in-memory JSON key-store with persistent backend. It's made for those times where you don't need a RDBMS like [MySQL](https://www.mysql.com/), or a NoSQL like [MongoDB](https://www.mongodb.com/) - basically when you just need a simple keystore. A really simple keystore. *JSONStore* is used in those times you don't need a distributed keystore like [etcd](https://coreos.com/etcd/docs/latest/), or
a remote keystore [Redis](https://redis.io/) or a local keystore like [Bolt](https://github.com/boltdb/bolt). Its really for those times where you just need a JSON file.

## Usage

First, install the library using:

```
go get -u -v github.com/schollz/jsonstore
```

Then you can add it to your program. Check out the examples, or see below for basic usage:

```golang
ks := new(jsonstore.JSONStore)

// set a key to any object you want
type Human struct {
  Name   string
  Height float64
}
err := ks.Set("human:1", Human{"Dante", 5.4})
if err != nil {
  panic(err)
}

// Saving will automatically gzip if .gz is provided
if err = jsonstore.Save(ks, "humans.json.gz"); err != nil {
  panic(err)
}

// Load any JSON / GZipped JSON
ks2, err := jsonstore.Open("humans.json.gz")
if err != nil {
  panic(err)
}

// get the data back via an interface
var human Human
err = ks2.Get("human:1", &human)
if err != nil {
  panic(err)
}
fmt.Println(human.Name) // Prints 'Dante'
```

The datastore on disk is then contains:

```bash
$ zcat humans.json.gz
{
"human:1": "{\"Name\":\"Dante\",\"Height\":5.4}"
}
```


**JSONStore** in the wild:

- [schollz/urls](https://github.com/schollz/urls) - URL shortening

# Dev

Benchmark against using Redis and BoltDB as KeyStores using Go1.8 (Intel i5-4310U CPU @ 2.00GHz). Take away is that setting/getting is faster in *JSONStore* (because its just a map), but opening is much slower (because its a file that is read into memory). So don't use this if you have to store 1,000,000+ things!

```
$ go test -bench=. tests/redis/* > redis.txt
$ go test -bench=. tests/bolt/* > bolt.txt
$ go test -bench=. > jsonstore.txt
$ benchcmp bolt.txt jsonstore.txt
benchmark                old ns/op     new ns/op     delta
BenchmarkSet-4           5471747       1847          -99.97%
BenchmarkGet-4           2424          1479          -38.99%
BenchmarkOpen100-4       11168         148035        +1225.53%
BenchmarkOpen10000-4     10095         19722376      +195267.77%
$ benchcmp redis.txt jsonstore.txt
benchmark                old ns/op     new ns/op     delta
BenchmarkSet-4           24717         1847          -92.53%
BenchmarkGet-4           22561         1479          -93.44%
BenchmarkOpen100-4       6221          148035        +2279.60%
BenchmarkOpen10000-4     4951          19722376      +398251.36%
```

# License

MIT
