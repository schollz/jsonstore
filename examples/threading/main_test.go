package main

import "testing"

var sites = []string{"http://bettermotherfuckingwebsite.com/", "http://motherfuckingwebsite.com/", "https://example.org/", "http://icanhazip.com/", "http://www.howmanypeopleareinspacerightnow.com/"}

func TestGet(t *testing.T) {
	get_parallel(sites)
	var ip string
	err := fs.Get("http://icanhazip.com/", &ip)
	if err != nil {
		t.Error(err)
	}
	if len(ip) == 0 {
		t.Errorf("No IP found?")
	}
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		get_n(sites)
	}
}

func BenchmarkGetParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		get_parallel(sites)
	}
}
