/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : picker.go

* Purpose :

* Creation Date : 05-15-2017

* Last Modified : Mon May 15 17:01:22 2017

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
	"time"
	// 	"strings"
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

func (p *Picker) Pick(r *socks5.Request) func(ctx context.Context, network, addr string) (net.Conn, error) {
	fqdn := r.DestAddr.FQDN
	dest := r.RealDestAddr().IP.String()

	// 	dialGen := func(proxy string) func(ctx context.Context, network, addr string) (net.Conn, error) {
	// 		return net.Dial(network, addr)
	// 	}
	found := color.Sprint("@{g}FOUND@{|}")
	notFound := color.Sprint("@{r}NOTFOUND@{|}")

	var auth *proxy.Auth

	if c, ok := myGeoConfig.MIP[dest]; ok {
		log.Println(found, "IP match", fqdn, dest, c.Upstream)
		if len(c.Upstream) == 0 {
			return func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial(network, addr)
			}
		}
		u := c.Upstream[0]
		if len(u.User) > 0 && len(u.Password) > 0 {
			auth = &proxy.Auth{User: u.User, Password: u.Password}
		}
		dialer, err := proxy.SOCKS5("tcp", u.Addr,
			auth,
			&net.Dialer{
				KeepAlive: 30 * time.Second,
			},
		)
		if err != nil {
			log.Println(err.Error())
		}
		return func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
			// 			return net.Dial(network, addr)
		}
		// 		return c.Upstream[0].dial
	}
	for k, c := range myGeoConfig.MCIDR {
		if subnettool.CIDRMatch(dest, k) {
			log.Println(found, "CIDR match", fqdn, dest, c.Upstream)
			if len(c.Upstream) == 0 {
				return func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.Dial(network, addr)
				}
			}
			u := c.Upstream[0]
			if len(u.User) > 0 && len(u.Password) > 0 {
				auth = &proxy.Auth{User: u.User, Password: u.Password}
			}
			dialer, err := proxy.SOCKS5("tcp", u.Addr,
				auth,
				&net.Dialer{
					KeepAlive: 30 * time.Second,
				},
			)
			if err != nil {
				log.Println(err.Error())
			}
			return func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
				// 				return net.Dial(network, addr)
			}
		}
	}
	if c, ok := myGeoConfig.MFQDN[fqdn]; ok {
		log.Println(found, "FQDN match", fqdn, dest, c.Upstream)
		if len(c.Upstream) == 0 {
			return func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial(network, addr)
			}
		}
		u := c.Upstream[0]
		if len(u.User) > 0 && len(u.Password) > 0 {
			auth = &proxy.Auth{User: u.User, Password: u.Password}
		}
		dialer, err := proxy.SOCKS5("tcp", u.Addr,
			auth,
			&net.Dialer{
				KeepAlive: 30 * time.Second,
			},
		)
		if err != nil {
			log.Println(err.Error())
		}
		return func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
			// 			return net.Dial(network, addr)
		}
	}
	for k, c := range myGeoConfig.MREFQDN {
		if glob.Glob(k, fqdn) {
			log.Println(found, "FQDN match", fqdn, dest, c.Upstream)
			if len(c.Upstream) == 0 {
				return func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.Dial(network, addr)
				}
			}
			u := c.Upstream[0]
			if len(u.User) > 0 && len(u.Password) > 0 {
				auth = &proxy.Auth{User: u.User, Password: u.Password}
			}
			dialer, err := proxy.SOCKS5("tcp", u.Addr,
				auth,
				&net.Dialer{
					KeepAlive: 30 * time.Second,
				},
			)
			if err != nil {
				log.Println(err.Error())
			}
			return func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
				// 				return net.Dial(network, addr)
			}
		}
	}
	country, _ := db.Country(net.ParseIP(dest))
	code := color.Sprintf("@{c}%s@{|}", country.Country.IsoCode)
	if c, ok := myGeoConfig.MGEO[country.Country.IsoCode]; ok {
		log.Println(found, "GEO match", fqdn, dest, code, c.Upstream)
		if len(c.Upstream) == 0 {
			return func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial(network, addr)
			}
		}
		u := c.Upstream[0]
		if len(u.User) > 0 && len(u.Password) > 0 {
			auth = &proxy.Auth{User: u.User, Password: u.Password}
		}
		dialer, err := proxy.SOCKS5("tcp", u.Addr,
			auth,
			&net.Dialer{
				KeepAlive: 30 * time.Second,
			},
		)
		if err != nil {
			log.Println(err.Error())
		}
		return func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
			// 			return net.Dial(network, addr)
		}
	}

	log.Println(notFound, "GEO", code, "use default", fqdn, dest)
	return func(ctx context.Context, net_, addr string) (net.Conn, error) {
		return net.Dial(net_, addr)
	}
}
