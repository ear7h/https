package main

import (
	"flag"
	"net/http/httputil"
	"net/url"
	"net/http"
)

var port = flag.String("destination", "8080", "port number")
var host = flag.String("host", "", "website we're serving")

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

	go func() {
		errc <- http.ListenAndServe(":80", http.RedirectHandler("https://"+*host+":443", http.StatusPermanentRedirect))
	}()

	go func() {
		errc <- http.ListenAndServe(":443", httputil.NewSingleHostReverseProxy(u))
	}()

	err = <- errc
	panic(err)
}
