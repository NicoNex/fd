package main

import (
	"os"
	"fmt"
	"sync"
	"regexp"
)

var matches chan string
var wg sync.WaitGroup

func printerRoutine() {
	for m := range matches {
		fmt.Println(m)
	}
}

func evaluate(path string, filename string) {
	if ok, _ := regexp.MatchString(os.Args[1], filename); ok {
		matches <- path
	}
}

func walkDir(root string) {
	defer wg.Done()
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
			wg.Add(1)
			go walkDir(path)
		}
		evaluate(path, file.Name())
	}
}

func main() {
	go printerRoutine()
	wg.Add(1)
	go walkDir(".")
	wg.Wait()
	close(matches)
}

func init() {
	matches = make(chan string)
}
