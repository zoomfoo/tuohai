package main

import (
	"log"

	api "tuohai/im_api"
	"tuohai/internal/svc"
)

type program struct {
	ImApi *api.ImApi
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
	opts := api.NewOptions()
	p.ImApi = api.New(opts)
	p.ImApi.Main()
	return nil
}

func (p *program) Stop() error {
	p.ImApi.Close()
	log.Println("停止")
	return nil
}
