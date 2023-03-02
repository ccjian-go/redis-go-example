package redis

import "net"

// net.Dail封装redis对服务器的连接
func Conn(network, address string) (net.Conn, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
