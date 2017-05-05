# usid

Unique sortable identifier

## Introduction

The USID creates a unique identifier using Unix timestamp (nanoseconds in this
case) and an entropy source. With the use of time, we naturally get a
sortable identifier as long as the machine times are synchronized (see `Compare`
method of USID).

Various types of entropy can be used to improve sorting if we want more high
grained positioning (see MachEntropy) or to improve uniqueness then something
like CryptoRndEntropy could be used.

## Types of entropy

The USID allows the creation with functions of entropies to best fit the
requirements of uniqueness.

#### MachEntropy

MachEntropy (machine entropy) uses two unique identifiers, the process id of the
application and a increment only atomic counter. The offset of the counter is
chosen at bootup to prevent conflicts between nodes.

#### RndEntropy

RndEntropy (random entropy) uses a random source using the time of now as the
seed.

#### SecRndEntropy

SecRndEntropy (secure random entropy) is the same as the random source, except
it uses mutexes to prevent data races when used concurrently.

#### CryptoRndEntropy

CryptoRndEntropy (cryptographic random entropy) uses a much more robust random
source for generating unique entropy.


## Example

```go
id := usid.MustNew(usid.Timestamp(time.Now()), usid.SecRndEntropy())
```

### Benchmarking

To run benchmarking, use the following command `go test -v -bench=. *.go -benchmem`

On my machine it outputs the following table, so the baseline is when we assume
that there is no entropy to use and in fact most if not all the ids will be
identical.

```
BenchmarkEntropy/Baseline-12         	30000000	        63.3 ns/op	 252.67 MB/s	      16 B/op	       1 allocs/op
BenchmarkEntropy/MachEntropy-12      	20000000	        85.4 ns/op	 187.40 MB/s	      16 B/op	       1 allocs/op
BenchmarkEntropy/RndEntropy-12       	10000000	       132 ns/op	 120.32 MB/s	      16 B/op	       1 allocs/op
BenchmarkEntropy/SecRndEntropy-12    	10000000	       131 ns/op	 121.22 MB/s	      16 B/op	       1 allocs/op
BenchmarkEntropy/CryptoRndEntropy-12 	 1000000	      1073 ns/op	  14.91 MB/s	      16 B/op	       1 allocs/op
```

### Commands

Possible series of commands to run:

```
go test -v *.go
go test -v -race *.go
go test -v -bench=. *.go -benchmem
```
