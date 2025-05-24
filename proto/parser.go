package proto

import (
	"bufio"
	"strconv"
)

func ReadArgs(r *bufio.Reader) ([]string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if len(line) < 3 {
		return nil, ERR_PROTOCAL
	}
	switch line[0] {
	default:
		return nil, ERR_PROTOCAL
	case '*':
		l, err := strconv.Atoi(line[1 : len(line)-2])
		if err != nil {
			return nil, err
		}
		var fields []string
		for ; l > 0; l-- {
			s, err := readString(r)
			if err != nil {
				return nil, err
			}
			fields = append(fields, s)
		}
		return fields, nil
	}
}

func readString(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	if len(line) < 3 {
		return "", ERR_PROTOCAL
	}

	switch line[0] {
	default:
		return "", ERR_PROTOCAL
	case '+', ':', '-':
		return line[1 : len(line)-2], nil
	case '$':
		length, err := strconv.Atoi(line[:len(line)-2])

		if err != nil {
			return "", ERR_PROTOCAL
		}

		if length < 0 {
			return "", nil
		}
		buf := make([]byte, length+2)
		pos := 0
		for pos < length+2 {
			n, err := r.Read(buf[pos:])
			if err != nil {
				return "", err
			}
			pos += n
		}
		return string(buf[:length]), nil
	}

}
