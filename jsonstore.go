package jsonstore

import (
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
	sync.RWMutex
}

func (s *JSONStore) Init() {
	s.location = "data.json"
	s.Data = make(map[string]interface{})
}

func (s *JSONStore) SetLocation(location string) {
	s.Lock()
	s.location = location
	s.Unlock()
}

func (s *JSONStore) Load() error {
	s.Lock()
	defer s.Unlock()
	if _, err := os.Stat(s.location); os.IsNotExist(err) {
		return errors.New("Location does not exist")
	} else {
		b, err2 := ioutil.ReadFile(s.location)
		if err != nil {
			return err
		}
		err2 = json.Unmarshal(b, &s.Data)
		return err2
	}
}

func (s *JSONStore) Save() error {
	s.Lock()
	defer s.Unlock()
	b, err := json.MarshalIndent(s.Data, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(s.location, b, 0644)
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
