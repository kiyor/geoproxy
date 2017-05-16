/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 05-15-2017

* Last Modified : Mon May 15 15:31:06 2017

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"flag"
	"github.com/kiyor/go-socks5"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	fListen = flag.String("l", "127.0.0.1:1080", "listen interface")
	fAuth   = flag.String("auth", "", "auth file.txt file.json or 'user:pass'")
	fConf   = flag.String("c", "./conf", "conf dir")
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for range c {
			log.Println("reload")
			Reload(*fConf)
		}
	}()

	conf := &socks5.Config{
		Picker: new(Picker),
	}
	if *fAuth != "" {
		cred := parseSocks5Auth(*fAuth)
		cator := socks5.UserPassAuthenticator{Credentials: cred}
		conf.AuthMethods = []socks5.Authenticator{cator}
	}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", *fListen); err != nil {
		panic(err)
	}
}
