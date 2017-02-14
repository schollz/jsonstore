package jsonstore

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
)

// NoSuchKeyError is thrown when calling Get with invalid key
type NoSuchKeyError struct {
	key string
}

func (err NoSuchKeyError) Error() string {
	return "jsonstore: no such key \"" + err.key + "\""
}

// JSONStore is the basic store object.
type JSONStore struct {
	Data map[string]json.RawMessage
	sync.RWMutex
}

// Open will load a jsonstore from a file.
func Open(filename string) (*JSONStore, error) {
	var err error
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	var w io.Reader
	if strings.HasSuffix(filename, ".gz") {
		w, err = gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
	} else {
		w = f
	}

	toOpen := make(map[string]string)
	err = json.NewDecoder(w).Decode(&toOpen)
	if err != nil {
		return nil, err
	}

	ks := new(JSONStore)
	ks.Data = make(map[string]json.RawMessage)
	for key := range toOpen {
		ks.Data[key] = json.RawMessage(toOpen[key])
	}
	return ks, nil
}

// Save writes the jsonstore to disk
func Save(ks *JSONStore, filename string) (err error) {
	ks.RLock()
	defer ks.RUnlock()
	toSave := make(map[string]string)
	for key := range ks.Data {
		toSave[key] = string(ks.Data[key])
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if strings.HasSuffix(filename, ".gz") {
		w := gzip.NewWriter(f)
		defer w.Close()
		enc := json.NewEncoder(w)
		enc.SetIndent("", " ")
		return enc.Encode(toSave)
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", " ")
	return enc.Encode(toSave)
}

// Set saves a value at the given key.
func (s *JSONStore) Set(key string, value interface{}) error {
	s.Lock()
	defer s.Unlock()
	if s.Data == nil {
		s.Data = make(map[string]json.RawMessage)
	}
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	s.Data[key] = json.RawMessage(b)
	return nil
}

// Get will return the value associated with a key.
func (s *JSONStore) Get(key string, v interface{}) error {
	s.RLock()
	defer s.RUnlock()
	b, ok := s.Data[key]
	if !ok {
		return NoSuchKeyError{key}
	}
	return json.Unmarshal(b, &v)
}

// GetAll is like a filter with a regexp.
func (s *JSONStore) GetAll(re *regexp.Regexp) map[string]json.RawMessage {
	s.RLock()
	defer s.RUnlock()
	results := make(map[string]json.RawMessage)
	for k, v := range s.Data {
		if re.MatchString(k) {
			results[k] = v
		}
	}
	return results
}

// Keys returns all the keys currently in map
func (s *JSONStore) Keys() []string {
	s.RLock()
	defer s.RUnlock()
	keys := make([]string, len(s.Data))
	i := 0
	for k := range s.Data {
		keys[i] = k
		i++
	}
	return keys
}

// Delete removes a key from the store.
func (s *JSONStore) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.Data, key)
}
