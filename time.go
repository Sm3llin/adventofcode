package adventofcode

import (
	"fmt"
	"time"
)

func Time(f func()) {
	start := time.Now()
	f()
	fmt.Println(time.Since(start))
}
