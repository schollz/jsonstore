package jsonstore

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type JSONStore struct {
	Data     map[string]interface{}
	location string
	gzip     bool
	sync.RWMutex
}

func (s *JSONStore) Init() {
	s.Lock()
	defer s.Unlock()
	s.location = "data.json.gz"
	s.Data = make(map[string]interface{})
	s.gzip = true
}

func (s *JSONStore) SetGzip(on bool) {
	s.Lock()
	defer s.Unlock()
	s.gzip = on
}

func (s *JSONStore) SetLocation(location string) {
	s.Lock()
	s.location = location
	s.Unlock()
}

func (s *JSONStore) Load() error {
	s.Lock()
	defer s.Unlock()
	var err error
	if _, err = os.Stat(s.location); os.IsNotExist(err) {
		err = errors.New("Location does not exist")
	} else {
		var b []byte
		if s.gzip {
			if !strings.Contains(s.location, ".gz") {
				s.location = s.location + ".gz"
			}
			b, err = readGzFile(s.location)
			if err != nil {
				return err
			}
		} else {
			b, err = ioutil.ReadFile(s.location)
			if err != nil {
				return err
			}
		}
		err = json.Unmarshal(b, &s.Data)
	}
	return err
}

func (s *JSONStore) Save() error {
	s.RLock()
	defer s.RUnlock()
	var err error
	b, err := json.MarshalIndent(s.Data, "", " ")
	if err != nil {
		return err
	}

	if s.gzip {
		var b2 bytes.Buffer
		w := gzip.NewWriter(&b2)
		w.Write(b)
		w.Close()
		err = ioutil.WriteFile(s.location, b2.Bytes(), 0644)
	} else {
		err = ioutil.WriteFile(s.location, b, 0644)
	}
	return err
}

func (s *JSONStore) Set(key string, value interface{}) error {
	s.set(key, value)
	s.Save()
	return nil
}

func (s *JSONStore) SetMem(key string, value interface{}) error {
	s.set(key, value)
	return nil
}

func (s *JSONStore) set(key string, value interface{}) error {
	s.Lock()
	defer s.Unlock()
	s.Data[key] = value
	return nil
}

func (s *JSONStore) Get(key string) (interface{}, error) {
	if strings.Contains(key, "*") {
		return s.getmany(key)
	} else {
		return s.getone(key)
	}
}

func (s *JSONStore) getmany(key string) (interface{}, error) {
	s.RLock()
	defer s.RUnlock()
	possible := []string{}
	for _, substring := range strings.Split(key, "*") {
		if strings.Contains(substring, "*") || len(substring) == 0 {
			continue
		}
		possible = append(possible, substring)
	}

	m := make(map[string]interface{})
	for key := range s.Data {
		for _, substring := range possible {
			if strings.Contains(key, substring) {
				m[key] = s.Data[key]
			}
		}
	}

	if len(m) == 0 {
		return -1, errors.New(key + " not found")
	}
	return m, nil
}

func (s *JSONStore) getone(key string) (interface{}, error) {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.Data[key]
	if !ok {
		return -1, errors.New(key + " not found")
	}
	return val, nil
}

// utils

// from http://stackoverflow.com/questions/16890648/how-can-i-use-golangs-compress-gzip-package-to-gzip-a-file
func readGzFile(filename string) ([]byte, error) {
	fi, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	if err != nil {
		return nil, err
	}
	defer fz.Close()

	s, err := ioutil.ReadAll(fz)
	if err != nil {
		return nil, err
	}
	return s, nil
}
