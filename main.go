package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
)

type Line struct {
	data       []byte
	lineNumber int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	lineNumber := 0

	lines := make(chan Line)
	formatted := make(chan []byte)

	wg := sync.WaitGroup{}
	workers := runtime.NumCPU()
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			for l := range lines {
				data := map[string]interface{}{}
				if checkError(json.Unmarshal(l.data, &data), l.data, l.lineNumber) {
					if b, err := json.MarshalIndent(data, "", " "); err == nil {
						formatted <- b
					}
				}
			}
			wg.Done()
		}()
	}

	go func() {
		for scanner.Scan() {
			lineNumber++
			b := scanner.Bytes()
			c := make([]byte, len(b))
			copy(c, b)
			lines <- Line{lineNumber: lineNumber, data: c}
		}
		close(lines)
		wg.Wait()
		close(formatted)
	}()

	for f := range formatted {
		os.Stdout.Write(f)
		os.Stdout.WriteString("\n")
	}
}

func checkError(e error, js []byte, lineNumber int) bool {
	if e == nil {
		return true
	}

	fmt.Fprintf(os.Stderr, "Json error at line %d : %v\nfrom string:\n%s\n", lineNumber, e, string(js))
	return false
}
