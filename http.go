/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : http.go

* Purpose :

* Creation Date : 05-15-2017

* Last Modified : Mon May 15 19:47:08 2017

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"fmt"
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
        return "SOCKS5 {{.}}; SOCKS {{.}}";
    }
});`

func handler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("s")
	if len(s) == 0 {
		fmt.Fprintf(w, "ok")
		return
	}
	s += ":" + strings.Split(*fListen, ":")[1]
	t, err := template.New("pac").Parse(TPL)
	if err != nil {
		fmt.Fprintf(w, "ok")
		return
	}
	t.Execute(w, s)
}

func runHttp() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":1081", nil)
}
