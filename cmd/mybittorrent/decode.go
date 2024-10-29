package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"unicode"
)

// Ensures gofmt doesn't remove the "os" encoding/json import (feel free to remove this!)
var _ = json.Marshal

type BencodeDecoder struct {
}

func NewBencodeDecoder() *BencodeDecoder {
	return &BencodeDecoder{}
}

func (d BencodeDecoder) DecodeBencode(bencodedString []rune) (decoded interface{}, remains []rune, err error) {
	runes := bencodedString

	switch {
	case unicode.IsDigit(runes[0]):
		decoded, remains, err = d.decodeString(bencodedString)
	case runes[0] == 'i':
		decoded, remains, err = d.decodeInteger(runes)
	case runes[0] == 'l':
		decoded, remains, err = d.decodeList(runes)
	default:
		panic("Invalid BencodeString")
	}
	return
}

// Example:
// - i5e -> 5
// - i10e -> 10
// - i-10e -> -10
func (BencodeDecoder) decodeInteger(runes []rune) (interface{}, []rune, error) {
	runes = runes[1:]

	var (
		buf    strings.Builder
		eIndex int
	)

	for i, r := range runes {
		if unicode.IsDigit(r) || r == '-' {
			buf.WriteRune(r)
		} else {
			eIndex = i
			break
		}
	}

	n, err := strconv.ParseInt(buf.String(), 10, 64)
	return n, runes[eIndex+1:], err
}

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func (BencodeDecoder) decodeString(bencodedString []rune) (interface{}, []rune, error) {
	var firstColonIndex int

	for i := 0; i < len(bencodedString); i++ {
		if bencodedString[i] == ':' {
			firstColonIndex = i
			break
		}
	}

	lengthStr := bencodedString[:firstColonIndex]

	length, err := strconv.Atoi(string(lengthStr))
	if err != nil {
		return "", nil, err
	}

	res := string(bencodedString[firstColonIndex+1 : firstColonIndex+1+length])
	return res, bencodedString[firstColonIndex+1+length:], nil
}

// Example:
// - l5:helloi52ee -> ["hello", 52]
func (d BencodeDecoder) decodeList(bencodedList []rune) (interface{}, []rune, error) {
	var list []interface{} = make([]interface{}, 0)

	bencodedList = bencodedList[1:]

	for len(bencodedList) != 0 {
		if bencodedList[0] == 'e' {
			return list, bencodedList[1:], nil
		}

		res, remains, err := d.DecodeBencode(bencodedList)
		if err != nil {
			return nil, nil, err
		}

		list = append(list, res)

		bencodedList = remains
	}

	return list, nil, nil
}
