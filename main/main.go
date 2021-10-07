package main

import (
	"flag"
	"fmt"
	"net/http"
	"sitemap"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url that you want to build a sitemap for")
	flag.Parse()
	resp, err := http.Get(*urlFlag)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	ret, _ := sitemap.Parse(resp.Body)
	fmt.Println(ret)
}
