package memcache

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type MCReq struct {
	Command  string
	Key      string
	Flags    string
	ExpireAt int64
	Data     []byte
	Noreply  bool
	Offset   int64
}

func Read(r *bufio.Reader) (req *MCReq, err error) {
	//set key 0 30 5
	//value
	lineBytes, _, err := r.ReadLine()
	if err != nil {
		return nil, err
	}

	lineStr := string(lineBytes)
	lineArray := strings.Fields(lineStr) // split
	fmt.Println(lineArray)
	if len(lineArray) < 1 {
		return nil, NewMCError(CLIENT_ERROR, fmt.Sprintf("too few params to command %s", lineArray[0]))
	}

	command := lineArray[0]
	switch command {
	case "set", "add", "replace":
		if len(lineArray) < 5 {
			return nil, NewMCError(CLIENT_ERROR, fmt.Sprintf("too few params to command %s", lineArray[0]))
		}
		req := &MCReq{}
		req.Command = command
		req.Key = lineArray[1]
		req.Flags = lineArray[2]
		req.ExpireAt, err = strconv.ParseInt(lineArray[3], 10, 64)
		if err != nil {
			return nil, err
		}
		dataLen, err := strconv.Atoi(lineArray[4])
		if err != nil {
			return nil, err
		}
		if len(lineArray) > 5 && lineArray[5] == "noreply" {
			req.Noreply = true
		}
		req.Data = make([]byte, dataLen)
		n, err := r.Read(req.Data)
		if err != nil {
			return nil, err
		}
		if n != dataLen {
			return nil, NewMCError(CLIENT_ERROR, fmt.Sprintf("Read only %d bytes of %d bytes of expected data", n, dataLen))
		}

		c, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if c != '\r' {
			return nil, NewMCError(CLIENT_ERROR, "expected \\r")
		}
		c, err = r.ReadByte()
		if err != nil {
			return nil, err
		}
		if c != '\n' {
			return nil, NewMCError(CLIENT_ERROR, "expected \\n")
		}

		// end := make([]byte, 2)
		// _, err = r.Read(end)
		// if err != nil {
		// 	return nil, err
		// }
		// if bytes.Equal(end, []byte("\r\n")) {
		// 	return nil, NewMCError(CLIENT_ERROR, "expected \\r\\n")
		// }

		return req, nil
		break
	case "get":
		//get key
		if len(lineArray) < 2 {
			return nil, NewMCError(CLIENT_ERROR, fmt.Sprintf("too few params to command %s", lineArray[0]))
		}
		req := &MCReq{}
		req.Command = command
		req.Key = lineArray[1]
		return req, nil
		break
	case "delete":
		//delete key
		if len(lineArray) < 2 {
			return nil, NewMCError(CLIENT_ERROR, fmt.Sprintf("too few params to command %s", lineArray[0]))
		}
		req := &MCReq{}
		req.Command = command
		req.Key = lineArray[1]
		return req, nil
		break
	case "version":
		// version
		req := &MCReq{}
		req.Command = command
		return req, nil
		break
	case "join":
		req := &MCReq{}
		req.Command = command
		req.Key = lineArray[1]
		return req, nil
		break
	case "quit":
		// quit
		req := &MCReq{}
		req.Command = command
		return req, nil
		break
	case "incr", "decr":
		// incr a 2
		if len(lineArray) < 3 {
			return nil, NewMCError(CLIENT_ERROR, fmt.Sprintf("too few params to command %s", lineArray[0]))
		}
		req := &MCReq{}
		req.Command = command
		req.Key = lineArray[1]
		req.Offset, err = strconv.ParseInt(lineArray[2], 10, 64)
		if err != nil {
			return nil, err
		}
		return req, nil
		break
	default:

	}
	return nil, NewMCError(ERROR, "")
}
