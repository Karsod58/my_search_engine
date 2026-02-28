package documents

func Sample() []Document {
	return []Document{
		NewDoc("doc1", "Go makes concurrency easy and powerful"),
		NewDoc("doc2", "Concurrency in Go is achieved using goroutines"),
		NewDoc("doc3","Rust provides fearless concurrency guarantees" ),
		NewDoc("doc4", "Go is a programming language"),
		NewDoc("doc5", "Distributed systems require careful concurrency control"),
		NewDoc("doc6", "Go makes it easy to build simple reliable software"),
	}
}