package core

import (
	"errors"
	"strconv"
)

// Decode the byte to
func Decode(data []byte) (any, error) {
	if len(data) == 0 {
		return nil, errors.New("no data found")
	}
	value, _, err := DecodeOne(data)
	return value, err
}

func DecodeOne(data []byte) (any, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data found")
	}
	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInt64(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	}
	return nil, 0, nil
}

func readArray(data []byte) (any, int, error) {
	pos := 1
	length, delta, err := readLength(data[pos:])
	if err != nil {
		return nil, 0, err
	}
	pos += delta
	elements := make([]any, length)
	for i := range elements {
		result, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elements[i] = result
		pos += delta
	}
	return elements, pos, nil
}

func readBulkString(data []byte) (any, int, error) {
	pos := 1
	length, delta, err := readLength(data[pos:])
	if err != nil {
		return nil, 0, err
	}
	pos += delta
	return string(data[pos:(pos + length)]), pos + length + 2, nil
}

func readInt64(data []byte) (any, int, error) {
	pos := 1
	var value int64 = 0
	for data[pos] != '\r' {
		pos += 1
	}
	value, err := strconv.ParseInt(string(data[1:pos]), 10, 64)
	if err != nil {
		return nil, 0, err
	}
	return value, pos + 2, nil
}

func readError(data []byte) (any, int, error) {
	return readSimpleString(data)
}

func readSimpleString(data []byte) (any, int, error) {
	pos := 1
	for data[pos] != '\r' {
		pos += 1
	}
	return string(data[1:pos]), pos + 2, nil
}

func readLength(data []byte) (int, int, error) {
	pos := 1
	for data[pos] != '\r' {
		pos += 1
	}
	value, err := strconv.Atoi(string(data[:pos]))
	if err != nil {
		return 0, 0, nil
	}
	return value, pos + 2, nil
}
