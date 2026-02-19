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