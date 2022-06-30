package intfact

import (
	"context"
	"errors"
	"io"
	"math/big"
)

func EcParallel(ctx context.Context, random io.Reader, n *big.Int, b, b1 uint32) (*big.Int, error) {
	pf := 100
	childctx, cancel := context.WithCancel(ctx)
	defer cancel()
	type result struct {
		fac *big.Int
		err error
	}
	resultC := make(chan result)
	for i := 0; i < pf; i++ {
		go func() {
			fac, err := Ec(childctx, random, n, b, b1)
			select {
			case <-childctx.Done():
				return
			case resultC <- result{fac, err}:
				return
			}
		}()
	}
	for finished := 0; finished < pf; finished++ {
		var r result
		select {
		case r = <-resultC:
			if r.err == nil {
				return r.fac, nil
			}
		case <-ctx.Done():
			return nil, errors.New("cancelled")
		}
	}
	return nil, errors.New("no factor found")
}
