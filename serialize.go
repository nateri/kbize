package netgo

import (
	"fmt"
	"os"
)


func Serialize(filename string, data interface{}) (err error) {
    f, err := os.Create(filename)
    if err != nil {
        return err
    }

    defer f.Close()
	
	fmt.Fprintf(f, "%+v", data)
	
	return nil
}
func SerializeAppend(filename string, data interface{}) (err error) {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }

    defer f.Close()
	
	fmt.Fprintf(f, "%+v", data)
	
	return nil
}


