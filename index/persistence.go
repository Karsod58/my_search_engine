package inverted_index

import (
	"encoding/json"
	"os"

)

func (i *InvertedIndex) Save(path string) error {
	file, err := os.Create(path)
	if err!=nil{
		return err
	}
	defer file.Close()
	encoder:=json.NewEncoder(file)
	encoder.SetIndent(""," ")
	return encoder.Encode(i)
}
func Load(path string) (*InvertedIndex,error){
	file,err:=os.Open(path)
	if err!=nil{
		return nil,err
	}
	defer file.Close()
	var idx InvertedIndex
	decoder:=json.NewDecoder(file)
	if err:=decoder.Decode(&idx); err!=nil {
		return nil,err
	}
 return &idx,nil
}