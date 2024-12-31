package solution

import (
	"log"
	"time"
)

func Track(start time.Time, msg string) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}
