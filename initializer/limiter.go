package initializer

import (
	"os"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

var Limiter *rate.Limiter = nil

func InitializeLimiter() {
	requestsPerSecond := 10
	if rps := os.Getenv("RATE_LIMIT_RPS"); rps != "" {
		if parsed, err := strconv.Atoi(rps); err == nil {
			requestsPerSecond = parsed
		}
	}

	Limiter = rate.NewLimiter(rate.Every(time.Second), requestsPerSecond)
}
