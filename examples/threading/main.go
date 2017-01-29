package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"runtime"

	"github.com/schollz/jsonstore"
)

var fs jsonstore.JSONStore

// Here’s the worker, of which we’ll run several concurrent instances. These workers will receive work on the jobs channel and send the corresponding results on results. We’ll sleep a second per job to simulate an expensive task.
func worker(id int, jobs <-chan string, results chan<- int) {
	for j := range jobs {
		fs.SetMem(j, get(j))
		results <- 1
	}
}

func get(site string) string {
	response, err := http.Get(site)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		htmlData, _ := ioutil.ReadAll(response.Body) //<--- here!
		return string(htmlData)
	}
	return ""
}

func get_n(sites []string) {
	fs.Load("data.json")
	for _, site := range sites {
		fs.SetMem(site, get(site))
	}
	fs.Save()
}

func get_parallel(sites []string) {
	fs.Load("data.json")
	runtime.GOMAXPROCS(runtime.NumCPU())
	// In order to use our pool of workers we need to send them work and collect their results. We make 2 channels for this.
	jobs := make(chan string, 100)
	results := make(chan int, 100)
	// This starts up 3 workers, initially blocked because there are no jobs yet.
	for w := 1; w <= runtime.NumCPU(); w++ {
		go worker(w, jobs, results)
	}
	// Here we send n jobs and then close that channel to indicate that’s all the work we have.
	for _, site := range sites {
		jobs <- site
	}
	close(jobs)
	// Finally we collect all the results of the work.
	for a := 1; a <= len(sites); a++ {
		<-results
	}
	fs.Save()
}
