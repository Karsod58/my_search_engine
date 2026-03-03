package crawler

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type Crawler struct {
	Visited   map[string]bool
	Domain    string
	MaxDepth  int
	MaxPages  int
	PageCount int

	mu sync.Mutex
}

func New(seed string, depth int) *Crawler {
	u, _ := url.Parse(seed)

	return &Crawler{
		Visited:  make(map[string]bool),
		Domain:   u.Host,
		MaxDepth: depth,
		MaxPages: 30,
	}
}

type job struct {
	url   string
	depth int
}

func (c *Crawler) Start(seed string, handler func(string, string)) {

	jobs := make(chan job, 100)
	var wg sync.WaitGroup

	workerCount := 5

	// Workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				c.process(j, jobs, handler)
			}
		}()
	}

	jobs <- job{url: seed, depth: 0}

	wg.Wait()
	close(jobs)
}

func (c *Crawler) process(j job, jobs chan job, handler func(string, string)) {

	if j.depth > c.MaxDepth {
		return
	}

	c.mu.Lock()

	if c.PageCount >= c.MaxPages {
		c.mu.Unlock()
		return
	}

	if c.Visited[j.url] {
		c.mu.Unlock()
		return
	}

	c.Visited[j.url] = true
	c.PageCount++
	c.mu.Unlock()

	time.Sleep(200 * time.Millisecond)

	client := &http.Client{}

	req, err := http.NewRequest("GET", j.url, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", "MiniSearchBot/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return
	}

	text := extractText(doc)

	if len(strings.TrimSpace(text)) > 0 {
		handler(j.url, text)
	}

	links := extractLinks(doc, j.url)

	for _, l := range links {
		jobs <- job{url: l, depth: j.depth + 1}
	}
}
func extractText(n *html.Node) string {

	// Skip script and style
	if n.Type == html.ElementNode &&
		(n.Data == "script" || n.Data == "style") {
		return ""
	}

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
					if err == nil &&
						u.Host == baseURL.Host &&
						(u.Scheme == "http" || u.Scheme == "https") {

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
