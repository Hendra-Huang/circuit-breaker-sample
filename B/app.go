package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	port                 = "2222"
	prometheusSummaryVec *prometheus.SummaryVec
)

func init() {
	// Construct logging.
	flag.Parse()

	prometheusSummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "B",
		Name:      "handler_request_milisecond",
		Help:      "Average of handler response time",
	}, []string{"handler", "method"})
	if err := prometheus.Register(prometheusSummaryVec); err != nil {
		log.Printf("Failed to register prometheus metrics: %s", err.Error())
	}
}

func main() {
	http.Handle("/metrics", prometheus.Handler())

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET")
		ctx := r.Context()
		r = r.WithContext(ctx)

		defer func(timeStart time.Time) {
			pattern := "/hello"
			method := "GET"

			// prometheus
			if prometheusSummaryVec != nil {
				prometheusSummaryVec.With(prometheus.Labels{"handler": "all", "method": method}).Observe(time.Since(timeStart).Seconds() * 1000)
				prometheusSummaryVec.With(prometheus.Labels{"handler": pattern, "method": method}).Observe(time.Since(timeStart).Seconds() * 1000)
			}
		}(time.Now())

		client := http.DefaultClient
		req, err := http.NewRequest("GET", "http://service-c:3333/world", nil)
		req = req.WithContext(ctx)
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		text := "hello " + string(body)
		fmt.Fprintf(w, text)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
