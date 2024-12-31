package solution

import (
	"log"
	"time"
)

func Track(start time.Time, msg string) {
	elapsed := time.Since(start)
	log.Printf("%v: %v\n", msg, elapsed)
	log.Printf("%v in seconds: %.6f\n", msg, elapsed.Seconds())
}
