package ant

import(
    "time"
    "math/rand"
)

func RandomNumber()(i int64){
    i = rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
    return
}

