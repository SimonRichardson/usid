package usid

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync/atomic"
	"time"
)

/*
A USID is a Unique Sortable Identifier
*/
type USID [18]byte

const (
	defaultEntropyOffset = 8
)

var (
	// ErrBigTime is returned when constructing an USID with a time that is larger
	// than MaxTime.
	ErrBigTime = errors.New("usid: time too big")
)

// New returns an USID with given Unix timestamp and an optional entropy source.
// Use the Timestamp function to convert time.Time to Unix timestamp.
func New(stamp uint64, entropy io.Reader) (id USID, err error) {
	if err = id.SetTimestamp(stamp); err != nil {
		return
	}

	if entropy != nil {
		_, err = entropy.Read(id[defaultEntropyOffset:])
	}

	return
}

// MustNew creates a new USID that panics on failure during creation.
func MustNew(stamp uint64, entropy io.Reader) USID {
	id, err := New(stamp, entropy)
	if err != nil {
		panic(err)
	}
	return id
}

// Timestamp returns the Unix time encoded in the USID
func (u USID) Timestamp() uint64 {
	return uint64(u[7]) | uint64(u[6])<<8 |
		uint64(u[5])<<16 | uint64(u[4])<<24 |
		uint64(u[3])<<32 | uint64(u[2])<<40 |
		uint64(u[1])<<48 | uint64(u[0])<<56
}

var maxTimestamp = USID{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}.Timestamp()

// SetTimestamp sets the timestamp component of the USID to the given Unix
// timestamp
func (u *USID) SetTimestamp(stamp uint64) error {
	if stamp > maxTimestamp {
		return ErrBigTime
	}

	(*u)[0] = byte(stamp >> 56)
	(*u)[1] = byte(stamp >> 48)
	(*u)[2] = byte(stamp >> 40)
	(*u)[3] = byte(stamp >> 32)
	(*u)[4] = byte(stamp >> 24)
	(*u)[5] = byte(stamp >> 16)
	(*u)[6] = byte(stamp >> 8)
	(*u)[7] = byte(stamp)

	return nil
}

// Entropy returns the encoded entropy with in the USID
func (u USID) Entropy() [10]byte {
	var b [10]byte
	for i := 0; i < 10; i++ {
		b[i] = u[defaultEntropyOffset+i]
	}
	return b
}

// SetEntropy set the USID Entropy to the passed byte slice.
func (u *USID) SetEntropy(b [10]byte) error {
	copy((*u)[defaultEntropyOffset:], b[:])
	return nil
}

// Compare returns an integer comparing id and other lexicographically.
// The result will be 0 if id == other, -1 if id < other and +1 if id > other.
func (u USID) Compare(other USID) int {
	return bytes.Compare(u[:], other[:])
}

func (u USID) String() string {
	return fmt.Sprintf("%x", string(u[:]))
}

// Timestamp concerts a time.Time into a Unix timestamp that USID can utilise.
func Timestamp(t time.Time) uint64 {
	return uint64(t.UnixNano())
}

// RndEntropy returns a random source of entropy for the creation of a USID
func RndEntropy() io.Reader {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// SecRndEntropy returns a random source that prevents issues with dead locks
// (race detector).
func SecRndEntropy() io.Reader {
	return newLockedRand(&lockedSource{
		src: rand.NewSource(time.Now().UnixNano()),
	})
}

// CryptoRndEntropy returns a cryptographic random source of entropy for the
// creation of a USID
func CryptoRndEntropy() io.Reader {
	return crand.Reader
}

// MachEntropy returns entropy that can be used to prevent collisions, yet be
// sortable in a more predictable manor.
func MachEntropy() io.Reader {
	return machReader{}
}

// processID and counter
var (
	// processID grabs application process id to ensure uniqueness between
	// multiple applications running on the same hardware.
	processID = os.Getpid()

	// Counter generates an atomically incremented integer. To prevent
	// collisions on different machines, a random offset it chosen first.
	counter = uint64(rand.New(rand.NewSource(time.Now().UnixNano())).Uint32())
)

// machReader abuses the io.Reader interface to allow the unique injection of
// machine information.
type machReader struct{}

func (m machReader) Read(b []byte) (n int, err error) {
	b[0] = byte(processID >> 8)
	b[1] = byte(processID)

	c := atomic.AddUint64(&counter, 1)
	binary.LittleEndian.PutUint64(b[2:], c)

	return 10, nil
}
