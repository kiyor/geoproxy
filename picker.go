/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : picker.go

* Purpose :

* Creation Date : 05-15-2017

* Last Modified : Tue 16 May 2017 08:03:21 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"context"
	"github.com/kiyor/go-socks5"
	"github.com/kiyor/subnettool"
	"github.com/oschwald/geoip2-golang"
	"github.com/ryanuber/go-glob"
	"github.com/wsxiaoys/terminal/color"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"strings"
	"time"
)

var db *geoip2.Reader

func init() {
	var err error
	db, err = geoip2.Open("./GeoLite2-City.mmdb")
	if err != nil {
		log.Fatalln(err.Error())
	}
}

type Picker struct {
}

func proxyDialer(p string, auth *proxy.Auth) func(ctx context.Context, network, addr string) (net.Conn, error) {
	if len(p) > 0 {
		dialer, err := proxy.SOCKS5("tcp", p,
			auth,
			&net.Dialer{
				KeepAlive: 30 * time.Second,
			},
		)
		if err != nil {
			log.Println(p, auth, err.Error())
		}
		return func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
	}
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return net.Dial(network, addr)
	}
}

func (p *Picker) Pick(r *socks5.Request) func(ctx context.Context, network, addr string) (net.Conn, error) {
	fqdn := r.DestAddr.FQDN
	dest := r.RealDestAddr().IP.String()

	// 	dialGen := func(proxy string) func(ctx context.Context, network, addr string) (net.Conn, error) {
	// 		return net.Dial(network, addr)
	// 	}
	found := color.Sprint("@{g}HIT@{|}")
	notFound := color.Sprint("@{r}MISS@{|}")
	from := r.RemoteAddr.IP.String()

	myGeoConfig.RLock()
	if c, ok := myGeoConfig.Cache[CacheKey{fqdn, dest}]; ok {
		myGeoConfig.RUnlock()
		u := c.Upstream[0]
		log.Println(found, "Cache match", from, fqdn, dest, u)
		return proxyDialer(u.Addr, u.auth)
	}
	myGeoConfig.RUnlock()

	if c, ok := myGeoConfig.MIP[dest]; ok {
		myGeoConfig.Lock()
		myGeoConfig.Cache[CacheKey{fqdn, dest}] = c
		myGeoConfig.Unlock()
		u := c.Upstream[0]

		log.Println(found, "IP match", from, fqdn, dest, u)
		return proxyDialer(u.Addr, u.auth)
	}
	for k, c := range myGeoConfig.MCIDR {
		if subnettool.CIDRMatch(dest, k) {
			myGeoConfig.Lock()
			myGeoConfig.Cache[CacheKey{fqdn, dest}] = c
			myGeoConfig.Unlock()
			u := c.Upstream[0]
			log.Println(found, "CIDR match", from, fqdn, dest, u)
			return proxyDialer(u.Addr, u.auth)
		}
	}
	if c, ok := myGeoConfig.MFQDN[fqdn]; ok {
		myGeoConfig.Lock()
		myGeoConfig.Cache[CacheKey{fqdn, dest}] = c
		myGeoConfig.Unlock()
		u := c.Upstream[0]
		log.Println(found, "FQDN match", from, fqdn, dest, u)
		return proxyDialer(u.Addr, u.auth)
	}
	for k, c := range myGeoConfig.MREFQDN {
		if glob.Glob(k, fqdn) {
			myGeoConfig.Lock()
			myGeoConfig.Cache[CacheKey{fqdn, dest}] = c
			myGeoConfig.Unlock()
			u := c.Upstream[0]
			log.Println(found, "FQDN match", from, fqdn, dest, u)
			return proxyDialer(u.Addr, u.auth)
		}
	}
	// 	country, _ := db.Country(net.ParseIP(dest))
	city, _ := db.City(net.ParseIP(dest))
	codes := []string{city.Country.IsoCode}
	for _, v := range city.Subdivisions {
		codes = append(codes, v.IsoCode)
	}
	code := color.Sprintf("@{c}%s@{|}", strings.Join(codes, "-"))
	if c, ok := myGeoConfig.MGEO[city.Country.IsoCode]; ok {
		myGeoConfig.Lock()
		myGeoConfig.Cache[CacheKey{fqdn, dest}] = c
		myGeoConfig.Unlock()
		u := c.Upstream[0]
		log.Println(found, "GEO match", from, fqdn, dest, code, u)
		return proxyDialer(u.Addr, u.auth)
	}

	log.Println(notFound, "GEO", code, "use default", from, fqdn, dest)
	myGeoConfig.Lock()
	myGeoConfig.Cache[CacheKey{fqdn, dest}] = myGeoConfig.Default
	myGeoConfig.Unlock()
	return func(ctx context.Context, net_, addr string) (net.Conn, error) {
		return net.Dial(net_, addr)
	}
}
