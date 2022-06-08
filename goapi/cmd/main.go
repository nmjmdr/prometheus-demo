package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func stats_interceptor() echo.MiddlewareFunc {
	pingRequestCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ping_request_count",
			Help: "Ping. No of requests",
		},
	)

	pingLatencyGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ping_request_latency_guage",
			Help: "Ping. Request latency guage",
		},
	)

	pingLatencySum := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ping_request_latency_sum",
			Help: "Ping. Request latency sum",
		},
	)

	pingLatencyHistogram := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "ping_request_latency_histogram",
		Help:    "Ping. Request latency histogram",
		Buckets: []float64{20, 40, 60, 80, 100, 120, 140, 160, 180, 200},
	})

	prometheus.MustRegister(pingRequestCounter)
	prometheus.MustRegister(pingLatencyHistogram)
	prometheus.MustRegister(pingLatencyGauge)
	prometheus.MustRegister(pingLatencySum)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			pingRequestCounter.Inc()
			// timer := prometheus.NewTimer(pingLatencyHistogram)
			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}
			if c.Request().URL.Path == "/ping" {
				f := float64(time.Since(start).Milliseconds())
				//timer.ObserveDuration()
				pingLatencyHistogram.Observe(f)
				pingLatencyGauge.Set(f)
				pingLatencySum.Add(f)
				fmt.Println("Ping time taken: ", f)
			}
			return nil
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	max := 100
	maxDelay := 190
	e := echo.New()
	defaultDelay := 10
	varianceDelay := 5
	timeDelay := 0

	flag := false
	stickyIndex := 0
	StickyLatencyCalls := 5

	e.Use(stats_interceptor())
	e.GET("/ping", func(c echo.Context) error {
		start := time.Now()
		n := rand.Intn(max)

		timeDelay = defaultDelay + rand.Intn(varianceDelay)

		if n%10 == 0 {
			timeDelay = maxDelay
			flag = true
		}

		if flag {
			stickyIndex++
		}

		if stickyIndex == StickyLatencyCalls {
			flag = false
			stickyIndex = 0
			timeDelay = 0
		}

		time.Sleep(time.Duration(timeDelay) * time.Millisecond)
		return c.String(http.StatusOK, fmt.Sprintf("%v pong, time taken (ms): %v",
			time.Now().Format(time.UnixDate), time.Since(start)))
	})

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.Logger.Fatal(e.Start(":1323"))
}
