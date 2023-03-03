package redis

import (
	"github.com/geek377148474/redis-go-example/util"
	"strconv"
	"strings"
)

/**
状态回复（status reply）的第一个字节是 "+"。
错误回复（error reply）的第一个字节是 "-"。
整数回复（integer reply）的第一个字节是 ":"。
批量回复（bulk reply）的第一个字节是 "$"。
多条批量回复（multi bulk reply）的第一个字节是 "*"。
*/

var StatusReply = []byte("+")[0]
var ErrorReply = []byte("-")[0]
var IntegerReply = []byte(":")[0]
var BulkReply = []byte("$")[0]
var MultiBulkReply = []byte("*")[0]

var OkReply = "ok"
var PongReply = "pong"

func GetRequest(args []string) []byte {
	req := []string{
		"*" + strconv.Itoa(len(args)),
	}

	for _, arg := range args {
		req = append(req, "$"+strconv.Itoa(len(arg)))
		req = append(req, arg)
	}

	str := strings.Join(req, "\r\n")

	return []byte(str + "\r\n")
}

func GetReply(reply []byte) (interface{}, error) {
	replyType := reply[0]
	switch replyType {
	case StatusReply:
		return doStatusReply(reply[1:])
	case ErrorReply:
		return doErrorReply(reply[1:])
	case IntegerReply:
		//return doIntegerReply(reply[1:])
		fallthrough
	case BulkReply:
		//return doBulkReply(reply[1:])
		fallthrough
	case MultiBulkReply:
		return doMultiBulkReply(reply[1:])
	default:
		return nil, nil
	}
}

func getFlagPos(search []rune, data []byte) int {
	return strings.Index(string(data), string(search))
}

func doMultiBulkReply(reply []byte) (string, error) {
	var result []string
	i := 0
	for {
		i++
		if i == 10 {
			break
		}
		pos := getFlagPos([]rune("\r"), reply)
		if pos < 0 {
			break
		}
		str := strings.Replace(string(reply[:pos]), "\n", "", -1)
		util.P(str)
		result = append(result, str)
		reply = reply[pos+1:]
	}

	return strings.Join(result, " "), nil
}

func doStatusReply(reply []byte) (string, error) {
	if len(reply) == 3 && reply[1] == 'O' && reply[2] == 'K' {
		return OkReply, nil
	}

	if len(reply) == 5 && reply[1] == 'P' && reply[2] == 'O' && reply[3] == 'N' && reply[4] == 'G' {
		return PongReply, nil
	}

	return string(reply), nil
}

func doErrorReply(reply []byte) (string, error) {
	return string(reply), nil
}

//func doIntegerReply(reply []byte) (int, error) {
//	pos := getFlagPos('\r', reply)
//	result, err := strconv.Atoi(string(reply[:pos]))
//	if err != nil {
//		return 0, err
//	}
//
//	return result, nil
//}
