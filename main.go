package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
)

type Line struct {
	data       []byte
	lineNumber int
}

type Formatted struct {
	stdout, stderr string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lineNumber := 0
	formattedLines := make(chan chan Formatted, runtime.NumCPU())

	process := func(line Line, out chan<- Formatted) {
		data := map[string]interface{}{}
		err := json.Unmarshal(line.data, &data)
		if err != nil {
			out <- Formatted{stderr: errorMessage(err, "json.Unmarshal", line)}
			return
		}
		b, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			out <- Formatted{stderr: errorMessage(err, "json.MarshalIndent", line)}
			return
		}
		out <- Formatted{stdout: string(b) + "\n"}
	}

	go func() {
		for scanner.Scan() {
			lineNumber++
			b := scanner.Bytes()
			c := make([]byte, len(b))
			copy(c, b)
			form := make(chan Formatted)
			formattedLines <- form
			go process(Line{lineNumber: lineNumber, data: c}, form)
		}
		close(formattedLines)
	}()

	for f := range formattedLines {
		formatted := <-f
		if formatted.stderr != "" {
			os.Stderr.WriteString(formatted.stderr)
		} else {
			os.Stdout.WriteString(formatted.stdout)
		}
	}
}

func errorMessage(err error, method string, line Line) string {
	return fmt.Sprintf("%s error at line %d : %v ; input string: '%s'\n", method, line.lineNumber, err, string(line.data))
}
