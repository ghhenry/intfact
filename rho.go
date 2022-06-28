package intfact

import (
	"context"
	"errors"
	"math/big"
)

var bigOne = big.NewInt(1)

// Rho tries to factor n with Pollard's rho method.
func Rho(ctx context.Context, n *big.Int) (*big.Int, error) {
	a := big.NewInt(2)
	c := &big.Int{}
	c.SetUint64(1)
	f := func(x *big.Int) *big.Int {
		r := new(big.Int).Mul(x, x)
		r.Add(r, c)
		r.Mod(r, n)
		return r
	}

	l := f(a)
	h := f(f(a))
	it := 0
	t := big.NewInt(1)
	for {
		it++
		select {
		case <-ctx.Done():
			return nil, errors.New("cancelled")
		default:
		}
		d := new(big.Int).Sub(h, l)
		d.Abs(d)
		t.Mul(t, d)
		t.Mod(t, n)
		if it%10 == 0 {
			// since GCD is expensive we check only every tenth iteration
			// this may collapse some small factors
			r := new(big.Int).GCD(nil, nil, t, n)
			if r.Cmp(bigOne) != 0 {
				if r.Cmp(n) != 0 {
					return r, nil
				} else {
					return nil, errors.New("no factor found")
				}
			}
			t.SetInt64(1)
		}
		l = f(l)
		h = f(f(h))
	}
}
