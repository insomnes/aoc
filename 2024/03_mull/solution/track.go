package solution

import (
	"log"
	"time"
)

func Track(msg string) func() {
	start := time.Now()
	return func() {
		log.Printf("%v: %v\n", msg, time.Since(start))
	}
}
