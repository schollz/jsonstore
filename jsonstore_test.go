package jsonstore

import (
	"io/ioutil"
	"os"
	"regexp"
	"testing"
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
	if len(ks.GetAll(regexp.MustCompile(`hello`))) != 3 {
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
