package usid

import (
	"io"
	"math"
	"testing"
	"time"
)

func TestConcurrency(t *testing.T) {
	t.Parallel()

	dispatch := func(reader io.Reader, count int) error {
		for i := 0; i < count; i++ {
			if _, err := New(Timestamp(time.Now()), reader); err != nil {
				return err
			}
		}
		return nil
	}
	runner := func(reader io.Reader, total int) error {
		var (
			num   = int(math.Sqrt(float64(total)))
			share = total / num

			errC = make(chan error, num)
		)

		for i := 0; i < num; i++ {
			go func() {
				errC <- dispatch(reader, share)
			}()
		}

		for i := 0; i < num; i++ {
			if err := <-errC; err != nil {
				return err
			}
		}

		return nil
	}

	// Test that using a source doesn't generate the same id consecutively.
	t.Run("MachEntropy", func(t *testing.T) {
		if err := runner(MachEntropy(), 10000); err != nil {
			t.Error(err)
		}
	})

	// Note: we can't use RndEntropy here as it causes race issues.

	t.Run("SecRndEntropy", func(t *testing.T) {
		if err := runner(SecRndEntropy(), 10000); err != nil {
			t.Error(err)
		}
	})

	t.Run("CryptoRndEntropy", func(t *testing.T) {
		if err := runner(CryptoRndEntropy(), 10000); err != nil {
			t.Error(err)
		}
	})
}
