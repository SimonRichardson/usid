package usid

import "testing"

func BenchmarkNonCachedEntropy(b *testing.B) {
	b.Run("Baseline", func(b *testing.B) {
		b.SetBytes(int64(len(USID{})))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = New(0, nil)
		}
	})

	b.Run("MachEntropy", func(b *testing.B) {
		b.SetBytes(int64(len(USID{})))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = New(0, MachEntropy())
		}
	})

	b.Run("RndEntropy", func(b *testing.B) {
		b.SetBytes(int64(len(USID{})))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = New(0, RndEntropy())
		}
	})

	b.Run("CryptoRndEntropy", func(b *testing.B) {
		b.SetBytes(int64(len(USID{})))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = New(0, CryptoRndEntropy())
		}
	})
}

func BenchmarkCachedEntropy(b *testing.B) {
	b.Run("Baseline", func(b *testing.B) {
		b.SetBytes(int64(len(USID{})))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = New(0, nil)
		}
	})

	b.Run("MachEntropy", func(b *testing.B) {
		entropy := MachEntropy()

		b.SetBytes(int64(len(USID{})))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = New(0, entropy)
		}
	})

	b.Run("RndEntropy", func(b *testing.B) {
		entropy := RndEntropy()

		b.SetBytes(int64(len(USID{})))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = New(0, entropy)
		}
	})

	b.Run("CryptoRndEntropy", func(b *testing.B) {
		entropy := CryptoRndEntropy()

		b.SetBytes(int64(len(USID{})))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = New(0, entropy)
		}
	})
}
