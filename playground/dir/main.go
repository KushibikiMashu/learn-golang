// ディレクトリ操作
package main

import (
	"fmt"
	"os"
)

func main() {
	dir := "./testdir"
	os.Mkdir(dir, os.ModePerm)
	os.Chmod(dir, 0777)

	recursiveMkdir(dir)
}

func recursiveMkdir(parent string) {
	// 数字で range する
	for i := range [10]int{} {
		dirname := fmt.Sprintf("%v/%d", parent, i+1)
		os.MkdirAll(dirname, os.ModePerm)
	}
}
