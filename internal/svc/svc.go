package svc

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var signalNotify = signal.Notify

type Service interface {
	Init() error
	Start() error
	Stop() error
}

func Run(service Service, sig ...os.Signal) error {
	if err := service.Init(); err != nil {
		return err
	}

	if err := service.Start(); err != nil {
		return err
	}

	if len(sig) == 0 {
		sig = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	signalChan := make(chan os.Signal, 1)
	signalNotify(signalChan, sig...)

	<-signalChan

	return service.Stop()
}

func HttpService(listener net.Listener, handler http.Handler) {
	ser := &http.Server{
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := ser.Serve(listener)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
