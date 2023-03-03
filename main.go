package main

import (
	"fmt"
	"github.com/geek377148474/redis-go-example/redis"
	"log"
	"net"
	"os"
)

const Network = "tcp"
const Address = "192.168.191.132:6379"
const Pwd = "123456"

func main() {
	args := os.Args[1:]
	if len(args) <= 0 {
		log.Fatalf("Os.Args <= 0")
	}
	handle(args)
}

func handle(args []string) {
	// 连接 Redis 服务器
	redisConn, err := redis.Conn(Network, Address)
	if err != nil {
		log.Fatalf("Conn err: %v", err)
	}
	defer redisConn.Close()

	if Pwd != "" {
		fmt.Println("==============AUTH end==============")
		doHandle(redisConn, []string{"auth", Pwd})
		fmt.Println("==============AUTH start==============")
	}
	doHandle(redisConn, args)
}

func doHandle(redisConn net.Conn, args []string) {
	// 获取请求协议
	reqCommand := redis.GetRequest(args)

	// 写入请求内容
	_, err := redisConn.Write(reqCommand)
	if err != nil {
		log.Fatalf("Conn Write err: %v", err)
	}

	// 读取回复
	command := make([]byte, 1024)
	n, err := redisConn.Read(command)
	if err != nil {
		log.Fatalf("Conn Read err: %v", err)
	}

	// 处理回复
	reply, err := redis.GetReply(command[:n])
	if err != nil {
		log.Fatalf("protocol.GetReply err: %v", err)
	}

	// 处理后的回复内容
	log.Printf("Reply: %v", reply)
	// 原始的回复内容
	//log.Printf("Command: %v", string(command[:n]))
}
