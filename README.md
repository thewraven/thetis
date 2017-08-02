# thetis
A simple ready-to-use http.Handler to register http-related metrics, using prometheus.

## Usage 

```go

func main() {
	//initialize your http server
	var server *http.ServeMux = createServer()

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
	http.ListenAndServe(":7070", monitorServer)
}

```

With this code, you can request /metrics for obtaining simple http information related to the requests.

See example/server.go for the full working example.
