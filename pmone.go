package intfact

import (
	"context"
	"errors"
	"github.com/ghhenry/primes"
	"math/big"
)

type p2helper struct {
	m      *big.Int
	powers []*big.Int
}

func newHelper(a, n *big.Int) *p2helper {
	a2 := new(big.Int).Mul(a, a)
	a2.Mod(a2, n)
	return &p2helper{
		m:      n,
		powers: []*big.Int{a2},
	}
}

func (h *p2helper) getPower(e uint32) *big.Int {
	i := int(e/2 - 1)
	for i >= len(h.powers) {
		k := new(big.Int).Mul(h.powers[0], h.powers[len(h.powers)-1])
		k.Mod(k, h.m)
		h.powers = append(h.powers, k)
	}
	return h.powers[i]
}

func PmOne(ctx context.Context, n *big.Int, b, b1 uint32) (fac *big.Int, err error) {
	var a = big.NewInt(3)
	gcd := newGcdtest(n, 20)
	phase1 := func(p uint32) bool {
		select {
		case <-ctx.Done():
			err = errors.New("cancelled")
			return true
		default:
		}
		exp := int64(p)
		for {
			ne := exp * int64(p)
			if ne > int64(b) {
				break
			}
			exp = ne
		}
		a.Exp(a, big.NewInt(exp), n)
		d := new(big.Int).Sub(a, bigOne)
		fac, err = gcd.test(d)
		if fac != nil || err != nil {
			return true
		}
		return false
	}
	primes.Iterate(2, b, phase1)
	if fac != nil || err != nil {
		return
	}

	// phase2
	var prev uint32
	h := newHelper(a, n)
	phase2 := func(p uint32) bool {
		select {
		case <-ctx.Done():
			err = errors.New("cancelled")
			return true
		default:
		}
		if prev == 0 {
			a.Exp(a, big.NewInt(int64(p)), n)
		} else {
			diff := p - prev
			a.Mul(a, h.getPower(diff))
			a.Mod(a, n)
		}
		d := new(big.Int).Sub(a, bigOne)
		fac, err = gcd.test(d)
		if fac != nil || err != nil {
			return true
		}
		prev = p
		return false
	}
	primes.Iterate(b+1, b1, phase2)
	if fac != nil || err != nil {
		return
	}
	fac, err = gcd.finish()
	if fac != nil || err != nil {
		return
	}
	return nil, errors.New("no factor found")
}
