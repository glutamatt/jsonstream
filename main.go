package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	for scanner.Scan() {
		data := map[string]interface{}{}
		dec := json.NewDecoder(strings.NewReader(scanner.Text()))
		if checkError(dec.Decode(&data)) {
			enc.Encode(data)
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
