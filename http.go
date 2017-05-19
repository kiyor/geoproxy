/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : http.go

* Purpose :

* Creation Date : 05-15-2017

* Last Modified : Wed 17 May 2017 06:14:49 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
)

const TPL = `var FindProxyForURL = function(init, profiles) {
    return function(url, host) {
        "use strict";
        var result = init, scheme = url.substr(0, url.indexOf(":"));
        do {
            result = profiles[result];
            if (typeof result === "function") result = result(url, host, scheme);
        } while (typeof result !== "string" || result.charCodeAt(0) === 43);
        return result;
    };
}("+geoproxy", {
    "+geoproxy": function(url, host, scheme) {
        "use strict";
        if (/^127\.0\.0\.1$/.test(host) || /^::1$/.test(host) || /^localhost$/.test(host)) return "DIRECT";
        if (scheme !== "http" && scheme !== "https") return "DIRECT";
        return "SOCKS5 {{.}}";
    }
});`

func handler(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.Host, ":")[0]
	s += ":" + strings.Split(*fListen, ":")[1]
	// 	t, err := template.New("pac").Parse(TPL)
	t, err := template.ParseFiles("./pac.tpl")
	if err != nil {
		fmt.Fprintf(w, "ok")
		return
	}
	w.Header().Add("Cache-Control", "max-age=300")
	t.Execute(w, s)
	log.Println(r.RemoteAddr, r.Method, r.URL.String())
}

func runHttp() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":1081", nil)
}
