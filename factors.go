package intfact

import (
	"errors"
	"math/big"
)

// Status describes the composite property
type Status int

// Status values
const (
	Unknown Status = iota
	ProbPrime
	Composite
	Prime
)

var (
	bigOne = big.NewInt(1)
)

var (
	errGcdIsN = errors.New("no factor found")
)

type gcdtest struct {
	n      *big.Int
	it     int
	period int
	acc    *big.Int
}

func newGcdtest(n *big.Int, period int) *gcdtest {
	return &gcdtest{
		n:      n,
		period: period,
		acc:    new(big.Int).Set(bigOne),
	}
}

func (t *gcdtest) test(a *big.Int) (fac *big.Int, err error) {
	t.it++
	t.acc.Mul(t.acc, a)
	t.acc.Mod(t.acc, t.n)
	if t.it >= t.period {
		fac, err = t.finish()
		if fac != nil || err != nil {
			return
		}
		t.it = 0
		t.acc.Set(bigOne)
	}
	return
}

func (t *gcdtest) finish() (fac *big.Int, err error) {
	d := new(big.Int).GCD(nil, nil, t.acc, t.n)
	if d.Cmp(bigOne) != 0 {
		if d.Cmp(t.n) == 0 {
			fac, err = nil, errGcdIsN
		} else {
			fac, err = d, nil
		}
	}
	return
}

// Fact contains information about the factors
type Fact struct {
	Fac  *big.Int
	Exp  uint
	Stat Status
	Next *Fact
}

// Factors contains the list of factors in increasing order
type Factors struct {
	// list of factors
	First *Fact
	// bound used for trial division. Factors < PBound^2 must be prime.
	PBound *big.Int
}

// NewFactors creates a fresh Factors structure for the number a.
func NewFactors(a *big.Int) *Factors {
	f := Fact{
		Fac:  a,
		Exp:  1,
		Stat: Unknown,
	}
	l := Factors{
		First:  &f,
		PBound: big.NewInt(1),
	}
	return &l
}

// RecordSplit removes fp and inserts new factors a and b.
// Since the new factors are smaller than fp, the list starting with
// *fp.Next is unchanged.
func (l *Factors) RecordSplit(fp **Fact, a, b *big.Int) {
	// remove the factor from the list
	f := *fp
	*fp = f.Next
	// create and insert new factors
	bnd := new(big.Int).Mul(l.PBound, l.PBound)
	fn := new(Fact)
	fn.Fac = a
	fn.Exp = f.Exp
	if bnd.Cmp(a) > 0 {
		fn.Stat = Prime
	} else {
		fn.Stat = Unknown
	}
	l.Insert(fn)
	fn = new(Fact)
	fn.Fac = b
	fn.Exp = f.Exp
	if bnd.Cmp(b) > 0 {
		fn.Stat = Prime
	} else {
		fn.Stat = Unknown
	}
	l.Insert(fn)
}

// IsComplete checks if the factorisation is complete.
// It returns 0 if there are still unknown or composite factors,
// 1 if all factors are at least probably prime, and
// 2 if all factors are prime.
func (l *Factors) IsComplete() int {
	res := 2
	for f := l.First; f != nil; f = f.Next {
		switch f.Stat {
		case Unknown, Composite:
			return 0
		case ProbPrime:
			res = 1
		}
	}
	return res
}

// PrimTest runs a primality test on the factors and updates their status.
// The test is done by calling func (*big.Int) ProbablyPrime(n)
// If retest is true checks again probably prime factors.
func (l *Factors) PrimTest(n int, retest bool) {
	for f := l.First; f != nil; f = f.Next {
		if f.Stat == Unknown || retest && f.Stat == ProbPrime {
			if f.Fac.ProbablyPrime(n) {
				f.Stat = ProbPrime
			} else {
				f.Stat = Composite
			}
		}
	}
}

// Insert is a low level function that adds a factor to the list.
// This operation does not preserve the product.
func (l *Factors) Insert(f *Fact) {
	var pp **Fact
	var cmp int
	for pp = &l.First; *pp != nil; pp = &(*pp).Next {
		cmp = (*pp).Fac.Cmp(f.Fac)
		if cmp == 0 {
			(*pp).Exp += f.Exp
			(*pp).Stat = mergeStat((*pp).Stat, f.Stat)
			return
		}
		if cmp == 1 {
			break
		}
	}
	f.Next = *pp
	*pp = f
}

func mergeStat(a, b Status) Status {
	switch {
	case a == Unknown:
		return b
	case b == Unknown:
		return a
	case a == ProbPrime:
		return b
	case b == ProbPrime:
		return a
	case a == b:
		return a
	}
	panic(errors.New("incompatible factor stati"))
}
