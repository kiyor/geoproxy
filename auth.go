/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : auth.go

* Purpose :

* Creation Date : 05-15-2017

* Last Modified : Mon May 15 15:26:52 2017

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"bufio"
	"encoding/json"
	"github.com/kiyor/go-socks5"
	"io/ioutil"
	"os"
	"strings"
)

func parseSocks5Auth(input string) socks5.StaticCredentials {
	if strings.Contains(input, " ") {
		p := strings.Split(input, " ")
		return socks5.StaticCredentials{
			p[0]: p[1],
		}
	}
	if strings.Contains(input, ":") {
		p := strings.Split(input, ":")
		return socks5.StaticCredentials{
			p[0]: p[1],
		}
	}
	cred := make(socks5.StaticCredentials)
	d, err := ioutil.ReadFile(input)
	if err != nil {
		return socks5.StaticCredentials{}
	}
	err = json.Unmarshal(d, &cred)
	if err != nil {
		lines, err := cleanFile(input)
		if err != nil {
			return socks5.StaticCredentials{}
		}
		for _, line := range lines {
			p := strings.Split(line, " ")
			if len(p) > 1 {
				cred[p[0]] = p[1]
			}
		}
		return cred
	}

	return cred
}

func cleanFile(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()

	var line []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		p := strings.Split(scanner.Text(), "#")
		if len(p[0]) > 0 {
			line = append(line, p[0])
		}
	}

	if err := scanner.Err(); err != nil {
		return line, err
	}
	return line, nil
}
