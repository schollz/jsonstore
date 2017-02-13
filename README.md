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

Benchmark against Redis and BoltDB using Go1.8 (Intel i5-4310U CPU @ 2.00GHz):

```
BenchmarkGet-4        	 1000000	      1481 ns/op
BenchmarkSet-4        	 2000000	       850 ns/op
BenchmarkSave-4       	    5000	    346912 ns/op
BenchmarkRedisSet-4   	  100000	     20673 ns/op
BenchmarkRedisGet-4   	  100000	     22280 ns/op
BenchmarkBoltSet-4    	     300	   4413128 ns/op
BenchmarkBoltGet-4    	  500000	      2473 ns/op
PASS
ok  	github.com/schollz/jsonstore	13.853s
```

# License

MIT
