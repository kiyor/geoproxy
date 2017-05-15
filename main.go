/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 05-15-2017

* Last Modified : Mon 15 May 2017 05:22:21 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"github.com/kiyor/go-socks5"
)

func main() {
	conf := &socks5.Config{
		Picker: new(Picker),
	}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", ":8899"); err != nil {
		panic(err)
	}
}
