/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : config_test.go

* Purpose :

* Creation Date : 05-14-2017

* Last Modified : Mon 15 May 2017 06:30:46 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"log"
	"testing"
)

func TestGeoConfig(t *testing.T) {
	conf, err := LoadConfig("./conf")
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(Json(conf))
}
