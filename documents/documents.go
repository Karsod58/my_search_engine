package documents
type Document struct {
    ID    string
    Text  string
    Title string
    URL   string
}
func NewDoc(id string, text string) Document {
	return Document{
		ID:   id,
		Text: text,
	}
}

func NewDocWithMeta(id, text, title, url string) Document {
	return Document{
		ID:    id,
		Text:  text,
		Title: title,
		URL:   url,
	}
}
func GetByID(docs []Document, id string) *Document{
	for _, d := range docs {
		if d.ID == id {
			return &d
		}
	}
	return nil
}