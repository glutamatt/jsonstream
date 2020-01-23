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
		js := scanner.Text()
		dec := json.NewDecoder(strings.NewReader(js))
		if checkError(dec.Decode(&data), js) {
			enc.Encode(data)
		}
	}
}

func checkError(e error, js string) bool {
	if e == nil {
		return true
	}

	fmt.Fprintf(os.Stderr, "Json error: %v\nfrom string:\n%s\n", e, js)
	return false
}
