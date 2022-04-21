package main

import (
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
)

type Server struct {
	server net.Listener
}

func InitServer() *Server {
	srv := &Server{}
	// If Port is not specified, use 3333 by default
	if _, err := strconv.Atoi(GlobalConfig.Port.ToString()); err != nil {
		GlobalConfig.Port.Value = "3333"
	}
	listner, err := net.Listen("tcp", ":"+GlobalConfig.Port.ToString())
	if err != nil {
		log.Fatal(err)
	}
	srv.server = listner
	return srv
}
