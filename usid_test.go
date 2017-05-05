package usid

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"testing/quick"
	"time"
)

func TestTimestampEquality(t *testing.T) {
	var (
		entropy  = RndEntropy()
		expected = time.Now()
		id       = MustNew(Timestamp(expected), entropy)
		actual   = time.Unix(0, int64(id.Timestamp()))
	)

	if !expected.Equal(actual) {
		t.Errorf("Expected: %v, actual: %v", expected, actual)
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	runner := func(reader io.Reader) error {
		var (
			num = 1000
			res = make([]USID, num, num)
		)
		for i := 0; i < num; i++ {
			res[i] = MustNew(Timestamp(time.Now()), reader)
		}

		// check for duplications
		for k0, v0 := range res {
			for k1, v1 := range res {
				if k0 == k1 {
					continue
				}

				if bytes.Equal(v0[:], v1[:]) {
					return errors.New("match")
				}
			}
		}
		return nil
	}

	// Test that using a source doesn't generate the same id consecutively.
	t.Run("MachEntropy", func(t *testing.T) {
		if err := runner(MachEntropy()); err != nil {
			t.Error(err)
		}
	})

	t.Run("RndEntropy", func(t *testing.T) {
		if err := runner(RndEntropy()); err != nil {
			t.Error(err)
		}
	})

	t.Run("SecRndEntropy", func(t *testing.T) {
		if err := runner(SecRndEntropy()); err != nil {
			t.Error(err)
		}
	})

	t.Run("CryptoRndEntropy", func(t *testing.T) {
		if err := runner(CryptoRndEntropy()); err != nil {
			t.Error(err)
		}
	})
}

func TestTime(t *testing.T) {
	t.Parallel()

	// Test Timestamp() returns correct time after encoding.
	t.Run("Timestamp()", func(t *testing.T) {
		fn := func() bool {
			var (
				stamp = time.Now()
				id    = MustNew(Timestamp(stamp), SecRndEntropy())
			)
			return id.Timestamp() == Timestamp(stamp)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	clamp := func(stamp uint64) uint64 {
		return stamp % maxTimestamp
	}

	// Test SetTime() returns correct time after encoding.
	t.Run("SetTime()", func(t *testing.T) {
		fn := func(stamp uint64) bool {
			var (
				id USID

				offset = clamp(stamp)
			)
			if err := id.SetTimestamp(offset); err != nil {
				t.Fatal(err)
			}

			return id.Timestamp() == offset
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestCompare(t *testing.T) {
	t.Parallel()

	// Test comparing the same id
	t.Run("Identical", func(t *testing.T) {
		fn := func() bool {
			var (
				stamp = time.Now()
				id    = MustNew(Timestamp(stamp), SecRndEntropy())
			)
			return id.Compare(id) == 0
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	// clamp makes sure that our duration isn't enormous.
	clamp := func(dur time.Duration) time.Duration {
		return dur % (time.Hour * 100)
	}
	// abs makes sure that we always have a positive duration.
	abs := func(dur time.Duration) time.Duration {
		if dur < 0 {
			return -dur
		}
		if dur == 0 {
			return 0
		}
		return dur
	}

	// Test comparing id is less
	t.Run("Less", func(t *testing.T) {
		fn := func(dur time.Duration) bool {
			var (
				offset  = clamp(abs(dur))
				stamp   = time.Now()
				entropy = SecRndEntropy()

				a = MustNew(Timestamp(stamp.Add(-offset)), entropy)
				b = MustNew(Timestamp(stamp), entropy)
			)

			return a.Compare(b) == -1
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	// Test comparing id is more
	t.Run("More", func(t *testing.T) {
		fn := func(dur time.Duration) bool {
			var (
				offset  = clamp(abs(dur))
				stamp   = time.Now()
				entropy = RndEntropy()

				a = MustNew(Timestamp(stamp), entropy)
				b = MustNew(Timestamp(stamp.Add(-offset)), entropy)
			)

			return a.Compare(b) == 1
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestEntropy(t *testing.T) {
	t.Parallel()

	t.Run("SetEntropy", func(t *testing.T) {
		fn := func(b [10]byte) bool {
			var id USID
			if err := id.SetEntropy(b); err != nil {
				t.Fatal(err)
			}

			expected, actual := b, id.Entropy()
			return bytes.Equal(expected[:], actual[:])
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	// Test comparing sunsequent reads don't lead to identical bytes.
	t.Run("RndEntropy", func(t *testing.T) {
		fn := func() bool {
			var (
				a, b    [10]byte
				entropy = RndEntropy()
			)
			if _, err := entropy.Read(a[:]); err != nil {
				t.Fatal(err)
			}
			if _, err := entropy.Read(b[:]); err != nil {
				t.Fatal(err)
			}
			return !bytes.Equal(a[:], b[:])
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	// Test comparing sunsequent reads don't lead to identical bytes.
	t.Run("SecRndEntropy", func(t *testing.T) {
		fn := func() bool {
			var (
				a, b    [10]byte
				entropy = SecRndEntropy()
			)
			if _, err := entropy.Read(a[:]); err != nil {
				t.Fatal(err)
			}
			if _, err := entropy.Read(b[:]); err != nil {
				t.Fatal(err)
			}
			return !bytes.Equal(a[:], b[:])
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	// Test comparing sunsequent reads don't lead to identical bytes.
	t.Run("MachEntropy", func(t *testing.T) {
		fn := func() bool {
			var (
				a, b    [10]byte
				entropy = MachEntropy()
			)
			if _, err := entropy.Read(a[:]); err != nil {
				t.Fatal(err)
			}
			if _, err := entropy.Read(b[:]); err != nil {
				t.Fatal(err)
			}
			return !bytes.Equal(a[:], b[:])
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestString(t *testing.T) {
	t.Parallel()

	t.Run("String", func(t *testing.T) {
		fn := func() bool {
			var (
				stamp   = time.Now()
				entropy = SecRndEntropy()

				id = MustNew(Timestamp(stamp), entropy)
			)

			return len(id.String()) == 36
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}
