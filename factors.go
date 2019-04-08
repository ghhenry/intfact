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
