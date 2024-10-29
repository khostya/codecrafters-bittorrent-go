package decoder

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
	case runes[0] == 'd':
		decoded, remains, err = d.decodeDict(runes)
	default:
		panic("Invalid BencodeString")
	}
	return
}

// Example:
// - i5e -> 5
// - i10e -> 10
// - i-10e -> -10
func (BencodeDecoder) decodeInteger(runes []rune) (int, []rune, error) {
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
	return int(n), runes[eIndex+1:], err
}

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func (BencodeDecoder) decodeString(bencodedString []rune) (string, []rune, error) {
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

	res := string(bencodedString[firstColonIndex+1 : min(len(bencodedString), firstColonIndex+1+length)])
	return res, bencodedString[min(len(bencodedString), firstColonIndex+1+length):], nil
}

// Example:
// - l5:helloi52ee -> ["hello", 52]
// - l5:helloi52ee -> []
// - lli4eei5ee -> [[4],5]
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

// Example:
// - d3:foo3:bar5:helloi52ee -> {"foo":"bar","hello":52}
// - d10:inner_dictd4:key16:value14:key2i42e8:list_keyl5:item15:item2i3eeee -> {"inner_dict":{"key1":"value1","key2":42,"list_key":["item1","item2",3]}}
func (d BencodeDecoder) decodeDict(list []rune) (interface{}, []rune, error) {
	list = list[1:]

	var dict = make(map[string]interface{})

	var (
		key string
		val interface{}
	)

	var isKey = true

	for i := 0; i != len(list); {
		r := list[i]

		if key == "" && val == nil && r == 'e' {
			list = list[1:]
			break
		}

		println(string(list))
		if isKey {
			decoded, l, err := d.decodeString(list)
			list = l

			if err != nil {
				return nil, nil, err
			}
			key = decoded
			isKey = false
		} else {
			decoded, l, err := d.DecodeBencode(list)
			list = l

			if err != nil {
				return nil, nil, err
			}

			val = decoded

			dict[key] = val

			isKey = true
			key = ""
			val = nil
		}
	}

	if len(dict) == 0 {
		return "{}", list, nil
	}

	return dict, list, nil
}
