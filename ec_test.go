package intfact

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"testing"
)

func TestRandom(t *testing.T) {
	r := &lcRandom{x: 2}
	p := make([]byte, 10)
	n, err := r.Read(p)
	if err != nil {
		t.Error("unexpected error", err)
	}
	if n != 10 {
		t.Error("invalid n", n)
	}
}

func onCurve(c *curve, p point) bool {
	if p.isZero() {
		return true
	}
	left := new(big.Int).Mul(p.x(), p.x())
	left.Add(left, c.a)
	left.Mul(left, p.x())
	left.Add(left, c.b)
	left.Mod(left, c.n)
	right := new(big.Int).Mul(p.y(), p.y())
	right.Mod(right, c.n)
	return left.Cmp(right) == 0
}

func TestRandCurve(t *testing.T) {
	r := &lcRandom{x: 2}
	c, p := randCurve(r, big.NewInt(47))
	if !onCurve(c, p) {
		t.Error("point is not on the curve")
	}
}

func Test_curve_double(t *testing.T) {
	type fields struct {
		n *big.Int
		a *big.Int
		b *big.Int
	}
	type args struct {
		a point
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    point
		wantErr bool
	}{
		{
			name: "zero",
			fields: fields{
				n: big.NewInt(47),
				a: big.NewInt(2),
				b: big.NewInt(3),
			},
			args: args{
				neutral{},
			},
			want:    neutral{},
			wantErr: false,
		},
		{
			name: "normal",
			fields: fields{
				n: big.NewInt(47),
				a: big.NewInt(2),
				b: big.NewInt(3),
			},
			args: args{
				ordinary{
					px: big.NewInt(0),
					py: big.NewInt(12),
				},
			},
			want: ordinary{
				px: big.NewInt(16),
				py: big.NewInt(18),
			},
			wantErr: false,
		},
		{
			name: "2-torsion",
			fields: fields{
				n: big.NewInt(47),
				a: big.NewInt(2),
				b: big.NewInt(3),
			},
			args: args{
				ordinary{
					px: big.NewInt(21),
					py: big.NewInt(0),
				},
			},
			want:    neutral{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &curve{
				n: tt.fields.n,
				a: tt.fields.a,
				b: tt.fields.b,
			}
			if !onCurve(c, tt.args.a) {
				t.Errorf("point a is not on the curve")
			}
			if !onCurve(c, tt.want) {
				t.Errorf("expected result is not on the curve")
			}
			got, err := c.double(tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("double() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("double() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_curve_rawAdd(t *testing.T) {
	type fields struct {
		n *big.Int
		a *big.Int
		b *big.Int
	}
	type args struct {
		a point
		b point
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    point
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				n: big.NewInt(47),
				a: big.NewInt(2),
				b: big.NewInt(3),
			},
			args: args{
				a: ordinary{
					px: big.NewInt(0),
					py: big.NewInt(12),
				},
				b: ordinary{
					px: big.NewInt(1),
					py: big.NewInt(10),
				},
			},
			want: ordinary{
				px: big.NewInt(3),
				py: big.NewInt(41),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &curve{
				n: tt.fields.n,
				a: tt.fields.a,
				b: tt.fields.b,
			}
			if !onCurve(c, tt.args.a) {
				t.Errorf("point a is not on the curve")
			}
			if !onCurve(c, tt.args.b) {
				t.Errorf("point b is not on the curve")
			}
			if !onCurve(c, tt.want) {
				t.Errorf("expected result is not on the curve")
			}
			got, err := c.rawAdd(tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("rawAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rawAdd() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderGf47(t *testing.T) {
	c := &curve{
		// this curve is isomorphic to Z2*Z24
		n: big.NewInt(47),
		a: big.NewInt(2),
		b: big.NewInt(3),
	}
	p := ordinary{
		// this point generates a subgroup of order 24
		px: big.NewInt(12),
		py: big.NewInt(4),
	}
	if ord := orderCalc(t, c, p); ord != 24 {
		t.Errorf("got order %v, want 24", ord)
	}
}

func TestOrderGf53(t *testing.T) {
	c := &curve{
		// this curve is isomorphic to Z42
		n: big.NewInt(53),
		a: big.NewInt(51),
		b: big.NewInt(42),
	}
	p := ordinary{
		// this is a generator of the curve
		px: big.NewInt(28),
		py: big.NewInt(46),
	}
	if ord := orderCalc(t, c, p); ord != 42 {
		t.Errorf("got order %v, want 42", ord)
	}
}

func TestOrderRandom(t *testing.T) {
	c, p := randCurve(&lcRandom{x: 1}, big.NewInt(2503))
	orderCalc(t, c, p)
}

func orderCalc(t *testing.T, c *curve, p point) int {
	fmt.Println("curve", c)
	var q point = neutral{}
	var err error
	order := 0
	for {
		order++
		q, err = c.add(q, p)
		if err != nil {
			t.Fatal("addition failed:", err)
		}
		if !onCurve(c, q) {
			t.Fatal("point is not on curve", q)
		}
		//fmt.Println(order, q)
		if q.isZero() {
			break
		}
	}
	return order
}

func TestFactorFound(t *testing.T) {
	c := &curve{
		n: big.NewInt(2491),
		a: big.NewInt(906),
		b: big.NewInt(956),
	}
	p := ordinary{
		px: big.NewInt(2276),
		py: big.NewInt(443),
	}
	q := ordinary{
		px: big.NewInt(421),
		py: big.NewInt(1041),
	}
	_, err := c.add(p, q)
	if err == nil {
		t.Fatal("an error was expected")
	}
	if e, ok := err.(factorError); ok {
		if e.f.Cmp(big.NewInt(47)) != 0 && e.f.Cmp(big.NewInt(53)) != 0 {
			t.Error("invalid factor", e.f)
		}
	} else {
		t.Error("unexpected error", err)
	}
}

func Test_curve_mult(t *testing.T) {
	type fields struct {
		n *big.Int
		a *big.Int
		b *big.Int
	}
	type args struct {
		p point
		m *big.Int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    point
		wantErr bool
	}{
		{
			name: "-1",
			fields: fields{
				n: big.NewInt(47),
				a: big.NewInt(2),
				b: big.NewInt(3),
			},
			args: args{
				p: ordinary{
					px: big.NewInt(3),
					py: big.NewInt(6),
				},
				m: big.NewInt(-1),
			},
			want: ordinary{
				px: big.NewInt(3),
				py: big.NewInt(41),
			},
			wantErr: false,
		},
		{
			name: "0",
			fields: fields{
				n: big.NewInt(47),
				a: big.NewInt(2),
				b: big.NewInt(3),
			},
			args: args{
				p: ordinary{
					px: big.NewInt(3),
					py: big.NewInt(6),
				},
				m: big.NewInt(0),
			},
			want:    neutral{},
			wantErr: false,
		},
		{
			name: "1",
			fields: fields{
				n: big.NewInt(47),
				a: big.NewInt(2),
				b: big.NewInt(3),
			},
			args: args{
				p: ordinary{
					px: big.NewInt(3),
					py: big.NewInt(6),
				},
				m: big.NewInt(1),
			},
			want: ordinary{
				px: big.NewInt(3),
				py: big.NewInt(6),
			},
			wantErr: false,
		},
		{
			name: "2",
			fields: fields{
				n: big.NewInt(47),
				a: big.NewInt(2),
				b: big.NewInt(3),
			},
			args: args{
				p: ordinary{
					px: big.NewInt(12),
					py: big.NewInt(4),
				},
				m: big.NewInt(2),
			},
			want: ordinary{
				px: big.NewInt(8),
				py: big.NewInt(25),
			},
			wantErr: false,
		},
		{
			name: "12",
			fields: fields{
				n: big.NewInt(47),
				a: big.NewInt(2),
				b: big.NewInt(3),
			},
			args: args{
				p: ordinary{
					px: big.NewInt(12),
					py: big.NewInt(4),
				},
				m: big.NewInt(12),
			},
			want: ordinary{
				px: big.NewInt(27),
				py: big.NewInt(0),
			},
			wantErr: false,
		},
		{
			name: "24",
			fields: fields{
				n: big.NewInt(47),
				a: big.NewInt(2),
				b: big.NewInt(3),
			},
			args: args{
				p: ordinary{
					px: big.NewInt(3),
					py: big.NewInt(6),
				},
				m: big.NewInt(24),
			},
			want:    neutral{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &curve{
				n: tt.fields.n,
				a: tt.fields.a,
				b: tt.fields.b,
			}
			got, err := c.mult(tt.args.p, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("mult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.equal(tt.want) {
				t.Errorf("mult() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEc(t *testing.T) {
	type args struct {
		ctx    context.Context
		random io.Reader
		n      *big.Int
		b      uint32
		b1     uint32
	}
	tests := []struct {
		name    string
		args    args
		wantFac *big.Int
		wantErr bool
	}{
		{
			name: "2491",
			args: args{
				ctx:    context.Background(),
				random: &lcRandom{x: 36},
				n:      big.NewInt(2491),
				b:      10,
				b1:     100,
			},
			wantFac: big.NewInt(53),
			wantErr: false,
		},
		{
			name: "n1",
			args: args{
				ctx:    context.Background(),
				random: &lcRandom{x: 13},
				n:      big.NewInt(43217358712783469),
				b:      1000,
				b1:     10000,
			},
			wantFac: big.NewInt(7420146347),
			wantErr: false,
		},
		{
			name: "f6",
			args: args{
				ctx:    context.Background(),
				random: &lcRandom{x: 14},
				n:      intval("18446744073709551617"),
				b:      1000,
				b1:     10000,
			},
			wantFac: big.NewInt(274177),
			wantErr: false,
		},
		{
			name: "f7",
			args: args{
				ctx:    context.Background(),
				random: &lcRandom{x: 16},
				n:      intval("340282366920938463463374607431768211457"),
				b:      10000,
				b1:     215000,
			},
			wantFac: big.NewInt(59649589127497217),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFac, err := Ec(tt.args.ctx, tt.args.random, tt.args.n, tt.args.b, tt.args.b1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFac, tt.wantFac) {
				t.Errorf("Ec() gotFac = %v, want %v", gotFac, tt.wantFac)
			}
		})
	}
}
