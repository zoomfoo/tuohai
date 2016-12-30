package pushnotify

import (
	"crypto/tls"
)

import (
// "strconv"
)

type Options struct {
	TLSCert     string
	TLSKey      string
	P12File     string
	P12Password string
	Production  bool
	RedisHost   string
	Subscribers string

	TLSMinVersion uint16
	MaxVersion    uint16
}

func NewOptions() *Options {
	return &Options{
		Production:  true,
		P12Password: "123456",
		RedisHost:   "127.0.0.1:6379",
		Subscribers: "pushnotify",

		TLSMinVersion: tls.VersionTLS10,
	}
}
