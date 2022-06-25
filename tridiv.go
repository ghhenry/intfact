package intfact

import (
	"math/big"

	"github.com/ghhenry/primes"
)

// TrialDivision tries to factor the list by trial division with small primes.
func (l *Factors) TrialDivision(bound uint32) {
	oldlist := l.First
	l.First = nil
	for f := oldlist; f != nil; f = f.Next {
		v := f.Fac
		primes.Iterate(uint32(l.PBound.Uint64()), bound, func(p uint32) bool {
			t := big.NewInt(int64(p))
			t2 := new(big.Int).Mul(t, t)
			for {
				if v.Cmp(t2) < 0 {
					l.Insert(&Fact{Fac: v, Exp: f.Exp, Stat: Prime})
					v = nil
					return true
				}
				if primes.Fastmod(v, p) != 0 {
					return false
				}
				d := new(big.Int).Div(v, t)
				l.Insert(&Fact{Fac: t, Exp: f.Exp, Stat: Prime})
				v = d
			}
		})
		if v != nil {
			l.Insert(&Fact{Fac: v, Exp: f.Exp, Stat: Unknown})
		}
	}
	l.PBound = big.NewInt(int64(bound))
}
