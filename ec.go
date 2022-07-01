package intfact

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/ghhenry/primes"
	"io"
	"math/big"
	"sync"
)

type lcRandom struct {
	x     uint32
	mutex sync.Mutex
}

// use a simple implementation, the quality of the random numbers does not matter
func (r *lcRandom) Read(p []byte) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for i := range p {
		r.x = 494131989*r.x + 998936465
		p[i] = byte(r.x >> 24)
	}
	return len(p), nil
}

type factorError struct {
	f *big.Int
}

func (e factorError) Error() string {
	return "modulus is not prime, a factor is: " + e.f.String()
}

type curve struct {
	n *big.Int
	a *big.Int
	b *big.Int
}

type point interface {
	isZero() bool
	equal(o point) bool
	x() *big.Int
	y() *big.Int
}

type neutral struct{}

func (neutral) isZero() bool { return true }

func (neutral) equal(o point) bool { return o.isZero() }

func (neutral) x() *big.Int {
	panic(errors.New("can't get coordinate of zero point"))
}

func (neutral) y() *big.Int {
	panic(errors.New("can't get coordinate of zero point"))
}

func (neutral) String() string {
	return "zero"
}

type ordinary struct {
	px, py *big.Int
}

func (ordinary) isZero() bool { return false }

func (p ordinary) equal(o point) bool {
	return !o.isZero() && p.x().Cmp(o.x()) == 0 && p.y().Cmp(o.y()) == 0
}

func (p ordinary) x() *big.Int { return p.px }

func (p ordinary) y() *big.Int { return p.py }

func (p ordinary) String() string {
	return fmt.Sprintf("(%v, %v)", p.px, p.py)
}

func randCurve(random io.Reader, n *big.Int) (*curve, point) {
	for {
		px, _ := rand.Int(random, n)
		py, _ := rand.Int(random, n)
		p := ordinary{px, py}
		a, _ := rand.Int(random, n)
		// calculate b as y**2 - x(x**2+a)
		b := new(big.Int).Mul(py, py)
		t0 := new(big.Int).Mul(px, px)
		t0.Add(t0, a)
		t0.Mod(t0, n)
		t0.Mul(t0, px)
		b.Sub(b, t0)
		b.Mod(b, n)
		c := &curve{n, a, b}
		if c.isNonSingular() {
			return c, p
		}
	}
}

func (c *curve) neg(a point) point {
	if a.isZero() {
		return a
	}
	y := new(big.Int).Neg(a.y())
	y.Mod(y, c.n)
	return ordinary{a.x(), y}
}

func (c *curve) add(a, b point) (point, error) {
	if a.isZero() {
		return b, nil
	}
	if b.isZero() {
		return a, nil
	}
	if a.equal(b) {
		return c.double(a)
	}
	if a.equal(c.neg(b)) {
		return neutral{}, nil
	}
	return c.rawAdd(a, b)
}

func (c *curve) double(a point) (point, error) {
	if a.isZero() {
		return a, nil
	}
	if a.y().Sign() == 0 {
		return neutral{}, nil
	}
	t0 := new(big.Int).Add(a.y(), a.y())
	t0.Mod(t0, c.n)
	t1 := new(big.Int).GCD(t0, nil, t0, c.n)
	if t1.Cmp(bigOne) != 0 {
		return nil, factorError{t1}
	}
	t1.Mul(a.x(), a.x())
	t1.Mod(t1, c.n)
	t1.Mul(t1, bigThree)
	t1.Add(t1, c.a)
	t0.Mul(t1, t0)
	t0.Mod(t0, c.n)
	t1.Mul(t0, t0)
	t1.Sub(t1, a.x())
	t1.Sub(t1, a.x())
	t1.Mod(t1, c.n)
	t2 := new(big.Int).Sub(a.x(), t1)
	t0.Mul(t0, t2)
	t0.Sub(t0, a.y())
	t0.Mod(t0, c.n)
	return ordinary{t1, t0}, nil
}

