package intfact

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
)

func TestEcParallel(t *testing.T) {
	// has a factor 59649589127497217
	fac, err := EcParallel(context.Background(), &lcRandom{x: 10}, intval("340282366920938463463374607431768211457"), 10000, 215000, 200)
	fmt.Println(fac, err)
}

func TestEcParallel1(t *testing.T) {
	// has a factor 24844087002162188818957
	fac, err := EcParallel(context.Background(), &lcRandom{x: 10}, intval("51975989010489170378207397897900519326437583043996359"), 20000, 550000, 500)
	fmt.Println(fac, err)
}

func TestEcParallel2(t *testing.T) {
	r := &lcRandom{x: 10}
	p1, _ := rand.Prime(r, 90)
	p2, _ := rand.Prime(r, 110)
	n := new(big.Int).Mul(p1, p2)
	fmt.Println(n, p1, p2)
	fac, err := EcParallel(context.Background(), r, n, 40000, 1400000, 2000)
	fmt.Println(fac, err)
}
