package crawler

import (
	"io"
     "golang.org/x/net/html"
	"net/http"
	"net/url"
	"strings"
)

type Crawler struct {
	Visited  map[string]bool
	Domain   string
	MaxDepth int
}

func New(seed string, depth int) *Crawler {
	u, _ := url.Parse(seed)
	return &Crawler{
		Visited: make(map[string]bool),
		Domain: u.Host,
		MaxDepth: depth,
	}
}
func (c *Crawler) Crawl(link string, depth int, handler func(string, string)) {

	if depth > c.MaxDepth {
		return
	}

	if c.Visited[link] {
		return
	}

	c.Visited[link] = true

	resp, err := http.Get(link)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return
	}

	text := extractText(doc)

	handler(link, text)

	links := extractLinks(doc, link)

	for _, l := range links {
		c.Crawl(l, depth+1, handler)
	}
}
func extractText(n *html.Node) string {

	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}

	var text string

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += " " + extractText(c)
	}

	return text
}
func extractLinks(n *html.Node, base string) []string {

	var links []string

	baseURL, _ := url.Parse(base)

	var f func(*html.Node)
	f = func(n *html.Node) {

		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					u, err := baseURL.Parse(attr.Val)
					if err == nil && u.Host == baseURL.Host {
						links = append(links, u.String())
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(n)
	return links
}