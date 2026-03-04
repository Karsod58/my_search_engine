package documents

func Sample() []Document {
	return []Document{
		NewDocWithMeta("doc1", "Go makes concurrency easy and powerful", "Go Concurrency", "https://golang.org/concurrency"),
		NewDocWithMeta("doc2", "Concurrency in Go is achieved using goroutines", "Goroutines Guide", "https://golang.org/goroutines"),
		NewDocWithMeta("doc3", "Rust provides fearless concurrency guarantees", "Rust Concurrency", "https://rust-lang.org/concurrency"),
		NewDocWithMeta("doc4", "Go is a programming language", "Go Language", "https://golang.org"),
		NewDocWithMeta("doc5", "Distributed systems require careful concurrency control", "Distributed Systems", "https://example.com/distributed"),
		NewDocWithMeta("doc6", "Go makes it easy to build simple reliable software", "Go Simplicity", "https://golang.org/simple"),
	}
}