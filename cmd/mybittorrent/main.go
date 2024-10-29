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

		switch decoded.(type) {
		case string:
			if len(decoded.(string)) == 0 {
				fmt.Println()
			} else if decoded.(string)[0] == '{' {
				fmt.Println(decoded.(string))
			} else {
				fmt.Println("\"" + decoded.(string) + "\"")
			}
		default:
			jsonOutput, _ := json.Marshal(decoded)
			fmt.Println(string(jsonOutput))
		}

	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
