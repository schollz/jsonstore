package jsonstore

import (
	"io/ioutil"
	"os"
	"testing"
)

func testFile() *os.File {
	f, err := ioutil.TempFile(".", "jsonstore")
	if err != nil {
		panic(err)
	}
	return f
}

func TestLoad(t *testing.T) {
	f := testFile()
	defer os.Remove(f.Name())
	ioutil.WriteFile(f.Name(), []byte(`{"hello":"world"}`), 0644)
	fs, err := Load(f.Name())
	if err != nil {
		t.Error(err)
	}
	if len(fs.Data) != 1 {
		t.Errorf("expected %d got %d", 1, len(fs.Data))
	}
	if world, ok := fs.Data["hello"]; !ok || string(world) != `world` {
		t.Errorf("expected %s got %s", "world", world)
	}
}

func TestGeneral(t *testing.T) {
	f := testFile()
	defer os.Remove(f.Name())
	fs := new(JSONStore)
	err := fs.Set("hello", "world")
	if err != nil {
		t.Error(err)
	}
	if err = Save(fs, "test.json"); err != nil {
		t.Error(err)
	}

	fs2, _ := Load("test.json")
	var a string
	var b string
	fs.Get("hello", &a)
	fs2.Get("hello", &b)
	if a != b {
		t.Errorf("expected '%s' got '%s'", a, b)
	}

	// Set a object, using a Gzipped JSON
	type Human struct {
		Name   string
		Height float64
	}
	fs.Set("human:1", Human{"Dante", 5.4})
	Save(fs, "test2.json.gz")
	fs2, _ = Load("test2.json.gz")
	var human Human
	fs2.Get("human:1", &human)
	if human.Height != 5.4 {
		t.Errorf("expected '%v', got '%v'", Human{"Dante", 5.4}, human)
	}
}

func BenchmarkGet(b *testing.B) {
	fs := new(JSONStore)
	fs.Set("data", 1234)
	fs.Set("name", "bob")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var num int
		fs.Get("data", &num)
	}
}
