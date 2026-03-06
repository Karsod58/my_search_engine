package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
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
		MaxPages: 3, 
	}
}

type job struct {
	url   string
	depth int
}
func (c *Crawler) Crawl(link string, depth int, handler func(string, string)) {

	if depth > c.MaxDepth {
		return
	}

	if c.PageCount >= c.MaxPages {
		return
	}

	if c.Visited[link] {
		return
	}

	c.Visited[link] = true
	c.PageCount++

	time.Sleep(300 * time.Millisecond)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", "MiniSearchBot/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Only process HTML
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
		handler(link, text)
	}

	links := extractLinks(doc, link)

	for _, l := range links {
		c.Crawl(l, depth+1, handler)
	}
}


func (c *Crawler) Start(seed string, handler func(string, string)) {

	jobs := make(chan job, 100)

	// wg waits for all workers to finish
	var wg sync.WaitGroup
	// jobWg tracks outstanding crawl jobs so we know when we're done
	var jobWg sync.WaitGroup

	workerCount := 5

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				// #region agent log
				logCrawlerEvent("crawler/crawler.go:Start-worker", "pre_process", map[string]interface{}{
					"url":   j.url,
					"depth": j.depth,
				}, "run1", "H1")
				// #endregion
				c.process(j, jobs, handler, &jobWg)
				// #region agent log
				logCrawlerEvent("crawler/crawler.go:Start-worker", "post_process", map[string]interface{}{
					"url":   j.url,
					"depth": j.depth,
				}, "run1", "H1")
				// #endregion
				jobWg.Done()
			}
		}()
	}

	// seed initial job
	jobWg.Add(1)
	// #region agent log
	logCrawlerEvent("crawler/crawler.go:Start", "seed_job", map[string]interface{}{
		"url":   seed,
		"depth": 0,
	}, "run1", "H1")
	// #endregion
	jobs <- job{url: seed, depth: 0}

	// wait until all crawl jobs (including newly spawned ones) are done
	// #region agent log
	logCrawlerEvent("crawler/crawler.go:Start", "wait_begin", map[string]interface{}{}, "run1", "H2")
	// #endregion
	jobWg.Wait()
	// #region agent log
	logCrawlerEvent("crawler/crawler.go:Start", "wait_end", map[string]interface{}{}, "run1", "H2")
	// #endregion

	// now it's safe to stop workers
	close(jobs)
	wg.Wait()
}

func (c *Crawler) process(j job, jobs chan job, handler func(string, string), jobWg *sync.WaitGroup) {

	// #region agent log
	logCrawlerEvent("crawler/crawler.go:process", "enter", map[string]interface{}{
		"url":        j.url,
		"depth":      j.depth,
		"pageCount":  c.PageCount,
		"maxPages":   c.MaxPages,
		"maxDepth":   c.MaxDepth,
	}, "run1", "H3")
	// #endregion

	if j.depth > c.MaxDepth {
		// #region agent log
		logCrawlerEvent("crawler/crawler.go:process", "skip_max_depth", map[string]interface{}{
			"url":   j.url,
			"depth": j.depth,
		}, "run1", "H3")
		// #endregion
		return
	}

	c.mu.Lock()

	if c.PageCount >= c.MaxPages {
		// #region agent log
		logCrawlerEvent("crawler/crawler.go:process", "skip_max_pages", map[string]interface{}{
			"url":       j.url,
			"pageCount": c.PageCount,
			"maxPages":  c.MaxPages,
		}, "run1", "H4")
		// #endregion
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

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

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
		// #region agent log
		logCrawlerEvent("crawler/crawler.go:process", "non_html", map[string]interface{}{
			"url":         j.url,
			"contentType": contentType,
		}, "run1", "H5")
		// #endregion
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
		// respect max depth here as well
		if j.depth+1 > c.MaxDepth {
			continue
		}
		// track each enqueued job so Start can know when we're finished
		jobWg.Add(1)
		// #region agent log
		logCrawlerEvent("crawler/crawler.go:process", "enqueue_child", map[string]interface{}{
			"parentUrl": j.url,
			"childUrl":  l,
			"depth":     j.depth + 1,
		}, "run1", "H6")
		// #endregion
		jobs <- job{url: l, depth: j.depth + 1}
	}
}

// #region agent log
func logCrawlerEvent(location, message string, data map[string]interface{}, runId, hypothesisId string) {
	f, err := os.OpenFile("debug-815281.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	entry := map[string]interface{}{
		"sessionId":    "815281",
		"id":           fmt.Sprintf("log_%d", time.Now().UnixNano()),
		"timestamp":    time.Now().UnixMilli(),
		"location":     location,
		"message":      message,
		"data":         data,
		"runId":        runId,
		"hypothesisId": hypothesisId,
	}

	b, err := json.Marshal(entry)
	if err != nil {
		return
	}

	_, _ = f.Write(append(b, '\n'))
}
// #endregion
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
