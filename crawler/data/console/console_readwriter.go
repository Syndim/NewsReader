package console

import (
	"fmt"
	"newsreader/crawler/data"
)

type ConsoleReadWriter struct {
}

func (j *ConsoleReadWriter) Read(type_ data.DataOperationType, args ...interface{}) (interface{}, error) {
	return "", nil
}

func (j *ConsoleReadWriter) Write(type_ data.DataOperationType, args ...interface{}) error {
	fmt.Print(type_)
	fmt.Print(" ")
	fmt.Println(args)
	return nil
}
