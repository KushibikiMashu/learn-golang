package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	filename := "./file.txt"
	content := []byte(
		`hello world
hello world
hello world`)
	ioutil.WriteFile(filename, content, os.ModePerm)

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}
