package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var concurrency int

var sites = []string{
	"http://trevorknott.io",
	"http://trevorknott.herokuapp.com",
	"https://evening-dusk-42226.herokuapp.com/posts",
	"https://www.google.com",
	"https://www.github.com",
	"https://www.golang.org",
	"http://www.twitter.com",
	"http://www.cloudflare.com",
	"http://www.flightstats.com",
	"http://www.linkedin.com/",
	"http://www.moz.com",
	"http://www.apple.com",
	"http://www.adobe.com",
	"http://www.blogspot.com",
	"http://www.tumblr.com",
	"http://www.youtube.com",
	"http://www.yahoo.com",
	"http://www.bing.com",
	"http://ec2-54-153-114-109.us-west-1.compute.amazonaws.com",
	"http://www.digg.com",
	"http://www.forbes.com",
	"http://www.dropbox.com",
	"http://www.yelp.com",
	"http://www.aol.com",
	"http://www.udemy.com",
	"http://www.webs.com",
	"http://www.egghead.io",
}

func init() {
	flag.IntVar(&concurrency, "con", 1, "concurrency flag. 1 is true")
	flag.Parse()
}

func main() {
	if concurrency != 1 {
		/* Non Concurrent */
		get()
	}
	if concurrency == 1 {
		/*  Concurrent  */
		concurrentlyGet()
	}
}

func get() {
	for _, s := range sites {
		res, _ := http.Get(s)
		fmt.Printf("%s %d\n", s, res.StatusCode)
		res.Body.Close()
	}
}

func concurrentlyGet() {
	log.Println(`Running w/ Concurrency.`)
	ch := make(chan string)

	for _, s := range sites {
		/* Concurrently call */
		go func(s string) {
			res, _ := http.Get(s)
			ch <- fmt.Sprintf("%s %d", s, res.StatusCode)
			res.Body.Close()
		}(s)
	}

	for range sites {
		fmt.Println(<-ch)
	}
}
