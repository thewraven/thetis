package main

import (
	"fmt"
	"net/http"

	"github.com/thewraven/thetis"
)

func createServer() *http.ServeMux {
	var mux http.ServeMux
	home := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello,world!")
	}
	bang := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Not working")
	}
	tea := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		fmt.Fprintln(w, "I'm a teapot")
	}

	mux.HandleFunc("/", home)
	mux.HandleFunc("/bang", bang)
	mux.HandleFunc("/tea", tea)

	return &mux
}

func main() {
	//initialize your http server
	server := createServer()

	//create a monitor registering the default http monitor
	//You can add custom monitors using the Add(thetis.Metric) method
	monitorServer := thetis.NewWithDefaults(server)

	//register all the http metrics via prometheus.MustRegister
	//returns the prometheus http.Handler that informs about
	//the collected metrics
	prometheusHandler := monitorServer.RegisterAll()

	//mount the prometheusHandler in the endpoint of your preference
	server.Handle("/metrics", prometheusHandler)

	//use the monitorServer as a regular http.Handler
	port := ":7070"
	fmt.Println("Server listening at ", port)
	http.ListenAndServe(port, monitorServer)
}
