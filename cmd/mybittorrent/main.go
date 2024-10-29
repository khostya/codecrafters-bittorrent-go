package main

import (
	"encoding/json"
	"fmt"
	"os"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

func main() {
	decoder := NewBencodeDecoder()

	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]
		decoded, _, err := decoder.DecodeBencode([]rune(bencodedValue))
		if err != nil {
			fmt.Println(err)
			return
		}
		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
