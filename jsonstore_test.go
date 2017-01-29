package jsonstore

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/boltdb/bolt"

	redis "gopkg.in/redis.v5"
)

func testFile() *os.File {
	f, err := ioutil.TempFile(".", "jsonstore")
	if err != nil {
		panic(err)
	}
	return f
}

func TestOpen(t *testing.T) {
	f := testFile()
	defer os.Remove(f.Name())
	ioutil.WriteFile(f.Name(), []byte(`{"hello":"world"}`), 0644)
	ks, err := Open(f.Name())
	if err != nil {
		t.Error(err)
	}
	if len(ks.Data) != 1 {
		t.Errorf("expected %d got %d", 1, len(ks.Data))
	}
	if world, ok := ks.Data["hello"]; !ok || string(world) != `world` {
		t.Errorf("expected %s got %s", "world", world)
	}
}

func TestGeneral(t *testing.T) {
	f := testFile()
	defer os.Remove(f.Name())
	ks := new(JSONStore)
	err := ks.Set("hello", "world")
	if err != nil {
		t.Error(err)
	}
	if err = Save(ks, f.Name()); err != nil {
		t.Error(err)
	}

	ks2, _ := Open(f.Name())
	var a string
	var b string
	ks.Get("hello", &a)
	ks2.Get("hello", &b)
	if a != b {
		t.Errorf("expected '%s' got '%s'", a, b)
	}

	// Set a object, using a Gzipped JSON
	type Human struct {
		Name   string
		Height float64
	}
	ks.Set("human:1", Human{"Dante", 5.4})
	Save(ks, "test2.json.gz")
	ks2, _ = Open("test2.json.gz")
	var human Human
	ks2.Get("human:1", &human)
	if human.Height != 5.4 {
		t.Errorf("expected '%v', got '%v'", Human{"Dante", 5.4}, human)
	}
}

func TestRegex(t *testing.T) {
	f := testFile()
	defer os.Remove(f.Name())
	ks := new(JSONStore)
	ks.Set("hello:1", "world1")
	ks.Set("hello:2", "world2")
	ks.Set("hello:3", "world3")
	ks.Set("world:1", "hello1")
	if len(ks.GetAll(regexp.MustCompile(`hello`))) != len(ks.Keys())-1 {
		t.Errorf("Problem getting all")
	}
}

func BenchmarkGet(b *testing.B) {
	ks := new(JSONStore)
	ks.Set("data", 1234)
	ks.Set("name", "bob")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var num int
		ks.Get("data", &num)
	}
}

func BenchmarkSet(b *testing.B) {
	ks := new(JSONStore)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ks.Set("data", 1234)
	}
}

func BenchmarkSave(b *testing.B) {
	ks := new(JSONStore)
	ks.Set("data", 1234)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Save(ks, "benchmark.json.gz")
	}
}

func TestRedis(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := client.Set("data", 1234, 0).Err()
	if err != nil {
		t.Errorf(err.Error())
	}
	val, err := client.Get("data").Result()
	if err != nil {
		t.Errorf(err.Error())
	}
	if val != "1234" {
		t.Errorf("Got %v instead of 1234", val)
	}
}

func BenchmarkRedisSet(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := client.Set("key", "value", 0).Err()
	if err != nil {
		fmt.Println("Problem setting Redis, is it installed?")
		fmt.Println(err)
		os.Exit(-1)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Set("data", 1234, 0)
	}
}

func BenchmarkRedisGet(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := client.Set("data", 1234, 0).Err()
	if err != nil {
		fmt.Println("Problem setting Redis, is it installed?")
		fmt.Println(err)
		os.Exit(-1)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Get("data").Result()
	}
}

func TestBolt(t *testing.T) {
	defer os.Remove("my.db")
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		err := b.Put([]byte("data"), []byte("1234"))
		return err
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	var result string
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		result = string(b.Get([]byte("data")))
		return nil
	})

	if result != "1234" {
		t.Errorf("Problem reading/writing with BoltDB")
	}
}

func BenchmarkBoltSet(b *testing.B) {
	defer os.Remove("my.db")
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("MyBucket"))
			err := b.Put([]byte("data"), []byte("1234"))
			return err
		})
	}
}

func BenchmarkBoltGet(b *testing.B) {
	defer os.Remove("my.db")
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		err := b.Put([]byte("data"), []byte("1234"))
		return err
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("MyBucket"))
			b.Get([]byte("data"))
			return nil
		})
	}
}
