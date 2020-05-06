package random

import (
	"math/rand"
	"time"
)

const s = 'a'
const c = 'z' - 'a' + 1

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

func Letter(len int) string {
	println(c, s)
	data := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(c) + s
		data[i] = byte(b)
	}
	return string(data)
}
