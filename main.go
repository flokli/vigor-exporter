package main

import (
	"flag"
	"github.com/T4cC0re/vigor-node-exporter/Vigor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

var username = flag.String("username", "", "username to authenticate to the Vigor")
var password = flag.String("password", "", "password to authenticate to the Vigor")
var ip = flag.String("ip", "", "ip the Vigor is reachable on")

var vigor *Vigor.Vigor

func loginIfError(err error) {
	if err != nil {
		print(err)
		vigor.Login(*username, *password)
	}
}

func main() {
	flag.Parse()

	var err error
	vigor, err = Vigor.New(*ip)
	if err != nil {
		panic(err)
	}
	vigor.Login(*username, *password)

	vigor.UpdateStatus()
	vigor.FetchStatus()

	go func() {
		for {
			time.Sleep(5 * time.Second)

			loginIfError(vigor.UpdateStatus())
			loginIfError(vigor.FetchStatus())
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9103", nil))
}
