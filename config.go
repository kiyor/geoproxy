/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : config.go

* Purpose :

* Creation Date : 05-14-2017

* Last Modified : Sun May 21 02:24:26 2017

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	// 	"context"
	"encoding/json"
	"github.com/kiyor/golib"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"strings"
	"sync"
	// 	"time"
)

var myGeoConfig *GeoConfig

func init() {
	err := Reload(*fConf)
	if err != nil {
		panic(err)
	}
}

// default will go to ""{}
type Upstream map[string]*UpstreamConfig

type UpstreamConfig struct {
	Upstream []*Up
}

type Up struct {
	Addr     string
	User     string
	Password string
	auth     *proxy.Auth
}

type geoConfig map[string]string

type CacheKey struct {
	IP   string
	FQDN string
}

type GeoConfig struct {
	Default *UpstreamConfig
	MIP     map[string]*UpstreamConfig
	MCIDR   map[string]*UpstreamConfig
	MFQDN   map[string]*UpstreamConfig
	MREFQDN map[string]*UpstreamConfig
	MGEO    map[string]*UpstreamConfig
	Cache   map[CacheKey]*UpstreamConfig
	*sync.RWMutex
}

func Reload(dir string) error {
	geo, err := LoadConfig(dir)
	if err != nil {
		return err
	}
	myGeoConfig = geo
	return nil
}

func LoadConfig(dir string) (*GeoConfig, error) {
	upstream_f := filepath.Join(dir, "upstream.json")
	b, err := ioutil.ReadFile(upstream_f)
	if err != nil {
		return nil, err
	}
	upstream := make(Upstream)
	err = golib.JsonUnmarshal(b, &upstream)
	if err != nil {
		return nil, err
	}

	for k, v := range upstream {
		for _, u := range v.Upstream {
			if len(u.User) > 0 && len(u.Password) > 0 {
				u.auth = &proxy.Auth{u.User, u.Password}
			}
		}
		if len(v.Upstream) == 0 {
			upstream[k].Upstream = append(upstream[k].Upstream, &Up{"", "", "", nil})
		}
	}

	/*
		for k, ups := range upstream {
			for _, up := range ups.Upstream {
				if len(up.Addr) == 0 {
					up.dial = func(ctx context.Context, net_, addr string) (net.Conn, error) {
						return net.Dial(net_, addr)
					}
					continue
				}
				// TODO: add support for more config
				log.Println(up.Addr)
				dialer, err := proxy.SOCKS5("tcp", up.Addr,
					nil,
					&net.Dialer{
						KeepAlive: 30 * time.Second,
					},
				)
				if err != nil {
					log.Println(err.Error())
				}
				up.dial = func(ctx context.Context, net_, addr string) (net.Conn, error) {
					log.Println("dial", addr)
					return dialer.Dial(net_, addr)
				}
				log.Println(k, up.dial)
			}
		}
	*/

	geo_f := filepath.Join(dir, "geo.json")
	b, err = ioutil.ReadFile(geo_f)
	if err != nil {
		return nil, err
	}
	geo := make(geoConfig)
	err = golib.JsonUnmarshal(b, &geo)
	if err != nil {
		return nil, err
	}

	Geo := &GeoConfig{
		MIP:     make(map[string]*UpstreamConfig),
		MCIDR:   make(map[string]*UpstreamConfig),
		MFQDN:   make(map[string]*UpstreamConfig),
		MREFQDN: make(map[string]*UpstreamConfig),
		MGEO:    make(map[string]*UpstreamConfig),
		Cache:   make(map[CacheKey]*UpstreamConfig),
		RWMutex: new(sync.RWMutex),
	}
	for k, v := range geo {
		if k == "default" {
			Geo.Default = upstream[v]
			continue
		}

		up, ok := upstream[v]
		if !ok {
			log.Println("upstream", v, "not found, use default")
			up = Geo.Default
		}
		if ip := net.ParseIP(k); ip != nil {
			Geo.MIP[k] = up
		} else if _, _, err := net.ParseCIDR(k); err == nil {
			Geo.MCIDR[k] = up
		} else if strings.Contains(k, ".") {
			if !strings.Contains(k, "*") {
				Geo.MFQDN[k] = up
			} else {
				Geo.MREFQDN[k] = up
			}
		} else {
			Geo.MGEO[k] = up
		}
	}
	log.Println("config load success")

	return Geo, nil
}

func Json(i interface{}) string {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		log.Println(err.Error())
	}
	return string(b)
}
