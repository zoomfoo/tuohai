package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"tuohai/internal/svc"
	"tuohai/internal/util"
	"tuohai/pushnotify"
)

func proxyFlagSet(opts *pushnotify.Options) *flag.FlagSet {
	flagSet := flag.NewFlagSet("proxy", flag.ExitOnError)
	// basic options
	flagSet.Bool("version", false, "Version information")
	flagSet.StringVar(&opts.TLSCert, "tls-cert", opts.TLSCert, "path to certificate file")
	flagSet.StringVar(&opts.TLSKey, "tls-key", opts.TLSKey, "path to key file")
	flagSet.StringVar(&opts.P12File, "p12", opts.P12File, "path to p12 file")
	flagSet.StringVar(&opts.P12Password, "p12-pwd", opts.P12Password, "path to p12-password")
	flagSet.StringVar(&opts.Subscribers, "subscribe", opts.Subscribers, "The theme of the subscription")
	flagSet.BoolVar(&opts.Production, "production", opts.Production, "The test environment or production environment(The default production environment)")
	return flagSet
}

type program struct {
	waitGroup util.WaitGroupWrapper
}

func main() {
	prg := &program{}
	if err := svc.Run(prg); err != nil {
		log.Printf("ERROR: %v", err)
	}
}

func (p *program) Init() error {
	return nil
}

func (p *program) Start() error {
	opts := pushnotify.NewOptions()
	flagSet := proxyFlagSet(opts)
	flagSet.Parse(os.Args[1:])

	if flagSet.Lookup("version").Value.(flag.Getter).Get().(bool) {
		fmt.Println(pushnotify.String("pushnotify"))
		os.Exit(0)
	}

	n, err := pushnotify.NewPNotify()
	if err != nil {
		return err
	}

	er := n.Main(opts)
	if er != nil {
		return err
	}

	return nil
}

func (p *program) Stop() error {
	return nil
}
