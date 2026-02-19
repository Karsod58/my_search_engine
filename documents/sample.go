package documents

func Sample() []Document {
	return []Document{
		NewDoc("doc1", "Go is expressive concise clean and efficient"),
		NewDoc("doc2", "Concurrency is not parallelism"),
		NewDoc("doc3", "Channels orchestrate communication"),
		NewDoc("doc4", "Go makes it easy to build simple reliable software"),
	}
}