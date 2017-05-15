/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : picker.go

* Purpose :

* Creation Date : 05-15-2017

* Last Modified : Mon 15 May 2017 07:35:18 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"context"
	"github.com/kiyor/subnettool"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"strings"
)

var db *geoip2.Reader

func init() {
	var err error
	db, err = geoip2.Open("./GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatalln(err.Error())
	}
}

type Picker struct {
}

func (p *Picker) Pick(fqdn, dest string) func(ctx context.Context, network, addr string) (net.Conn, error) {
	dest = strings.Split(dest, ":")[0]

	if c, ok := myGeoConfig.MIP[dest]; ok {
		log.Println("found IP match", fqdn, dest, c)
		return c.Upstream[0].dial
	}
	for k, c := range myGeoConfig.MCIDR {
		if subnettool.CIDRMatch(dest, k) {
			log.Println("found CIDR match", fqdn, dest, c)
			return c.Upstream[0].dial
		}
	}
	if c, ok := myGeoConfig.MFQDN[fqdn]; ok {
		log.Println("found FQDN match", fqdn, dest, c)
		return c.Upstream[0].dial
	}
	country, _ := db.Country(net.ParseIP(dest))
	log.Println(dest, country.Country.IsoCode)
	if c, ok := myGeoConfig.MGEO[country.Country.IsoCode]; ok {
		log.Println("found GEO match", fqdn, dest, c)
		return c.Upstream[0].dial
	}

	log.Println("not found use default", fqdn, dest)
	return func(ctx context.Context, net_, addr string) (net.Conn, error) {
		return net.Dial(net_, addr)
	}
}
