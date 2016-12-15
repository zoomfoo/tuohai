package main

import (
	"log"

	api "tuohai/file_api"
	"tuohai/file_api/options"
	"tuohai/internal/svc"
)

type program struct {
}

func main() {
	if err := svc.Run(&program{}); err != nil {
		log.Print("ERROR: ")
		log.Fatal(err)
	}
}

func (p *program) Init() error {
	log.Println("启动")
	return nil
}

func (p *program) Start() error {
	opts := options.NewOptions()
	api.New(opts).Main()
	return nil
}

func (p *program) Stop() error {
	log.Println("停止")
	return nil
}
