package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	delay                = flag.Int64("delay", 50, "in milliseconds")
	port                 = "3333"
	prometheusSummaryVec *prometheus.SummaryVec
)

func init() {
	// Construct logging.
	flag.Parse()

	prometheusSummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "C",
		Name:      "handler_request_milisecond",
		Help:      "Average of handler response time",
	}, []string{"handler", "method", "httpCode"})
	if err := prometheus.Register(prometheusSummaryVec); err != nil {
		log.Printf("Failed to register prometheus metrics: %s", err.Error())
	}
}

func main() {
	http.Handle("/metrics", prometheus.Handler())

	http.HandleFunc("/world", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		r = r.WithContext(ctx)
		httpCode := "500"

		defer func(timeStart time.Time) {
			pattern := "/world"
			method := "GET"

			// prometheus
			if prometheusSummaryVec != nil {
				prometheusSummaryVec.With(prometheus.Labels{"handler": pattern, "method": method, "httpCode": httpCode}).Observe(time.Since(timeStart).Seconds() * 1000)
			}
		}(time.Now())

		stopChannel := time.After(time.Millisecond * time.Duration(*delay))
	myLoop:
		for {
			select {
			case <-stopChannel:
				break myLoop
			case <-ctx.Done():
				httpCode = "408"
				return
			default:
				time.Sleep(time.Microsecond * 5)
			}
		}

		httpCode = "200"
		fmt.Fprintf(w, "world")
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
