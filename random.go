package ant

import (
	"math/rand"
	"time"
)

func RandomNumber() (i int64) {
	i = rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
	return
}
