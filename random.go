package x

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"time"
)

// RandomNumber used to get a random number
func RandomNumber() (i int64) {
	i = rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
	return
}

// Random used to get a random number between [min, max)
func Random(min, max int64) int64 {
	if max <= min {
		return min
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63n(max-min) + min
}

func CryptoRandom(min, max int64) (int64, error) {
	if max <= min {
		return min, nil
	}
	bi, err := crand.Int(crand.Reader, big.NewInt(max-min))
	if err != nil {
		return 0, err
	}
	return bi.Int64() + min, nil
}
