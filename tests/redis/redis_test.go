package jsonstore

import (
	"encoding/json"
	"strconv"
	"testing"

	redis "gopkg.in/redis.v5"
)

type Human struct {
	Name   string
	Height float64
}

func TestRedis(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()

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

func BenchmarkSet(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bJSON, _ := json.Marshal(Human{"Dante", 5.4})
		err := client.Set("human:"+strconv.Itoa(i), bJSON, 0).Err()
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkGet(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()
	bJSON, _ := json.Marshal(Human{"Dante", 5.4})
	err := client.Set("human:1", bJSON, 0).Err()
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, _ := client.Get("human:1").Result()
		var human Human
		json.Unmarshal([]byte(v), &human)
	}
}

func BenchmarkOpen100(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	for i := 1; i < 100; i++ {
		err := client.Set("hello:"+strconv.Itoa(i), "world"+strconv.Itoa(i), 0).Err()
		if err != nil {
			panic(err)
		}
	}
	client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		client.Close()
	}
}

func BenchmarkOpen10000(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	for i := 1; i < 10000; i++ {
		err := client.Set("hello:"+strconv.Itoa(i), "world"+strconv.Itoa(i), 0).Err()
		if err != nil {
			panic(err)
		}
	}
	client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		client.Close()
	}
}
