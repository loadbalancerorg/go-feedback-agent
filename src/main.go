package main

import (
	"github.com/kardianos/service"
	"log"
)

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) run() {
	// Do work here
	InitConfig()
	srv := InitServer()

	for {
		conn, err := srv.server.Accept()
		if err != nil {
			log.Println(err)
		}
		go handleClient(conn)
	}
}

func (p *program) Stop(s service.Service) error {
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
