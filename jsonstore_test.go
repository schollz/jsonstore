package jsonstore

import (
	"encoding/json"
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
	if world, ok := fs.Data["hello"]; !ok || string(world) != `"world"` {
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
	if err = Save(fs, f.Name()); err != nil {
		t.Error(err)
	}
	b, _ := json.Marshal(fs.Data)
	fs, _ = Load(f.Name())
	b2, _ := json.Marshal(fs.Data)
	if string(b) != string(b2) {
		t.Errorf("expected %s got %s", b, b2)
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
