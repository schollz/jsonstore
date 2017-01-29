package jsonstore

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type JSONStore struct {
	Data        map[string][]byte
	location    string
	initialized bool
	gzip        bool
	sync.RWMutex
}

// SetGzip will toggle Gzip compression
func (s *JSONStore) SetGzip(on bool) {
	s.Lock()
	defer s.Unlock()
	s.gzip = on
	if s.gzip && !strings.Contains(s.location, ".gz") {
		s.location = s.location + ".gz"
	} else if !s.gzip && strings.Contains(s.location, ".gz") {
		s.location = strings.Replace(s.location, ".gz", "", 1)
	}
}

// Load will load the data from the current file
func (s *JSONStore) Load(location string) error {
	s.Lock()
	defer s.Unlock()
	s.gzip = true
	s.location = location
	if s.gzip && !strings.Contains(s.location, ".gz") {
		s.location = s.location + ".gz"
	} else if !s.gzip && strings.Contains(s.location, ".gz") {
		s.location = strings.Replace(s.location, ".gz", "", 1)
	}

	s.Data = make(map[string][]byte)
	s.gzip = true

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
		data := make(map[string]string)
		err = json.Unmarshal(b, &data)
		if err != nil {
			return err
		}

		for key := range data {
			s.Data[key] = compressByte([]byte(data[key]))
		}
	}
	return nil
}

// Save will save the current data to the location, adding Gzip compression if
// is enabled (it is by default)
func (s *JSONStore) Save() error {
	s.RLock()
	defer s.RUnlock()
	var err error

	// Decompress data to save it so its readable on disk
	data := make(map[string]string)
	for key := range s.Data {
		data[key] = string(decompressByte(s.Data[key]))
	}

	b, err := json.MarshalIndent(data, "", " ")
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

// Set will set a key to a value, and then save go disk
func (s *JSONStore) Set(key string, value interface{}) error {
	s.set(key, value)
	s.Save()
	return nil
}

// SetMem will set a key to a value, but not save to disk
func (s *JSONStore) SetMem(key string, value interface{}) error {
	s.set(key, value)
	return nil
}

func (s *JSONStore) set(key string, value interface{}) error {
	s.Lock()
	defer s.Unlock()

	// Marshal the data into JSON
	bJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Compress into map
	s.Data[key] = compressByte(bJSON)

	return nil
}

// Get will return the value associated with a key
// if the key contains a `*`, like `name:*`, it will a map[string]interface{}
// where each key is a key containing `*` and its corresponding value
func (s *JSONStore) Get(key string, v interface{}) error {
	s.RLock()
	defer s.RUnlock()
	compressedVal, ok := s.Data[key]
	if !ok {
		return errors.New(key + " not found")
	}
	val := decompressByte(compressedVal)
	return json.Unmarshal(val, &v)
}

// func (s *JSONStore) getmany(key string, v interface{}) error {
// 	s.RLock()
// 	defer s.RUnlock()
// 	possible := []string{}
// 	for _, substring := range strings.Split(key, "*") {
// 		if strings.Contains(substring, "*") || len(substring) == 0 {
// 			continue
// 		}
// 		possible = append(possible, substring)
// 	}
//
// 	for key := range s.Data {
// 		for _, substring := range possible {
// 			if strings.Contains(key, substring) {
// 				v[key] = s.Data[key]
// 			}
// 		}
// 	}
//
// 	if len(m) == 0 {
// 		return errors.New(key + " not found")
// 	}
// 	return nil
// }

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

// compressByte returns a compressed byte slice.
func compressByte(src []byte) []byte {
	compressedData := new(bytes.Buffer)
	compress(src, compressedData, 9)
	return compressedData.Bytes()
}

// decompressByte returns a decompressed byte slice.
func decompressByte(src []byte) []byte {
	compressedData := bytes.NewBuffer(src)
	deCompressedData := new(bytes.Buffer)
	decompress(compressedData, deCompressedData)
	return deCompressedData.Bytes()
}

// compress uses flate to compress a byte slice to a corresponding level
func compress(src []byte, dest io.Writer, level int) {
	compressor, _ := flate.NewWriter(dest, level)
	compressor.Write(src)
	compressor.Close()
}

// compress uses flate to decompress an io.Reader
func decompress(src io.Reader, dest io.Writer) {
	decompressor := flate.NewReader(src)
	io.Copy(dest, decompressor)
	decompressor.Close()
}
