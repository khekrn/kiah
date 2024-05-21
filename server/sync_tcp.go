package server

import (
	"coding2fun.in/kiah/config"
	"coding2fun.in/kiah/core"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

func RunSyncTCPServer() {
	log.Println("starting a synchronous tcp server on", config.Host, config.Port)
	connectedClients := 0
	listener, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		connectedClients += 1
		log.Println("client connected with address: ", conn.RemoteAddr(), "concurrent clients", connectedClients)
		for {
			cmd, err := readCommand(conn)
			if err != nil {
				conn.Close()
				connectedClients -= 1
				log.Println("client disconnected", conn.RemoteAddr(), "concurrent clients", connectedClients)
				if err == io.EOF {
					break
				}
				log.Println("error: ", err.Error())
			}
			log.Println("command :", cmd)
			respond(cmd, conn)
		}
	}
}

func respond(cmd *core.RedisCommand, conn net.Conn) {
	err := core.EvalAndRespond(cmd, conn)
	if err != nil {
		respondError(err, conn)
	}
}

func respondError(err error, conn net.Conn) {
	_, _ = conn.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func readCommand(conn net.Conn) (*core.RedisCommand, error) {
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	tokens, err := core.DecodeArrayString(buf[:n])
	if err != nil {
		return nil, err
	}
	return &core.RedisCommand{
		Command: strings.ToUpper(tokens[0]),
		Args:    tokens[1:],
	}, nil
}
