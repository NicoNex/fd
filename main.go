package main

import (
	"os"
	"fmt"
	"sync"
	"regexp"
	// "path/filepath"
)

// var matches chan string
var wg sync.WaitGroup

func evaluate(filename string) {
	if ok, _ := regexp.MatchString(os.Args[1], filename); ok {
		println(filename)
		matches <- filename
	}
}

func walkDir(root string) {
	f, err := os.Open(root)
	if err != nil {
		return
	}

	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return
	}

	for _, file := range fileInfo {
		path := fmt.Sprintf("%s/%s", root, file.Name())
		if file.IsDir() {
			go walkDir(path)
		} else {
			// matches <- path
			go fmt.Println(path)
		}
	}
}

func main() {
	walkDir(".")
}

// func init() {
// 	matches = make(chan string)
// }
