package documents
type Document struct {
    ID   string
    Text string
}
func NewDoc(id string,text string) Document{
 return  Document{
ID: id,
Text: text,
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