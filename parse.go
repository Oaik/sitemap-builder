package sitemap

import (
	"io"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	// Text string
	// Text not needed in this
}

func dfs(n *html.Node) []*html.Node {
	var ret []*html.Node
	if n.Type == html.ElementNode && n.Data == "a" {
		ret = append(ret, n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, dfs(c)...)
	}
	return ret
}

func buildLink(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	return ret
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	linkNodes := dfs(doc)
	var allLinks []Link
	for _, n := range linkNodes {
		allLinks = append(allLinks, buildLink(n))
	}
	return allLinks, nil
}
