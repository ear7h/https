package main

import (
	"flag"
	"net/http/httputil"
	"net/url"
	"net/http"
)

var port = flag.String("d", "8080", "destination port number")
var host = flag.String("h", "", "website we're serving")
var cert = flag.String("c", "", "cert location")
var key = flag.String("k", "", "key location")

func main() {
	flag.Parse()

	if len(*port) * len(*host) == 0 {
		flag.Usage()
		return
	}

	u, err := url.Parse("http://localhost:" + *port)
	if err != nil {
		panic(err)
	}

	errc := make(chan error, 1)

	// forward to https
	go func() {
		errc <- http.ListenAndServe(":80", http.RedirectHandler("https://"+*host+":443", http.StatusPermanentRedirect))
	}()

	//serve https
	go func() {
		rprox := httputil.NewSingleHostReverseProxy(u)
		d1 := rprox.Director
		d2 := func(req *http.Request) {
			d1(req)
			req.Header.Add("X-Forwarded-Proto", "https")
		}

		rprox.Director = d2

		errc <- http.ListenAndServeTLS(":443", *cert, *key, rprox)
	}()

	err = <- errc
	panic(err)
}
