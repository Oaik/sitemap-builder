package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sitemap"
	"strings"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com/", "the url that you want to build a sitemap for")
	maxDepth := flag.Int("depth", 3, "maximum number of links deep to traverse")
	flag.Parse()
	links := bfs(*urlFlag, *maxDepth)
	fmt.Println(links)
}

func bfs(startStr string, maxDepth int) []string {
	links := get(startStr)
	var allLinks []string
	mp := make(map[string]bool)
	for _, link := range links {
		mp[link] = true
	}
	allLinks = append(allLinks, links...)
	for i := 0; i < maxDepth; i++ {
		var nwLink []string
		for _, link := range links {
			currentLinks := get(link)
			for _, currentLink := range currentLinks {
				if _, ok := mp[currentLink]; !ok {
					nwLink = append(nwLink, currentLink)
					mp[currentLink] = true
				}
			}
		}
		links = nwLink
		allLinks = append(allLinks, links...)
	}
	return allLinks
}

func get(urlString string) []string {
	resp, err := http.Get(urlString)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	base := reqUrl.Scheme + "://" + reqUrl.Host

	links := hrefs(resp.Body, base)
	return filter(links, withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := sitemap.Parse(r)
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		}
	}
	return ret
}

func filter(links []string, filterFun func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if filterFun(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(prefix string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, prefix)
	}
}
