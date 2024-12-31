package solution

import (
	"log"
	"time"
)

func Track(msg string) func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		log.Printf("%v: %v\n", msg, elapsed)
		log.Printf("%v in seconds: %.6f\n", msg, elapsed.Seconds())
	}
}
