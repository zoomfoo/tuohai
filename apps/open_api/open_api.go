package main

import (
	"log"

	"tuohai/internal/svc"
	api "tuohai/open_api"
)

type program struct {
	OpenApi *api.OpenApi
}

func main() {
	if err := svc.Run(&program{}); err != nil {
		log.Print("ERROR：")
		log.Fatal(err)
	}
}

func (p *program) Init() error {
	log.Println("启动")
	return nil
}

func (p *program) Start() error {
	opts := api.NewOptions()
	open := api.NewOpenApi(opts)
	open.Main()
	return nil
}

func (p *program) Stop() error {
	log.Println("停止")
	return nil
}
