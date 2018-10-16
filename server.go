package main

import (
	"log"
	"net"
)

type Server struct {
	server net.Listener
}

func InitServer() *Server {
	srv := &Server{}
	listner, err := net.Listen("tcp", ":" + GlobalConfig.Port.ToString())
	if err != nil {
		log.Fatal(err)
	}
	srv.server = listner
	return srv
}
