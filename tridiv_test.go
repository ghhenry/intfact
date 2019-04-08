package intfact

import (
	"math/big"
	"testing"
)

func TestTridiv(t *testing.T) {
	type tf struct {
		f *big.Int
		e uint
		s Status
	}
	var tests = []struct {
		n *big.Int
		b uint32
		l []tf
	}{
		{
			n: big.NewInt(1),
			b: 3,
			l: []tf{
				{
					f: big.NewInt(1),
					e: 1,
					s: Prime,
				},
			},
		},
		{
			n: big.NewInt(2),
			b: 3,
			l: []tf{
				{
					f: big.NewInt(2),
					e: 1,
					s: Prime,
				},
			},
		},
		{
			n: big.NewInt(4),
			b: 3,
			l: []tf{
				{
					f: big.NewInt(2),
					e: 2,
					s: Prime,
				},
			},
		},
		{
			n: big.NewInt(6),
			b: 3,
			l: []tf{
				{
					f: big.NewInt(2),
					e: 1,
					s: Prime,
				},
				{
					f: big.NewInt(3),
					e: 1,
					s: Prime,
				},
			},
		},
		{
			n: big.NewInt(20),
			b: 3,
			l: []tf{
				{
					f: big.NewInt(2),
					e: 2,
					s: Prime,
				},
				{
					f: big.NewInt(5),
					e: 1,
					s: Prime,
				},
			},
		},
		{
			n: big.NewInt(44),
			b: 3,
			l: []tf{
				{
					f: big.NewInt(2),
					e: 2,
					s: Prime,
				},
				{
					f: big.NewInt(11),
					e: 1,
					s: Unknown,
				},
			},
		},
		{
			n: big.NewInt(44),
			b: 5,
			l: []tf{
				{
					f: big.NewInt(2),
					e: 2,
					s: Prime,
				},
				{
					f: big.NewInt(11),
					e: 1,
					s: Prime,
				},
			},
		},
	}
	for _, test := range tests {
		l := NewFactors(test.n)
		l.TrialDivision(test.b)
		if big.NewInt(int64(test.b)).Cmp(l.PBound) != 0 {
			t.Errorf("got bound %v, want %v", l.PBound, test.b)
		}
		i := 0
		for p := l.First; p != nil; p = p.Next {
			if i >= len(test.l) {
				t.Errorf("got too many factors, want %v", len(test.l))
			}
			if test.l[i].f.Cmp(p.Fac) != 0 {
				t.Errorf("got %v factor %v, want %v", i, p.Fac, test.l[i].f)
			}
			if test.l[i].e != p.Exp {
				t.Errorf("got %v exponent %v, want %v", i, p.Exp, test.l[i].e)
			}
			if test.l[i].s != p.Stat {
				t.Errorf("got %v status %v, want %v", i, p.Stat, test.l[i].s)
			}
			i++
		}
		if i != len(test.l) {
			t.Errorf("got %v factors, want %v", i, len(test.l))
		}
	}
}
