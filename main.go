package main

import (
	"coding2fun.in/kiah/config"
	"coding2fun.in/kiah/server"
	"flag"
	"log"
)

func configureFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for kiah server")
	flag.IntVar(&config.Port, "port", 7379, "port got the kiah server")
	flag.Parse()
}
func main() {
	configureFlags()
	log.Println("spinning new server 🌅...")
	server.RunSyncTCPServer()
}
