package main

import (
	"encoding/json"
	"fmt"
	"github.com/codecrafters-io/bittorrent-starter-go/cmd/mybittorrent/decoder"
	"github.com/codecrafters-io/bittorrent-starter-go/cmd/mybittorrent/domain"
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
	decoded, _, err := decoder.DecodeBencode(bencodedValue)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(str(decoded))
}

func info() error {
	filePath := os.Args[2]

	metadata, err := domain.NewTorrent(filePath)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Tracker URL: %v", metadata.Announce))
	fmt.Println(fmt.Sprintf("Length: %v", metadata.Info.Length))

	infoHash, err := metadata.InfoHash()
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Info Hash: %x", infoHash))
	fmt.Println(fmt.Sprintf("Piece Length: %v", metadata.Info.PieceLength))
	fmt.Println("Piece Hashes:")

	println(len(metadata.Info.Pieces))
	for i := 0; i < len(metadata.Info.Pieces); i += 20 {
		fmt.Printf("%x\n", metadata.Info.Pieces[i:i+20])
	}
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
