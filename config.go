/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : config.go

* Purpose :

* Creation Date : 05-14-2017

* Last Modified : Mon 15 May 2017 05:34:05 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"context"
	"encoding/json"
	"github.com/kiyor/golib"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var myGeoConfig *GeoConfig

func init() {
	Reload("./conf")
}

// default will go to ""{}
type Upstream map[string]*UpstreamConfig

type UpstreamConfig struct {
	Upstream []Up
}

type Up struct {
	Addr string
	dial func(ctx context.Context, network, addr string) (net.Conn, error)
}

type geoConfig map[string]string

type GeoConfig struct {
	Default *UpstreamConfig
	MIP     map[string]*UpstreamConfig
	MCIDR   map[string]*UpstreamConfig
	MFQDN   map[string]*UpstreamConfig
	MGEO    map[string]*UpstreamConfig
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

	for _, ups := range upstream {
		for _, up := range ups.Upstream {
			if len(up.Addr) == 0 {
				up.dial = func(ctx context.Context, net_, addr string) (net.Conn, error) {
					return net.Dial(net_, addr)
				}
				continue
			}
			// TODO: add support for more config
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
				return dialer.Dial(net_, addr)
			}
		}
	}

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
		MGEO:    make(map[string]*UpstreamConfig),
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
			Geo.MFQDN[k] = up
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