func (c *curve) rawAdd(a, b point) (point, error) {
	// we can assume that a and b are not zero and that a.x != b.x
	t0 := new(big.Int).Sub(a.x(), b.x())
	t0.Mod(t0, c.n)
	t1 := new(big.Int).GCD(t0, nil, t0, c.n)
	if t1.Cmp(bigOne) != 0 {
		return nil, factorError{t1}
	}
	t1.Sub(a.y(), b.y())
	t0.Mul(t1, t0)
	t0.Mod(t0, c.n)
	t1.Mul(t0, t0)
	t1.Sub(t1, a.x())
	t1.Sub(t1, b.x())
	t1.Mod(t1, c.n) // this is the x-coordinate of the sum
	t2 := new(big.Int).Sub(a.x(), t1)
	t0.Mul(t0, t2)
	t0.Sub(t0, a.y())
	t0.Mod(t0, c.n) // this is the y-coordinate of the sum
	return ordinary{t1, t0}, nil
}

func (c *curve) isNonSingular() bool {
	t0 := new(big.Int).Mul(c.a, c.a)
	t0.Mod(t0, c.n)
	t0.Mul(t0, c.a)
	t0.Mod(t0, c.n)
	t0.Add(t0, t0)
	t0.Add(t0, t0)
	t1 := new(big.Int).Mul(c.b, c.b)
	t1.Mod(t1, c.n)
	t1.Mul(t1, big.NewInt(27))
	t0.Add(t0, t1)
	t0.Mod(t0, c.n)
	return t0.Sign() != 0
}

func (c *curve) String() string {
	return fmt.Sprintf("a=%v, b=%v over GF(%v)", c.a, c.b, c.n)
}

func (c *curve) mult(p point, m *big.Int) (point, error) {
	if m.Sign() < 0 {
		m = new(big.Int).Neg(m)
		p = c.neg(p)
	} else {
		// we change m, so make a copy
		m = new(big.Int).Set(m)
	}
	var r point = neutral{}
	for m.Sign() != 0 {
		var err error
		if m.Bit(0) == 1 {
			m.SetBit(m, 0, 0)
			r, err = c.add(r, p)
		} else {
			m.Rsh(m, 1)
			p, err = c.double(p)
		}
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

type p2echelper struct {
	c         *curve
	multiples []point
}

func newEcHelper(a point, c *curve) (*p2echelper, error) {
	p2, err := c.double(a)
	if err != nil {
		return nil, err
	}
	return &p2echelper{
		c:         c,
		multiples: []point{p2},
	}, nil
}

func (h *p2echelper) getMultiple(e uint32) (point, error) {
	i := int(e/2 - 1)
	for i >= len(h.multiples) {
		k, err := h.c.add(h.multiples[0], h.multiples[len(h.multiples)-1])
		if err != nil {
			return nil, err
		}
		h.multiples = append(h.multiples, k)
	}
	return h.multiples[i], nil
}

func Ec(ctx context.Context, random io.Reader, n *big.Int, b, b1 uint32) (fac *big.Int, err error) {
	c, pt := randCurve(random, n)
	phase1 := func(p uint32) bool {
		select {
		case <-ctx.Done():
			err = errors.New("cancelled")
			return true
		default:
		}
		mult := int64(p)
		for {
			ne := mult * int64(p)
			if ne > int64(b) {
				break
			}
			mult = ne
		}
		pt, err = c.mult(pt, big.NewInt(mult))
		if err != nil {
			if e, ok := err.(factorError); ok {
				//fmt.Println("factor found in phase1 at", p)
				fac = e.f
				err = nil
			}
			return true
		}
		if pt.isZero() {
			err = errors.New("no factor found")
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
	h, err := newEcHelper(pt, c)
	if err != nil {
		if e, ok := err.(factorError); ok {
			fac = e.f
			err = nil
		}
		return
	}
	phase2 := func(p uint32) bool {
		select {
		case <-ctx.Done():
			err = errors.New("cancelled")
			return true
		default:
		}
		if prev == 0 {
			pt, err = c.mult(pt, big.NewInt(int64(p)))
		} else {
			diff := p - prev
			var ptInc point
			ptInc, err = h.getMultiple(diff)
			if err == nil {
				pt, err = c.add(pt, ptInc)
			}
		}
		if err != nil {
			if e, ok := err.(factorError); ok {
				//fmt.Println("factor found in phase2 at", p)
				fac = e.f
				err = nil
			}
			return true
		}
		prev = p
		return false
	}
	primes.Iterate(b+1, b1, phase2)
	if fac != nil || err != nil {
		return
	}
	return nil, errors.New("no factor found")
}
