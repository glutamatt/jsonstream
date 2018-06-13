package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const inputdelimiter = '\n'

func main() {
	reader := bufio.NewReader(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	for {
		input, err := reader.ReadString(inputdelimiter)
		if checkError(err) {
			data := map[string]interface{}{}
			dec := json.NewDecoder(strings.NewReader(input))
			if checkError(dec.Decode(&data)) {
				enc.Encode(data)
			}
		}
	}
}

func checkError(e error) bool {
	if e == nil {
		return true
	}
	println("error")
	fmt.Println(e)
	return false
}
