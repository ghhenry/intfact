package intfact

import (
	"context"
	"math/big"
	"reflect"
	"testing"
)

func intval(s string) *big.Int {
	var a = &big.Int{}
	a.SetString(s, 10)
	return a
}

func TestPmOne(t *testing.T) {
	type args struct {
		ctx context.Context
		n   *big.Int
		b   uint32
		b1  uint32
	}
	tests := []struct {
		name    string
		args    args
		wantFac *big.Int
		wantErr bool
	}{
		{
			name: "phase1",
			args: args{
				ctx: context.Background(),
				n:   big.NewInt(41 * 3803),
				b:   10,
				b1:  100,
			},
			wantFac: big.NewInt(41),
			wantErr: false,
		},
		{
			name: "phase2",
			args: args{
				ctx: context.Background(),
				n:   big.NewInt(3607 * 3803),
				b:   10,
				b1:  700,
			},
			wantFac: big.NewInt(3607),
			wantErr: false,
		},
		{
			name: "n1",
			args: args{
				ctx: context.Background(),
				n:   big.NewInt(43217358712783469),
				b:   1000,
				b1:  10000,
			},
			wantFac: big.NewInt(7420146347),
			wantErr: false,
		},
		{
			name: "f6",
			args: args{
				ctx: context.Background(),
				n:   intval("18446744073709551617"),
				b:   300,
				b1:  1000,
			},
			wantFac: big.NewInt(274177),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFac, err := PmOne(tt.args.ctx, tt.args.n, tt.args.b, tt.args.b1)
			if (err != nil) != tt.wantErr {
				t.Errorf("PmOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFac, tt.wantFac) {
				t.Errorf("PmOne() gotFac = %v, want %v", gotFac, tt.wantFac)
			}
		})
	}
}
