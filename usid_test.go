package usid

import (
	"bytes"
	"testing"
	"testing/quick"
	"time"
)

func TestTime(t *testing.T) {
	t.Parallel()

	// Test Timestamp() returns correct time after encoding.
	t.Run("Timestamp()", func(t *testing.T) {
		fn := func() bool {
			var (
				stamp = time.Now()
				id    = MustNew(Timestamp(stamp), RndEntropy())
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
				id    = MustNew(Timestamp(stamp), RndEntropy())
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
				entropy = RndEntropy()

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
}

func TestString(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		fn := func() bool {
			var (
				stamp   = time.Now()
				entropy = RndEntropy()

				id = MustNew(Timestamp(stamp), entropy)
			)

			return len(id.String()) == 32
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}
