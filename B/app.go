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
	}, []string{"handler", "method", "httpCode"})
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
		httpCode := "500"

		defer func(timeStart time.Time) {
			pattern := "/hello"
			method := "GET"

			// prometheus
			if prometheusSummaryVec != nil {
				prometheusSummaryVec.With(prometheus.Labels{"handler": pattern, "method": method, "httpCode": httpCode}).Observe(time.Since(timeStart).Seconds() * 1000)
			}
		}(time.Now())

		client := http.DefaultClient
		req, err := http.NewRequest("GET", "http://service-c:3333/world", nil)
		if err != nil {
			httpCode = "500"
			log.Println(err)
			return
		}
		req = req.WithContext(ctx)
		errChan := make(chan error)
		respChan := make(chan *http.Response)
		go func() {
			resp, err := client.Do(req)
			if err != nil {
				errChan <- err
			} else {
				respChan <- resp
			}
		}()

		select {
		case err := <-errChan:
			httpCode = "500"
			log.Println(err)
		case resp := <-respChan:
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				httpCode = "500"
				return
			}
			text := "hello " + string(body)
			httpCode = "200"
			fmt.Fprintf(w, text)
		case <-ctx.Done():
			httpCode = "408"
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
