package core

import (
	"errors"
	"fmt"
	"net"
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

func DecodeArrayString(bytes []byte) ([]string, error) {
	value, err := Decode(bytes)
	if err != nil {
		return nil, err
	}
	castValue := value.([]any)
	tokens := make([]string, len(castValue))
	for i := range tokens {
		tokens[i] = castValue[i].(string)
	}
	return tokens, nil
}

func EvalAndRespond(cmd *RedisCommand, conn net.Conn) error {
	switch cmd.Command {
	case "PING":
		return evalPing(cmd.Args, conn)
	default:
		return evalPing(cmd.Args, conn)
	}
}

func Encode(value any, isSimpleString bool) []byte {
	switch val := value.(type) {
	case string:
		if isSimpleString {
			return []byte(fmt.Sprintf("+%s\r\n", val))
		}
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val))
	}
	return []byte{}
}

func evalPing(args []string, conn net.Conn) error {
	if len(args) >= 2 {
		return errors.New("ERR wrong number of arguments for 'ping' command")
	}
	var res []byte
	if len(args) == 0 {
		res = Encode("PONG", true)
	} else {
		res = Encode(args[0], false)
	}
	_, err := conn.Write(res)
	return err
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
