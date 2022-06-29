package intfact

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"
)

func TestRho(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	// n=2^2^7+1
	// takes 1227.18s in original version
	// takes 731.99s in optimized version
	//n, _ := new(big.Int).SetString("340282366920938463463374607431768211457", 10)
	// n=2^67-1
	//n, _ := new(big.Int).SetString("147573952589676412927", 10)
	n, _ := new(big.Int).SetString("43217358712783469", 10)
	d, err := Rho(ctx, n)
	fmt.Printf("d=%v, e=%v\n", d, err)
}
