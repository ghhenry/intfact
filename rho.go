package intfact

import (
	"context"
	"errors"
	"math/big"
)

// Rho tries to factor n with Pollard's rho method.
func Rho(ctx context.Context, n *big.Int) (fac *big.Int, err error) {
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
	gcd := newGcdtest(n, 20)
	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("cancelled")
		default:
		}
		d := new(big.Int).Sub(h, l)
		d.Abs(d)
		fac, err = gcd.test(d)
		if fac != nil || err != nil {
			return
		}
		l = f(l)
		h = f(f(h))
	}
}
