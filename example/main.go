package main

import (
	"fmt"
	"sync"

	"github.com/schollz/jsonstore"
)

func main() {
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

	// Saving will automatically gzip if .gz is provided,
	// and can be performed in a wait group
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = jsonstore.Save(ks, "test.json.gz"); err != nil {
			panic(err)
		}
	}()
	wg.Wait()

	// Load any JSON / GZipped JSON
	ks2, err := jsonstore.Open("test.json.gz")
	if err != nil {
		panic(err)
	}

	// get the data back via an interface
	var human Human
	err = ks2.Get("human:1", &human)
	if err != nil {
		panic(err)
	}
	fmt.Println(human.Name)
}
