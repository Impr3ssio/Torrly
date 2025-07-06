package bencode

import (
	"bytes"
	"fmt"
	"strconv"
)

type BDecoded interface{}

func DecodeBencode(input string) (BDecoded, error) {
	r := bytes.NewReader([]byte(input))
	return decode(r)
}

func decode(enc *bytes.Reader) (BDecoded, error) {
	b, err := enc.ReadByte()
	if err != nil {
		return nil, err
	}

	switch b {
	case 'i':
		return decodeInt(enc)

	default:
		return nil, fmt.Errorf("unexpected byte: %q", b)
	}
}

func decodeInt(enc *bytes.Reader) (int, error) {
	var num []byte

	for {
		letter, err := enc.ReadByte()
		if err != nil {
			return 0, err
		}

		if letter == 'e' {
			break
		}
		num = append(num, letter)
	}

	return strconv.Atoi(string(num))
}
