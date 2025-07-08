package bencode

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
)

type BValue interface{}

var (
	ErrBlankInput       = errors.New("input cannot be an empty string")
	ErrUnexpectedEOF    = errors.New("unexpected end of input")
	ErrInvalidFormat    = errors.New("invalid bencode format")
	ErrInvalidInteger   = errors.New("invalid integer in bencode")
	ErrInvalidStringLen = errors.New("invalid string length")
	ErrNonStringKey     = errors.New("dictionary keys must be strings")
)

func DecodeBencode(input string) (BValue, error) {
	if len(input) <= 0 {
		return nil, ErrBlankInput
	}

	r := bytes.NewReader([]byte(input))
	if r.Len() == 0 {
		return nil, ErrUnexpectedEOF
	}

	var res []BValue

	for r.Len() > 0 {
		val, err := decode(r)
		if err != nil {
			return nil, err
		}
		res = append(res, val)
	}

	if len(res) == 1 {
		return res[0], nil
	}
	return res, nil
}

func decode(r *bytes.Reader) (BValue, error) {
	if r.Len() == 0 {
		return nil, ErrUnexpectedEOF
	}

	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	switch {
	case b == 'l':
		return extractList(r)
	case b == 'i':
		return extractInteger(r)
	case b == 'd':
		return extractDictionary(r)
	case b >= '0' && b <= '9':
		if err := r.UnreadByte(); err != nil {
			return nil, err
		}
		return extractString(r)
	default:
		return nil, fmt.Errorf("%w: unexpected byte %q", ErrInvalidFormat, b)
	}
}

// Format-> <length>:<data>
func extractString(r *bytes.Reader) (BValue, error) {
	var len []byte

	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		if b == ':' {
			break
		}

		if b < '0' || b > '9' {
			return nil, ErrInvalidFormat
		}

		len = append(len, b)
	}

	length, err := strconv.Atoi(string(len))
	if err != nil || length < 0 {
		return nil, ErrInvalidStringLen
	}

	strArr := make([]byte, length)
	if _, err = io.ReadFull(r, strArr); err != nil {
		return nil, err
	}
	return string(strArr), nil
}

// Format-> i<integer>e, e.g i69e -> 69
func extractInteger(r *bytes.Reader) (int, error) {
	var num []byte

	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		if b == 'e' {
			break
		}

		if b < '0' || b > '9' {
			return 0, ErrInvalidInteger
		}

		num = append(num, b)
	}

	if len(num) == 0 || (len(num) == 1 && num[0] == '-') {
		return 0, ErrInvalidInteger
	}

	parsed, err := strconv.Atoi(string(num))
	if err != nil {
		return 0, fmt.Errorf("%w: could not parse integer", ErrInvalidFormat)
	}
	return parsed, nil
}

// Format-> l<bencoded values>e
func extractList(r *bytes.Reader) (BValue, error) {
	var list []BValue

	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		if b == 'e' {
			break
		}

		r.UnreadByte()

		val, err := decode(r)
		if err != nil {
			return nil, err
		}
		list = append(list, val)
	}

	return list, nil
}

// Format-> d<bencoded string><bencoded element>e
func extractDictionary(r *bytes.Reader) (map[string]BValue, error) {
	var dict = make(map[string]BValue)

	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		if b == 'e' {
			break
		}

		if err := r.UnreadByte(); err != nil {
			return nil, err
		}

		keyRaw, err := extractString(r)
		if err != nil {
			return nil, err
		}

		keyStr, ok := keyRaw.(string)
		if !ok {
			return nil, ErrNonStringKey
		}

		val, err := decode(r)
		if err != nil {
			return nil, err
		}

		dict[keyStr] = val
	}
	return dict, nil
}

func Encode(v BValue) ([]byte, error) {
	var buf bytes.Buffer
	err := encode(&buf, v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encode(w *bytes.Buffer, v BValue) error {
	switch val := v.(type) {
	case string:
		_, err := fmt.Fprintf(w, "%d:%s", len(val), val)
		return err
	case int:
		_, err := fmt.Fprintf(w, "i%de", val)
		return err
	case []BValue:
		w.WriteByte('l')
		for _, item := range val {
			if err := encode(w, item); err != nil {
				return err
			}
		}
		w.WriteByte('e')
		return nil
	case map[string]BValue:
		w.WriteByte('d')

		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			if err := encode(w, k); err != nil {
				return err
			}
			if err := encode(w, val[k]); err != nil {
				return err
			}
		}

		w.WriteByte('e')
		return nil
	default:
		return fmt.Errorf("unsupported bencode type: %T", v)
	}
}
