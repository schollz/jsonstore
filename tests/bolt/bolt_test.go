package jsonstore

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/boltdb/bolt"
)

type Human struct {
	Name   string
	Height float64
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

func BenchmarkSet(b *testing.B) {
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
			bJSON, _ := json.Marshal(Human{"Dante", 5.4})
			err := b.Put([]byte("data"+strconv.Itoa(i)), bJSON)
			return err
		})
	}
}

func BenchmarkGet(b *testing.B) {
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
		bJSON, _ := json.Marshal(Human{"Dante", 5.4})
		err := b.Put([]byte("data"), bJSON)
		return err
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("MyBucket"))
			dat := b.Get([]byte("data"))
			var human Human
			json.Unmarshal([]byte(dat), &human)
			return nil
		})
	}
}

func BenchmarkOpen100(b *testing.B) {
	defer os.Remove("my.db")
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	for i := 1; i < 100; i++ {
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("MyBucket"))
			err := b.Put([]byte("hello:"+strconv.Itoa(i)), []byte("world"+strconv.Itoa(i)))
			return err
		})
	}
	db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db, err := bolt.Open("my.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		db.Close()
	}
}

func BenchmarkOpen10000(b *testing.B) {
	defer os.Remove("my.db")
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		for i := 1; i < 10000; i++ {
			err := b.Put([]byte("hello:"+strconv.Itoa(i)), []byte("world"+strconv.Itoa(i)))
			return err
		}
		return nil
	})
	db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db, err := bolt.Open("my.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		db.Close()
	}
}
