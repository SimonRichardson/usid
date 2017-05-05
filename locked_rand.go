package usid

import (
	"errors"
	"math/rand"
	"sync"
)

// LockedRand allows random to be used in concurrent cases, where the default
// rand.Rand does not protect against race conditions.
//
// Note: this code is lifted from stdlib, but is not available to the user!?
type lockedRand struct {
	src     rand.Source
	readVal int64
	readPos int8
}

func newLockedRand(src rand.Source) *lockedRand {
	return &lockedRand{src: src}
}

func (r *lockedRand) Read(p []byte) (n int, err error) {
	if lk, ok := r.src.(*lockedSource); ok {
		return lk.read(p, &r.readVal, &r.readPos)
	}
	return 0, errors.New("invalid source")
}

func read(p []byte, int63 func() int64, readVal *int64, readPos *int8) (n int, err error) {
	pos := *readPos
	val := *readVal
	for n = 0; n < len(p); n++ {
		if pos == 0 {
			val = int63()
			pos = 7
		}
		p[n] = byte(val)
		val >>= 8
		pos--
	}
	*readPos = pos
	*readVal = val
	return
}

type lockedSource struct {
	mutex sync.Mutex
	src   rand.Source
}

func (r *lockedSource) Int63() (n int64) {
	r.mutex.Lock()
	n = r.src.Int63()
	r.mutex.Unlock()
	return
}

func (r *lockedSource) Seed(seed int64) {
	r.mutex.Lock()
	r.src.Seed(seed)
	r.mutex.Unlock()
}

func (r *lockedSource) read(p []byte, readVal *int64, readPos *int8) (n int, err error) {
	r.mutex.Lock()
	n, err = read(p, r.src.Int63, readVal, readPos)
	r.mutex.Unlock()
	return
}
