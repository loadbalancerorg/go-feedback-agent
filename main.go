package main

import (
	"github.com/kardianos/service"
	"log"
)

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	// Do work here
	InitConfig()
	srv := InitServer()
	// accept connection on port

	// run loop forever (or until ctrl-c)
	for {
		conn, err := srv.server.Accept()
		if err != nil {
			log.Println(err)
		}
		go handleClient(conn)
	}
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "FeedBackService",
		DisplayName: "TCP Feedback Service",
		Description: "This is a go service to provide system stats",
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

//http://decouvric.cluster013.ovh.net/golang/thirdparty/divers/creer-un-service-golang-avec-kardianos.html