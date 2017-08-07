package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/afex/hystrix-go/hystrix"
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
		hystrix.ConfigureCommand("getWorld", hystrix.CommandConfig{
			Timeout:               1500,
			MaxConcurrentRequests: 100,
			ErrorPercentThreshold: 25,
		})
		respChan := make(chan *http.Response)
		errChan := hystrix.Go("getWorld", func() error {
			req, err := http.NewRequest("GET", "http://service-c:3333/world", nil)
			if err != nil {
				return err
			}
			req = req.WithContext(ctx)
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			respChan <- resp

			return nil
		}, nil)
		select {
		case <-ctx.Done():
			httpCode = "408"
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
		case err := <-errChan:
			if err == nil {
				httpCode = "200"
			} else if err == hystrix.ErrCircuitOpen || err == hystrix.ErrMaxConcurrency {
				httpCode = "503"
			} else if err == hystrix.ErrTimeout {
				httpCode = "408"
			} else {
				httpCode = "500"
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
			}
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
