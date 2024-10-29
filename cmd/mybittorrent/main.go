package main

import (
	"encoding/json"
	"fmt"
	"github.com/codecrafters-io/bittorrent-starter-go/cmd/mybittorrent/decoder"
	"io"
	"os"
)

func main() {

	command := os.Args[1]

	if command == "decode" {
		decode()
	} else if command == "info" {
		err := info()
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}

func decode() {
	decoder := decoder.NewBencodeDecoder()
	bencodedValue := os.Args[2]
	decoded, _, err := decoder.DecodeBencode([]rune(bencodedValue))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(str(decoded))
}

func info() error {
	decoder := decoder.NewBencodeDecoder()
	filePath := os.Args[2]

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	println(string(data))

	if err != nil {
		return err
	}

	decoded, _, err := decoder.DecodeBencode([]rune(string(data)))
	println(str(decoded), err)

	if err != nil {
		return err
	}

	/*
		var metadata domain.Metadata
		s := str(decoded)
		err = json.NewDecoder(strings.NewReader(s)).Decode(&metadata)
		if err != nil {
			return err
		}
	*/

	var m = decoded.(map[string]interface{})
	fmt.Println(fmt.Sprintf("Tracker URL: %v", m["announce"]))
	fmt.Println(fmt.Sprintf("Length: %v", m["info"].(map[string]interface{})["length"]))

	return nil
}

func str(decoded interface{}) string {
	switch decoded.(type) {
	case string:
		if len(decoded.(string)) == 0 {
			return ""
		} else if decoded.(string)[0] == '{' {
			return decoded.(string)
		} else {
			return "\"" + decoded.(string) + "\""
		}
	default:
		jsonOutput, _ := json.Marshal(decoded)
		s := string(jsonOutput)
		return s
	}
}
