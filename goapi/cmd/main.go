package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	var pingcounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ping_request_count",
			Help: "No of ping requests",
		},
	)

	e := echo.New()
	e.GET("/ping", func(c echo.Context) error {
		pingcounter.Inc()
		return c.String(http.StatusOK, fmt.Sprintf("%v pong", time.Now().Format(time.UnixDate)))
	})

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	prometheus.MustRegister(pingcounter)

	e.Logger.Fatal(e.Start(":1323"))
}
