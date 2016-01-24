package magopie

import (
	"fmt"
	"io/ioutil"
)

// SaveToFile saves 'data' to the file 'filePath'
func SaveToFile(filePath string, data []byte) error {
	return ioutil.WriteFile(filePath, data, 0644)
}

// ReadFromFile returns the bytes from 'filePath'
func ReadFromFile(filePath string) []byte {
	var data []byte
	var err error
	if data, err = ioutil.ReadFile(filePath); err != nil {
		fmt.Println(err)
	}
	return data
}
