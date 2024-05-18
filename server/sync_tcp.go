package server

import (
	"coding2fun.in/kiah/config"
	"io"
	"log"
	"net"
	"strconv"
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
			if err = respond(cmd, conn); err != nil {
				log.Println("error while writing to connection : ", err.Error())
				break
			}
		}
	}
}

func respond(cmd string, conn net.Conn) error {
	if _, err := conn.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
}

func readCommand(conn net.Conn) (string, error) {
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}
