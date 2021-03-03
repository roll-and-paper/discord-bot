package dices

import (
	"math/rand"
	"sync"
	"time"
)

func SetRandomSeed(newSeed int64) { random.Seed(newSeed) }

var random = rand.New(&internalRand{src: rand.NewSource(time.Now().UnixNano())})

type internalRand struct {
	src rand.Source
	lk  sync.Mutex
}

func (r *internalRand) Int63() int64 {
	r.lk.Lock()
	defer r.lk.Unlock()
	return r.src.Int63()
}

func (r *internalRand) Seed(seed int64) {
	r.lk.Lock()
	defer r.lk.Unlock()
	r.src.Seed(seed)
}
