package main

import (
	"encoding/json"
	"strconv"
	"unicode"
)

// Ensures gofmt doesn't remove the "os" encoding/json import (feel free to remove this!)
var _ = json.Marshal

type BencodeDecoder struct {
}

func NewBencodeDecoder() *BencodeDecoder {
	return &BencodeDecoder{}
}

func (d BencodeDecoder) DecodeBencode(bencodedString string) (interface{}, error) {
	runes := []rune(bencodedString)

	switch {
	case unicode.IsDigit(runes[0]):
		return d.decodeString(bencodedString)
	case runes[0] == 'i' && runes[len(runes)-1] == 'e':
		return d.decodeInteger(runes)
	default:
		panic("Invalid BencodeString")
	}
}

// Example:
// - i5e -> 5
// - i10e -> 10
func (BencodeDecoder) decodeInteger(runes []rune) (interface{}, error) {
	runes = runes[1 : len(runes)-1]

	return strconv.ParseInt(string(runes), 10, 64)
}

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func (BencodeDecoder) decodeString(bencodedString string) (interface{}, error) {
	var firstColonIndex int

	for i := 0; i < len(bencodedString); i++ {
		if bencodedString[i] == ':' {
			firstColonIndex = i
			break
		}
	}

	lengthStr := bencodedString[:firstColonIndex]

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", err
	}

	return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
}
